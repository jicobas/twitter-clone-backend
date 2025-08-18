package usecases

import (
	"context"
	"twitter-clone-backend/internal/domain"
	"twitter-clone-backend/internal/ports"
)

// FollowUseCase handles business logic related to following
type FollowUseCase struct {
	followRepo ports.FollowRepository
	userRepo   ports.UserRepository
	cache      ports.CacheService
	logger     ports.Logger
}

// NewFollowUseCase creates a new instance of the use case
func NewFollowUseCase(
	followRepo ports.FollowRepository,
	userRepo ports.UserRepository,
	cache ports.CacheService,
	logger ports.Logger,
) *FollowUseCase {
	return &FollowUseCase{
		followRepo: followRepo,
		userRepo:   userRepo,
		cache:      cache,
		logger:     logger,
	}
}

// FollowUser allows a user to follow another user
func (uc *FollowUseCase) FollowUser(ctx context.Context, followerID, followeeID string) error {
	// Validate the follow relationship
	follow, err := domain.NewFollow(followerID, followeeID)
	if err != nil {
		return err
	}

	// Verify that both users exist
	followerExists, err := uc.userRepo.Exists(ctx, followerID)
	if err != nil {
		uc.logger.Error("failed to check follower existence", err, "followerID", followerID)
		return err
	}
	if !followerExists {
		return domain.ErrUserNotFound
	}

	followeeExists, err := uc.userRepo.Exists(ctx, followeeID)
	if err != nil {
		uc.logger.Error("failed to check followee existence", err, "followeeID", followeeID)
		return err
	}
	if !followeeExists {
		return domain.ErrUserNotFound
	}

	// Atomic operation: verify + create in a single transaction
	if err := uc.followRepo.FollowIfNotExists(ctx, follow.FollowerID, follow.FolloweeID); err != nil {
		if err == domain.ErrAlreadyFollowing {
			return err
		}
		uc.logger.Error("failed to create follow relationship", err, "followerID", followerID, "followeeID", followeeID)
		return err
	}

	// Invalidate follower's timeline cache asynchronously (not critical path)
	if uc.cache != nil {
		go func() {
			if err := uc.cache.InvalidateTimeline(context.Background(), followerID); err != nil {
				uc.logger.Warn("failed to invalidate timeline cache", "error", err, "followerID", followerID)
			}
		}()
	}

	uc.logger.Info("user followed successfully", "followerID", followerID, "followeeID", followeeID)
	return nil
}

// UnfollowUser allows a user to stop following another user
func (uc *FollowUseCase) UnfollowUser(ctx context.Context, followerID, followeeID string) error {
	if followerID == "" || followeeID == "" {
		return domain.ErrInvalidUserID
	}

	if followerID == followeeID {
		return domain.ErrCannotFollowSelf
	}

	// Atomic operation: verify + delete in a single transaction
	if err := uc.followRepo.UnfollowIfExists(ctx, followerID, followeeID); err != nil {
		if err == domain.ErrNotFollowing {
			return err
		}
		uc.logger.Error("failed to unfollow user", err, "followerID", followerID, "followeeID", followeeID)
		return err
	}

	// Invalidate follower's timeline cache asynchronously (not critical path)
	if uc.cache != nil {
		go func() {
			if err := uc.cache.InvalidateTimeline(context.Background(), followerID); err != nil {
				uc.logger.Warn("failed to invalidate timeline cache", "error", err, "followerID", followerID)
			}
		}()
	}

	uc.logger.Info("user unfollowed successfully", "followerID", followerID, "followeeID", followeeID)
	return nil
}

// GetFollowers gets the list of followers for a user
func (uc *FollowUseCase) GetFollowers(ctx context.Context, userID string) ([]string, error) {
	if userID == "" {
		return nil, domain.ErrInvalidUserID
	}

	followers, err := uc.followRepo.GetFollowers(ctx, userID)
	if err != nil {
		uc.logger.Error("failed to get followers", err, "userID", userID)
		return nil, err
	}

	uc.logger.Info("followers retrieved", "userID", userID, "followersCount", len(followers))
	return followers, nil
}

// GetFollowing gets the list of users that a user follows
func (uc *FollowUseCase) GetFollowing(ctx context.Context, userID string) ([]string, error) {
	if userID == "" {
		return nil, domain.ErrInvalidUserID
	}

	following, err := uc.followRepo.GetFollowing(ctx, userID)
	if err != nil {
		uc.logger.Error("failed to get following", err, "userID", userID)
		return nil, err
	}

	uc.logger.Info("following retrieved", "userID", userID, "followingCount", len(following))
	return following, nil
}

// IsFollowing verifies if a user is following another user
func (uc *FollowUseCase) IsFollowing(ctx context.Context, followerID, followeeID string) (bool, error) {
	if followerID == "" || followeeID == "" {
		return false, domain.ErrInvalidUserID
	}

	isFollowing, err := uc.followRepo.IsFollowing(ctx, followerID, followeeID)
	if err != nil {
		uc.logger.Error("failed to check following status", err, "followerID", followerID, "followeeID", followeeID)
		return false, err
	}

	return isFollowing, nil
}
