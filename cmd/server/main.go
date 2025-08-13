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
	// Cargar configuración
	cfg := config.LoadConfig()

	// Inicializar logger
	appLogger := logger.NewLogger()
	appLogger.Info("Starting Twitter Clone Backend", "port", cfg.Port, "storage", cfg.StorageType)

	// Inicializar repositorios
	repo := memory.NewRepositories()

	// Inicializar casos de uso
	tweetUseCase := usecases.NewTweetUseCase(repo, repo, repo, nil, appLogger)
	followUseCase := usecases.NewFollowUseCase(repo, repo, nil, appLogger)

	// Inicializar handlers HTTP
	handlers := httpAdapters.NewHandlers(tweetUseCase, followUseCase)

	// Configurar rutas
	router := httpAdapters.SetupRoutes(handlers)

	// Iniciar servidor
	appLogger.Info("Server starting", "port", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
