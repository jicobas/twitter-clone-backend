package domain

import "time"

// Follow represents a following relationship between users
type Follow struct {
	FollowerID string    `json:"follower_id"` // User who follows
	FolloweeID string    `json:"followee_id"` // User who is followed
	CreatedAt  time.Time `json:"created_at"`
}

// NewFollow creates a new following relationship
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

// IsValid verifies if the following relationship is valid
func (f *Follow) IsValid() bool {
	return f.FollowerID != "" &&
		f.FolloweeID != "" &&
		f.FollowerID != f.FolloweeID
}
