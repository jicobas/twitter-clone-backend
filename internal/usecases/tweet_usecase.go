package usecases

import (
	"context"
	"twitter-clone-backend/internal/domain"
	"twitter-clone-backend/internal/ports"
)

// TweetUseCase handles business logic related to tweets
type TweetUseCase struct {
	tweetRepo  ports.TweetRepository
	followRepo ports.FollowRepository
	userRepo   ports.UserRepository
	cache      ports.CacheService
	logger     ports.Logger
}

// NewTweetUseCase creates a new instance of the use case
func NewTweetUseCase(
	tweetRepo ports.TweetRepository,
	followRepo ports.FollowRepository,
	userRepo ports.UserRepository,
	cache ports.CacheService,
	logger ports.Logger,
) *TweetUseCase {
	return &TweetUseCase{
		tweetRepo:  tweetRepo,
		followRepo: followRepo,
		userRepo:   userRepo,
		cache:      cache,
		logger:     logger,
	}
}

// CreateTweet creates a new tweet
func (uc *TweetUseCase) CreateTweet(ctx context.Context, userID, content string) (*domain.Tweet, error) {
	// Verify that the user exists
	exists, err := uc.userRepo.Exists(ctx, userID)
	if err != nil {
		uc.logger.Error("failed to check user existence", err, "userID", userID)
		return nil, err
	}
	if !exists {
		return nil, domain.ErrUserNotFound
	}

	// Create the tweet
	tweet, err := domain.NewTweet(userID, content)
	if err != nil {
		return nil, err
	}

	// Persist the tweet
	if err := uc.tweetRepo.Create(ctx, tweet); err != nil {
		uc.logger.Error("failed to create tweet", err, "tweetID", tweet.ID)
		return nil, err
	}

	// Invalidate followers' timeline cache
	if uc.cache != nil {
		uc.invalidateFollowersTimeline(ctx, userID)
	}

	uc.logger.Info("tweet created successfully", "tweetID", tweet.ID, "userID", userID)
	return tweet, nil
}

// GetTimeline gets a user's timeline
func (uc *TweetUseCase) GetTimeline(ctx context.Context, userID string, limit int) ([]*domain.Tweet, error) {
	if limit <= 0 || limit > domain.MaxTimelineLimit {
		limit = domain.MaxTimelineLimit
	}

	// Try to get from cache first
	if uc.cache != nil {
		tweets, err := uc.cache.GetTimeline(ctx, userID)
		if err == nil && tweets != nil {
			uc.logger.Debug("timeline served from cache", "userID", userID)
			return tweets, nil
		}
	}

	// Get users being followed
	following, err := uc.followRepo.GetFollowing(ctx, userID)
	if err != nil {
		uc.logger.Error("failed to get following users", err, "userID", userID)
		return nil, err
	}

	// Include tweets from the user themselves
	following = append(following, userID)

	// Get timeline tweets
	tweets, err := uc.tweetRepo.GetTimeline(ctx, following, limit)
	if err != nil {
		uc.logger.Error("failed to get timeline", err, "userID", userID)
		return nil, err
	}

	// Save to cache
	if uc.cache != nil {
		if err := uc.cache.SetTimeline(ctx, userID, tweets); err != nil {
			uc.logger.Warn("failed to cache timeline", "error", err, "userID", userID)
		}
	}

	uc.logger.Info("timeline retrieved", "userID", userID, "tweetsCount", len(tweets))
	return tweets, nil
}

// GetUserTweets gets all tweets from a specific user
func (uc *TweetUseCase) GetUserTweets(ctx context.Context, userID string) ([]*domain.Tweet, error) {
	tweets, err := uc.tweetRepo.GetByUserID(ctx, userID)
	if err != nil {
		uc.logger.Error("failed to get user tweets", err, "userID", userID)
		return nil, err
	}

	uc.logger.Info("user tweets retrieved", "userID", userID, "tweetsCount", len(tweets))
	return tweets, nil
}

// invalidateFollowersTimeline invalidates the timeline cache of followers
func (uc *TweetUseCase) invalidateFollowersTimeline(ctx context.Context, userID string) {
	followers, err := uc.followRepo.GetFollowers(ctx, userID)
	if err != nil {
		uc.logger.Warn("failed to get followers for cache invalidation", "error", err, "userID", userID)
		return
	}

	for _, followerID := range followers {
		if err := uc.cache.InvalidateTimeline(ctx, followerID); err != nil {
			uc.logger.Warn("failed to invalidate follower timeline", "error", err, "followerID", followerID)
		}
	}
}
