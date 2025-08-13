package test

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"twitter-clone-backend/internal/adapters/memory"
	"twitter-clone-backend/internal/usecases"
	"twitter-clone-backend/pkg/logger"
)

func TestConcurrentTweetCreation(t *testing.T) {
	// Setup
	repo := memory.NewRepositories()
	logger := logger.NewLogger()
	tweetUseCase := usecases.NewTweetUseCase(repo, repo, repo, nil, logger)

	const numGoroutines = 100
	const tweetsPerGoroutine = 10

	var wg sync.WaitGroup
	ctx := context.Background()

	// Función que crea tweets concurrentemente
	createTweets := func(userID string, startIndex int) {
		defer wg.Done()
		for i := 0; i < tweetsPerGoroutine; i++ {
			content := fmt.Sprintf("Tweet %d from %s", startIndex+i, userID)
			_, err := tweetUseCase.CreateTweet(ctx, userID, content)
			if err != nil {
				t.Errorf("Error creating tweet: %v", err)
			}
		}
	}

	// Ejecutar goroutines concurrentes
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go createTweets("user1", i*tweetsPerGoroutine)
	}

	wg.Wait()

	// Verificar que se crearon todos los tweets
	tweets, err := tweetUseCase.GetUserTweets(ctx, "user1")
	if err != nil {
		t.Fatalf("Error getting user tweets: %v", err)
	}

	expectedCount := numGoroutines * tweetsPerGoroutine
	if len(tweets) != expectedCount {
		t.Errorf("Expected %d tweets, got %d", expectedCount, len(tweets))
	}
}

func TestConcurrentFollowOperations(t *testing.T) {
	// Setup
	repo := memory.NewRepositories()
	logger := logger.NewLogger()
	followUseCase := usecases.NewFollowUseCase(repo, repo, nil, logger)

	const numGoroutines = 50
	ctx := context.Background()

	var wg sync.WaitGroup
	var successfulFollows int32
	var successfulUnfollows int32

	// Test 1: Múltiples goroutines intentando hacer follow a la misma relación
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			err := followUseCase.FollowUser(ctx, "user1", "user2")
			if err == nil {
				atomic.AddInt32(&successfulFollows, 1)
			}
		}()
	}
	wg.Wait()

	// Solo debería haber exactamente 1 follow exitoso
	if successfulFollows != 1 {
		t.Errorf("Expected exactly 1 successful follow, got %d", successfulFollows)
	}

	// Verificar que efectivamente está siguiendo
	isFollowing, err := followUseCase.IsFollowing(ctx, "user1", "user2")
	if err != nil {
		t.Fatalf("Error checking follow status: %v", err)
	}
	if !isFollowing {
		t.Error("Expected user1 to be following user2")
	}

	// Test 2: Múltiples goroutines intentando hacer unfollow
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			err := followUseCase.UnfollowUser(ctx, "user1", "user2")
			if err == nil {
				atomic.AddInt32(&successfulUnfollows, 1)
			}
		}()
	}
	wg.Wait()

	// Solo debería haber exactamente 1 unfollow exitoso
	if successfulUnfollows != 1 {
		t.Errorf("Expected exactly 1 successful unfollow, got %d", successfulUnfollows)
	}

	// Verificar que efectivamente ya no está siguiendo
	isFollowing, err = followUseCase.IsFollowing(ctx, "user1", "user2")
	if err != nil {
		t.Fatalf("Error checking follow status: %v", err)
	}
	if isFollowing {
		t.Error("Expected user1 to not be following user2 after unfollow")
	}
}

func TestConcurrentTimelineAccess(t *testing.T) {
	// Setup
	repo := memory.NewRepositories()
	logger := logger.NewLogger()
	tweetUseCase := usecases.NewTweetUseCase(repo, repo, repo, nil, logger)
	followUseCase := usecases.NewFollowUseCase(repo, repo, nil, logger)

	ctx := context.Background()

	// Crear algunos tweets
	for i := 0; i < 10; i++ {
		content := fmt.Sprintf("Test tweet %d", i)
		_, err := tweetUseCase.CreateTweet(ctx, "user1", content)
		if err != nil {
			t.Fatalf("Error creating tweet: %v", err)
		}
	}

	// Follow user1 desde user2
	err := followUseCase.FollowUser(ctx, "user2", "user1")
	if err != nil {
		t.Fatalf("Error following user: %v", err)
	}

	const numGoroutines = 100
	var wg sync.WaitGroup

	// Función que lee timeline concurrentemente
	readTimeline := func(userID string) {
		defer wg.Done()
		_, err := tweetUseCase.GetTimeline(ctx, userID, 50)
		if err != nil {
			t.Errorf("Error getting timeline: %v", err)
		}
	}

	// Ejecutar lecturas concurrentes
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go readTimeline("user2")
	}

	wg.Wait()
}
