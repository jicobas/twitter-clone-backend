package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	httpAdapters "twitter-clone-backend/internal/adapters/http"
	"twitter-clone-backend/internal/adapters/memory"
	"twitter-clone-backend/internal/usecases"
	"twitter-clone-backend/pkg/logger"
)

// TestAPI runs basic integration tests for the API
func TestAPI(t *testing.T) {
	// Setup
	appLogger := logger.NewLogger()
	repo := memory.NewRepositories()
	tweetUseCase := usecases.NewTweetUseCase(repo, repo, repo, nil, appLogger)
	followUseCase := usecases.NewFollowUseCase(repo, repo, nil, appLogger)
	handlers := httpAdapters.NewHandlers(tweetUseCase, followUseCase)
	router := httpAdapters.SetupRoutes(handlers)

	// Start test server
	server := &http.Server{
		Addr:    ":8081",
		Handler: router,
	}

	go func() {
		server.ListenAndServe()
	}()

	// Wait for server to start
	time.Sleep(100 * time.Millisecond)

	defer server.Close()

	baseURL := "http://localhost:8081"

	// Test 1: Health check
	t.Run("Health check", func(t *testing.T) {
		resp, err := http.Get(baseURL + "/health")
		if err != nil {
			t.Fatalf("Failed to make health check request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})

	// Test 2: Create tweet
	t.Run("Create tweet", func(t *testing.T) {
		tweetData := map[string]string{
			"content": "Hello, World! This is my first tweet.",
		}

		jsonData, _ := json.Marshal(tweetData)

		req, _ := http.NewRequest("POST", baseURL+"/tweets", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-ID", "user1")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to create tweet: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			body, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 201, got %d. Body: %s", resp.StatusCode, body)
		}
	})

	// Test 3: Get timeline
	t.Run("Get timeline", func(t *testing.T) {
		resp, err := http.Get(baseURL + "/users/user1/timeline")
		if err != nil {
			t.Fatalf("Failed to get timeline: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})

	// Test 4: Follow user
	t.Run("Follow user", func(t *testing.T) {
		followData := map[string]string{
			"followee_id": "user2",
		}

		jsonData, _ := json.Marshal(followData)

		req, _ := http.NewRequest("POST", baseURL+"/users/following", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-ID", "user1")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to follow user: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 200, got %d. Body: %s", resp.StatusCode, body)
		}
	})

	// Test 5: Get followers
	t.Run("Get followers", func(t *testing.T) {
		resp, err := http.Get(baseURL + "/users/user2/followers")
		if err != nil {
			t.Fatalf("Failed to get followers: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})
}
