package domain

import "time"

// User represents a user in the system
type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

// NewUser creates a new user
func NewUser(id, username string) *User {
	return &User{
		ID:        id,
		Username:  username,
		CreatedAt: time.Now(),
	}
}

// IsValid verifies if the user is valid
func (u *User) IsValid() bool {
	return u.ID != "" && u.Username != ""
}
