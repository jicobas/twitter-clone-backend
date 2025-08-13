package domain

import "time"

// User representa un usuario en el sistema
type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

// NewUser crea un nuevo usuario
func NewUser(id, username string) *User {
	return &User{
		ID:        id,
		Username:  username,
		CreatedAt: time.Now(),
	}
}

// IsValid verifica si el usuario es válido
func (u *User) IsValid() bool {
	return u.ID != "" && u.Username != ""
}
