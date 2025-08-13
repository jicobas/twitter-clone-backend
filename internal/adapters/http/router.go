package http

import (
	"net/http"
	"strings"
)

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

	// API routes
	mux.HandleFunc("/api/v1/tweets", handlers.CreateTweet)
	mux.HandleFunc("/api/v1/tweets/user/", handlers.GetUserTweets)
	mux.HandleFunc("/api/v1/timeline/", handlers.GetTimeline)
	mux.HandleFunc("/api/v1/follow", handlers.FollowUser)
	mux.HandleFunc("/api/v1/follow/", handlers.UnfollowUser)
	mux.HandleFunc("/api/v1/users/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if strings.HasSuffix(path, "/followers") {
			handlers.GetFollowers(w, r)
		} else if strings.HasSuffix(path, "/following") {
			handlers.GetFollowing(w, r)
		} else {
			http.NotFound(w, r)
		}
	})

	return corsHandler(mux)
}
