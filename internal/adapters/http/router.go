package http

import (
	"net/http"
	"strings"
)

// methodHandler wraps an HTTP handler with method validation
func methodHandler(method string, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte(`{"error": "Method not allowed"}`))
			return
		}
		handler(w, r)
	}
}

// SetupRoutes configura todas las rutas de la API usando net/http
func SetupRoutes(handlers *Handlers) http.Handler {
	mux := http.NewServeMux()

	// CORS middleware
	corsHandler := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, X-User-ID")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	})

	// API routes - Clean REST endpoints con validación de métodos
	mux.HandleFunc("/tweets", methodHandler("POST", handlers.CreateTweet))
	mux.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if strings.HasSuffix(path, "/tweets") {
			methodHandler("GET", handlers.GetUserTweets)(w, r)
		} else if strings.HasSuffix(path, "/timeline") {
			methodHandler("GET", handlers.GetTimeline)(w, r)
		} else if strings.HasSuffix(path, "/followers") {
			methodHandler("GET", handlers.GetFollowers)(w, r)
		} else if strings.HasSuffix(path, "/following") {
			methodHandler("GET", handlers.GetFollowing)(w, r)
		} else {
			http.NotFound(w, r)
		}
	})
	mux.HandleFunc("/users/following", methodHandler("POST", handlers.FollowUser))
	mux.HandleFunc("/users/following/", methodHandler("DELETE", handlers.UnfollowUser))

	return corsHandler(mux)
}
