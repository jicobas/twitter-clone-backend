package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"twitter-clone-backend/internal/usecases"
)

// Handlers contiene los manejadores HTTP
type Handlers struct {
	tweetUseCase  *usecases.TweetUseCase
	followUseCase *usecases.FollowUseCase
}

// NewHandlers crea una nueva instancia de los manejadores
func NewHandlers(tweetUseCase *usecases.TweetUseCase, followUseCase *usecases.FollowUseCase) *Handlers {
	return &Handlers{
		tweetUseCase:  tweetUseCase,
		followUseCase: followUseCase,
	}
}

// Tweet request/response structures
type CreateTweetRequest struct {
	Content string `json:"content"`
}

type TweetResponse struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

type FollowRequest struct {
	FolloweeID string `json:"followee_id"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type FollowersResponse struct {
	Followers []string `json:"followers"`
}

type FollowingResponse struct {
	Following []string `json:"following"`
}

// writeJSON escribe una respuesta JSON
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// writeError escribe una respuesta de error JSON
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, ErrorResponse{Error: message})
}

// CreateTweet maneja la creación de un nuevo tweet
func (h *Handlers) CreateTweet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		writeError(w, http.StatusBadRequest, "X-User-ID header is required")
		return
	}

	var req CreateTweetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if req.Content == "" {
		writeError(w, http.StatusBadRequest, "Content is required")
		return
	}

	if len(req.Content) > 280 {
		writeError(w, http.StatusBadRequest, "Content exceeds 280 characters")
		return
	}

	tweet, err := h.tweetUseCase.CreateTweet(r.Context(), userID, req.Content)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := TweetResponse{
		ID:        tweet.ID,
		UserID:    tweet.UserID,
		Content:   tweet.Content,
		CreatedAt: tweet.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}

	writeJSON(w, http.StatusCreated, response)
}

// GetTimeline obtiene el timeline de un usuario
func (h *Handlers) GetTimeline(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extraer userID del path (asumiendo formato /timeline/{userID})
	userID := r.URL.Path[len("/api/v1/timeline/"):]
	if userID == "" {
		writeError(w, http.StatusBadRequest, "userID parameter is required")
		return
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 50 // valor por defecto
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil {
			limit = parsedLimit
		}
	}

	tweets, err := h.tweetUseCase.GetTimeline(r.Context(), userID, limit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var response []TweetResponse
	for _, tweet := range tweets {
		response = append(response, TweetResponse{
			ID:        tweet.ID,
			UserID:    tweet.UserID,
			Content:   tweet.Content,
			CreatedAt: tweet.CreatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}

	writeJSON(w, http.StatusOK, response)
}

// GetUserTweets obtiene todos los tweets de un usuario
func (h *Handlers) GetUserTweets(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extraer userID del path (asumiendo formato /tweets/user/{userID})
	userID := r.URL.Path[len("/api/v1/tweets/user/"):]
	if userID == "" {
		writeError(w, http.StatusBadRequest, "userID parameter is required")
		return
	}

	tweets, err := h.tweetUseCase.GetUserTweets(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var response []TweetResponse
	for _, tweet := range tweets {
		response = append(response, TweetResponse{
			ID:        tweet.ID,
			UserID:    tweet.UserID,
			Content:   tweet.Content,
			CreatedAt: tweet.CreatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}

	writeJSON(w, http.StatusOK, response)
}

// FollowUser permite a un usuario seguir a otro
func (h *Handlers) FollowUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	followerID := r.Header.Get("X-User-ID")
	if followerID == "" {
		writeError(w, http.StatusBadRequest, "X-User-ID header is required")
		return
	}

	var req FollowRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if req.FolloweeID == "" {
		writeError(w, http.StatusBadRequest, "followee_id is required")
		return
	}

	err := h.followUseCase.FollowUser(r.Context(), followerID, req.FolloweeID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, MessageResponse{Message: "successfully followed user"})
}

// UnfollowUser permite a un usuario dejar de seguir a otro
func (h *Handlers) UnfollowUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	followerID := r.Header.Get("X-User-ID")
	if followerID == "" {
		writeError(w, http.StatusBadRequest, "X-User-ID header is required")
		return
	}

	// Extraer followeeID del path (asumiendo formato /follow/{followeeID})
	followeeID := r.URL.Path[len("/api/v1/follow/"):]
	if followeeID == "" {
		writeError(w, http.StatusBadRequest, "followeeID parameter is required")
		return
	}

	err := h.followUseCase.UnfollowUser(r.Context(), followerID, followeeID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, MessageResponse{Message: "successfully unfollowed user"})
}

// GetFollowers obtiene los seguidores de un usuario
func (h *Handlers) GetFollowers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extraer userID del path (asumiendo formato /users/{userID}/followers)
	path := r.URL.Path
	userID := ""
	if len(path) > len("/api/v1/users/") {
		endIndex := len(path) - len("/followers")
		if endIndex > len("/api/v1/users/") {
			userID = path[len("/api/v1/users/"):endIndex]
		}
	}

	if userID == "" {
		writeError(w, http.StatusBadRequest, "userID parameter is required")
		return
	}

	followers, err := h.followUseCase.GetFollowers(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, FollowersResponse{Followers: followers})
}

// GetFollowing obtiene los usuarios que sigue un usuario
func (h *Handlers) GetFollowing(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extraer userID del path (asumiendo formato /users/{userID}/following)
	path := r.URL.Path
	userID := ""
	if len(path) > len("/api/v1/users/") {
		endIndex := len(path) - len("/following")
		if endIndex > len("/api/v1/users/") {
			userID = path[len("/api/v1/users/"):endIndex]
		}
	}

	if userID == "" {
		writeError(w, http.StatusBadRequest, "userID parameter is required")
		return
	}

	following, err := h.followUseCase.GetFollowing(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, FollowingResponse{Following: following})
}
