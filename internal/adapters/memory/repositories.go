package memory

import (
	"context"
	"sort"
	"sync"
	"twitter-clone-backend/internal/domain"
)

// Repositories implementa los repositorios en memoria
type Repositories struct {
	tweets  map[string]*domain.Tweet
	users   map[string]*domain.User
	follows map[string]map[string]bool // followerID -> followeeID -> true
	mu      sync.RWMutex
}

// NewRepositories crea una nueva instancia de repositorios en memoria
func NewRepositories() *Repositories {
	repo := &Repositories{
		tweets:  make(map[string]*domain.Tweet),
		users:   make(map[string]*domain.User),
		follows: make(map[string]map[string]bool),
	}

	// Agregar algunos usuarios de ejemplo para testing
	repo.seedUsers()

	return repo
}

// seedUsers agrega usuarios de ejemplo
func (r *Repositories) seedUsers() {
	users := []*domain.User{
		domain.NewUser("user1", "alice"),
		domain.NewUser("user2", "bob"),
		domain.NewUser("user3", "charlie"),
	}

	for _, user := range users {
		r.users[user.ID] = user
	}
}

// TweetRepository methods

func (r *Repositories) Create(ctx context.Context, tweet *domain.Tweet) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.tweets[tweet.ID] = tweet
	return nil
}

func (r *Repositories) GetByID(ctx context.Context, id string) (*domain.Tweet, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tweet, exists := r.tweets[id]
	if !exists {
		return nil, domain.ErrTweetNotFound
	}

	return tweet, nil
}

func (r *Repositories) GetByUserID(ctx context.Context, userID string) ([]*domain.Tweet, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var tweets []*domain.Tweet
	for _, tweet := range r.tweets {
		if tweet.UserID == userID {
			tweets = append(tweets, tweet)
		}
	}

	// Ordenar por fecha de creación (más reciente primero)
	sort.Slice(tweets, func(i, j int) bool {
		return tweets[i].CreatedAt.After(tweets[j].CreatedAt)
	})

	return tweets, nil
}

func (r *Repositories) GetTimeline(ctx context.Context, userIDs []string, limit int) ([]*domain.Tweet, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	userIDMap := make(map[string]bool)
	for _, id := range userIDs {
		userIDMap[id] = true
	}

	var tweets []*domain.Tweet
	for _, tweet := range r.tweets {
		if userIDMap[tweet.UserID] {
			tweets = append(tweets, tweet)
		}
	}

	// Ordenar por fecha de creación (más reciente primero)
	sort.Slice(tweets, func(i, j int) bool {
		return tweets[i].CreatedAt.After(tweets[j].CreatedAt)
	})

	// Aplicar límite
	if limit > 0 && len(tweets) > limit {
		tweets = tweets[:limit]
	}

	return tweets, nil
}

func (r *Repositories) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tweets[id]; !exists {
		return domain.ErrTweetNotFound
	}

	delete(r.tweets, id)
	return nil
}

// FollowRepository methods

func (r *Repositories) Follow(ctx context.Context, followerID, followeeID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.follows[followerID] == nil {
		r.follows[followerID] = make(map[string]bool)
	}

	r.follows[followerID][followeeID] = true
	return nil
}

func (r *Repositories) FollowIfNotExists(ctx context.Context, followerID, followeeID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Verificar atómicamente si ya está siguiendo
	if r.follows[followerID] != nil && r.follows[followerID][followeeID] {
		return domain.ErrAlreadyFollowing
	}

	// Si no existe, crear la relación
	if r.follows[followerID] == nil {
		r.follows[followerID] = make(map[string]bool)
	}

	r.follows[followerID][followeeID] = true
	return nil
}

func (r *Repositories) Unfollow(ctx context.Context, followerID, followeeID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.follows[followerID] != nil {
		delete(r.follows[followerID], followeeID)
		if len(r.follows[followerID]) == 0 {
			delete(r.follows, followerID)
		}
	}

	return nil
}

func (r *Repositories) UnfollowIfExists(ctx context.Context, followerID, followeeID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Verificar atómicamente si está siguiendo
	if r.follows[followerID] == nil || !r.follows[followerID][followeeID] {
		return domain.ErrNotFollowing
	}

	// Si existe, eliminarlo
	delete(r.follows[followerID], followeeID)
	if len(r.follows[followerID]) == 0 {
		delete(r.follows, followerID)
	}

	return nil
}

func (r *Repositories) GetFollowers(ctx context.Context, userID string) ([]string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var followers []string
	for followerID, following := range r.follows {
		if following[userID] {
			followers = append(followers, followerID)
		}
	}

	return followers, nil
}

func (r *Repositories) GetFollowing(ctx context.Context, userID string) ([]string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var following []string
	if r.follows[userID] != nil {
		for followeeID := range r.follows[userID] {
			following = append(following, followeeID)
		}
	}

	return following, nil
}

func (r *Repositories) IsFollowing(ctx context.Context, followerID, followeeID string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.follows[followerID] == nil {
		return false, nil
	}

	return r.follows[followerID][followeeID], nil
}

// UserRepository methods

func (r *Repositories) CreateUser(ctx context.Context, user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.users[user.ID] = user
	return nil
}

func (r *Repositories) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, domain.ErrUserNotFound
	}

	return user, nil
}

func (r *Repositories) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.Username == username {
			return user, nil
		}
	}

	return nil, domain.ErrUserNotFound
}

func (r *Repositories) Exists(ctx context.Context, id string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.users[id]
	return exists, nil
}
