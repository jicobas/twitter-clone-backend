package domain

import "errors"

// Domain errors
var (
	ErrInvalidUserID    = errors.New("invalid user ID")
	ErrEmptyContent     = errors.New("tweet content cannot be empty")
	ErrContentTooLong   = errors.New("tweet content exceeds maximum length")
	ErrUserNotFound     = errors.New("user not found")
	ErrTweetNotFound    = errors.New("tweet not found")
	ErrAlreadyFollowing = errors.New("already following this user")
	ErrNotFollowing     = errors.New("not following this user")
	ErrCannotFollowSelf = errors.New("cannot follow yourself")
)

// Business constants
const (
	MaxTweetLength   = 280
	MaxTimelineLimit = 100
)
