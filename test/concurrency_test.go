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

	// Function that creates tweets concurrently
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

	// Execute concurrent goroutines
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go createTweets("user1", i*tweetsPerGoroutine)
	}

	wg.Wait()

	// Verify that all tweets were created
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

	// Test 1: Multiple goroutines trying to follow the same relationship
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

	// There should be exactly 1 successful follow
	if successfulFollows != 1 {
		t.Errorf("Expected exactly 1 successful follow, got %d", successfulFollows)
	}

	// Verify that they are indeed following
	isFollowing, err := followUseCase.IsFollowing(ctx, "user1", "user2")
	if err != nil {
		t.Fatalf("Error checking follow status: %v", err)
	}
	if !isFollowing {
		t.Error("Expected user1 to be following user2")
	}

	// Test 2: Multiple goroutines trying to unfollow
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

	// There should be exactly 1 successful unfollow
	if successfulUnfollows != 1 {
		t.Errorf("Expected exactly 1 successful unfollow, got %d", successfulUnfollows)
	}

	// Verify that they are indeed no longer following
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

	// Create some tweets
	for i := 0; i < 10; i++ {
		content := fmt.Sprintf("Test tweet %d", i)
		_, err := tweetUseCase.CreateTweet(ctx, "user1", content)
		if err != nil {
			t.Fatalf("Error creating tweet: %v", err)
		}
	}

	// Follow user1 from user2
	err := followUseCase.FollowUser(ctx, "user2", "user1")
	if err != nil {
		t.Fatalf("Error following user: %v", err)
	}

	const numGoroutines = 100
	var wg sync.WaitGroup

	// Function that reads timeline concurrently
	readTimeline := func(userID string) {
		defer wg.Done()
		_, err := tweetUseCase.GetTimeline(ctx, userID, 50)
		if err != nil {
			t.Errorf("Error getting timeline: %v", err)
		}
	}

	// Execute concurrent reads
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go readTimeline("user2")
	}

	wg.Wait()
}
