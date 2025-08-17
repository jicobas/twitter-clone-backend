package main

import (
	"log"
	"net/http"
	httpAdapters "twitter-clone-backend/internal/adapters/http"
	"twitter-clone-backend/internal/adapters/memory"
	"twitter-clone-backend/internal/config"
	"twitter-clone-backend/internal/usecases"
	"twitter-clone-backend/pkg/logger"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize logger
	appLogger := logger.NewLogger()
	appLogger.Info("Starting Twitter Clone Backend", "port", cfg.Port, "storage", cfg.StorageType)

	// Initialize repositories
	repo := memory.NewRepositories()

	// Initialize use cases
	tweetUseCase := usecases.NewTweetUseCase(repo, repo, repo, nil, appLogger)
	followUseCase := usecases.NewFollowUseCase(repo, repo, nil, appLogger)

	// Initialize HTTP handlers
	handlers := httpAdapters.NewHandlers(tweetUseCase, followUseCase)

	// Configure routes
	router := httpAdapters.SetupRoutes(handlers)

	// Start server
	appLogger.Info("Server starting", "port", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
