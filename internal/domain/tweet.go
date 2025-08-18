package domain

import (
	"time"

	"github.com/google/uuid"
)

// generateID generates a unique ID using UUID v4
func generateID() string {
	return uuid.New().String()
}

// Tweet represents a tweet in the system
type Tweet struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// NewTweet creates a new tweet with validations
func NewTweet(userID, content string) (*Tweet, error) {
	if userID == "" {
		return nil, ErrInvalidUserID
	}

	if content == "" {
		return nil, ErrEmptyContent
	}

	if len(content) > MaxTweetLength {
		return nil, ErrContentTooLong
	}

	return &Tweet{
		ID:        generateID(),
		UserID:    userID,
		Content:   content,
		CreatedAt: time.Now(),
	}, nil
}

// IsValid verifies if the tweet is valid
func (t *Tweet) IsValid() bool {
	return t.ID != "" &&
		t.UserID != "" &&
		t.Content != "" &&
		len(t.Content) <= MaxTweetLength
}
