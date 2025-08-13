package ports

import (
	"context"
	"twitter-clone-backend/internal/domain"
)

// TweetRepository define las operaciones para tweets
type TweetRepository interface {
	Create(ctx context.Context, tweet *domain.Tweet) error
	GetByID(ctx context.Context, id string) (*domain.Tweet, error)
	GetByUserID(ctx context.Context, userID string) ([]*domain.Tweet, error)
	GetTimeline(ctx context.Context, userIDs []string, limit int) ([]*domain.Tweet, error)
	Delete(ctx context.Context, id string) error
}

// FollowRepository define las operaciones relacionadas con seguimientos
type FollowRepository interface {
	Follow(ctx context.Context, followerID, followeeID string) error
	FollowIfNotExists(ctx context.Context, followerID, followeeID string) error
	Unfollow(ctx context.Context, followerID, followeeID string) error
	UnfollowIfExists(ctx context.Context, followerID, followeeID string) error
	GetFollowers(ctx context.Context, userID string) ([]string, error)
	GetFollowing(ctx context.Context, userID string) ([]string, error)
	IsFollowing(ctx context.Context, followerID, followeeID string) (bool, error)
}

// UserRepository define las operaciones para usuarios
type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) error
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
	GetUserByUsername(ctx context.Context, username string) (*domain.User, error)
	Exists(ctx context.Context, id string) (bool, error)
}
