package domain

import "time"

// Follow representa una relación de seguimiento entre usuarios
type Follow struct {
	FollowerID string    `json:"follower_id"` // Usuario que sigue
	FolloweeID string    `json:"followee_id"` // Usuario que es seguido
	CreatedAt  time.Time `json:"created_at"`
}

// NewFollow crea una nueva relación de seguimiento
func NewFollow(followerID, followeeID string) (*Follow, error) {
	if followerID == "" || followeeID == "" {
		return nil, ErrInvalidUserID
	}

	if followerID == followeeID {
		return nil, ErrCannotFollowSelf
	}

	return &Follow{
		FollowerID: followerID,
		FolloweeID: followeeID,
		CreatedAt:  time.Now(),
	}, nil
}

// IsValid verifica si la relación de seguimiento es válida
func (f *Follow) IsValid() bool {
	return f.FollowerID != "" &&
		f.FolloweeID != "" &&
		f.FollowerID != f.FolloweeID
}
