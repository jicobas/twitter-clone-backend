package ports

import (
	"context"
	"twitter-clone-backend/internal/domain"
)

// CacheService defines cache operations
type CacheService interface {
	GetTimeline(ctx context.Context, userID string) ([]*domain.Tweet, error)
	SetTimeline(ctx context.Context, userID string, tweets []*domain.Tweet) error
	InvalidateTimeline(ctx context.Context, userID string) error
}

// Logger defines logging operations
type Logger interface {
	Info(msg string, args ...interface{})
	Error(msg string, err error, args ...interface{})
	Debug(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
}
