package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"twitter-clone-backend/internal/usecases"
)

// Handlers contains HTTP handlers
type Handlers struct {
	tweetUseCase  *usecases.TweetUseCase
	followUseCase *usecases.FollowUseCase
}

// NewHandlers creates a new instance of handlers
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

// writeJSON writes a JSON response
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// writeError writes a JSON error response
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, ErrorResponse{Error: message})
}

// extractUserIDFromPath extracts userID from paths like /users/{userID}/timeline, /users/{userID}/tweets, etc.
func extractUserIDFromPath(path, suffix string) string {
	prefix := "/users/"
	if len(path) <= len(prefix) {
		return ""
	}

	endIndex := len(path) - len(suffix)
	if endIndex <= len(prefix) {
		return ""
	}

	return path[len(prefix):endIndex]
}

// CreateTweet handles the creation of a new tweet
func (h *Handlers) CreateTweet(w http.ResponseWriter, r *http.Request) {
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

// GetTimeline gets a user's timeline
func (h *Handlers) GetTimeline(w http.ResponseWriter, r *http.Request) {
	// Extract userID from path (format: /users/{userID}/timeline)
	userID := extractUserIDFromPath(r.URL.Path, "/timeline")

	if userID == "" {
		writeError(w, http.StatusBadRequest, "userID parameter is required")
		return
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 50 // default value
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

// GetUserTweets gets all tweets from a user
func (h *Handlers) GetUserTweets(w http.ResponseWriter, r *http.Request) {
	// Extract userID from path (format: /users/{userID}/tweets)
	userID := extractUserIDFromPath(r.URL.Path, "/tweets")

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

// FollowUser allows a user to follow another user
func (h *Handlers) FollowUser(w http.ResponseWriter, r *http.Request) {
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

// UnfollowUser allows a user to stop following another user
func (h *Handlers) UnfollowUser(w http.ResponseWriter, r *http.Request) {
	followerID := r.Header.Get("X-User-ID")
	if followerID == "" {
		writeError(w, http.StatusBadRequest, "X-User-ID header is required")
		return
	}

	// Extract followeeID from path (format: /users/following/{followeeID})
	followeeID := r.URL.Path[len("/users/following/"):]
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

// GetFollowers gets the followers of a user
func (h *Handlers) GetFollowers(w http.ResponseWriter, r *http.Request) {
	// Extract userID from path (format: /users/{userID}/followers)
	userID := extractUserIDFromPath(r.URL.Path, "/followers")

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

// GetFollowing gets the users that a user follows
func (h *Handlers) GetFollowing(w http.ResponseWriter, r *http.Request) {
	// Extract userID from path (format: /users/{userID}/following)
	userID := extractUserIDFromPath(r.URL.Path, "/following")

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
