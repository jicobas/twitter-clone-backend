package usecases

import (
	"context"
	"twitter-clone-backend/internal/domain"
	"twitter-clone-backend/internal/ports"
)

// TweetUseCase maneja la lógica de negocio relacionada con tweets
type TweetUseCase struct {
	tweetRepo  ports.TweetRepository
	followRepo ports.FollowRepository
	userRepo   ports.UserRepository
	cache      ports.CacheService
	logger     ports.Logger
}

// NewTweetUseCase crea una nueva instancia del caso de uso
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

// CreateTweet crea un nuevo tweet
func (uc *TweetUseCase) CreateTweet(ctx context.Context, userID, content string) (*domain.Tweet, error) {
	// Verificar que el usuario existe
	exists, err := uc.userRepo.Exists(ctx, userID)
	if err != nil {
		uc.logger.Error("failed to check user existence", err, "userID", userID)
		return nil, err
	}
	if !exists {
		return nil, domain.ErrUserNotFound
	}

	// Crear el tweet
	tweet, err := domain.NewTweet(userID, content)
	if err != nil {
		return nil, err
	}

	// Persistir el tweet
	if err := uc.tweetRepo.Create(ctx, tweet); err != nil {
		uc.logger.Error("failed to create tweet", err, "tweetID", tweet.ID)
		return nil, err
	}

	// Invalidar cache de timeline de seguidores
	if uc.cache != nil {
		uc.invalidateFollowersTimeline(ctx, userID)
	}

	uc.logger.Info("tweet created successfully", "tweetID", tweet.ID, "userID", userID)
	return tweet, nil
}

// GetTimeline obtiene el timeline de un usuario
func (uc *TweetUseCase) GetTimeline(ctx context.Context, userID string, limit int) ([]*domain.Tweet, error) {
	if limit <= 0 || limit > domain.MaxTimelineLimit {
		limit = domain.MaxTimelineLimit
	}

	// Intentar obtener del cache primero
	if uc.cache != nil {
		tweets, err := uc.cache.GetTimeline(ctx, userID)
		if err == nil && tweets != nil {
			uc.logger.Debug("timeline served from cache", "userID", userID)
			return tweets, nil
		}
	}

	// Obtener usuarios que sigue
	following, err := uc.followRepo.GetFollowing(ctx, userID)
	if err != nil {
		uc.logger.Error("failed to get following users", err, "userID", userID)
		return nil, err
	}

	// Incluir tweets del propio usuario
	following = append(following, userID)

	// Obtener tweets del timeline
	tweets, err := uc.tweetRepo.GetTimeline(ctx, following, limit)
	if err != nil {
		uc.logger.Error("failed to get timeline", err, "userID", userID)
		return nil, err
	}

	// Guardar en cache
	if uc.cache != nil {
		if err := uc.cache.SetTimeline(ctx, userID, tweets); err != nil {
			uc.logger.Warn("failed to cache timeline", "error", err, "userID", userID)
		}
	}

	uc.logger.Info("timeline retrieved", "userID", userID, "tweetsCount", len(tweets))
	return tweets, nil
}

// GetUserTweets obtiene todos los tweets de un usuario específico
func (uc *TweetUseCase) GetUserTweets(ctx context.Context, userID string) ([]*domain.Tweet, error) {
	tweets, err := uc.tweetRepo.GetByUserID(ctx, userID)
	if err != nil {
		uc.logger.Error("failed to get user tweets", err, "userID", userID)
		return nil, err
	}

	uc.logger.Info("user tweets retrieved", "userID", userID, "tweetsCount", len(tweets))
	return tweets, nil
}

// invalidateFollowersTimeline invalida el cache de timeline de los seguidores
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
