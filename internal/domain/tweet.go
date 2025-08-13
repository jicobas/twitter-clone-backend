package domain

import (
	"crypto/rand"
	"fmt"
	"time"
)

// generateID genera un ID único
func generateID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

// Tweet representa un tweet en el sistema
type Tweet struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// NewTweet crea un nuevo tweet con validaciones
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

// IsValid verifica si el tweet es válido
func (t *Tweet) IsValid() bool {
	return t.ID != "" &&
		t.UserID != "" &&
		t.Content != "" &&
		len(t.Content) <= MaxTweetLength
}
