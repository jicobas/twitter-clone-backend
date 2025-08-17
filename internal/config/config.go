package config

import (
	"os"
	"strconv"
)

// Config contains all application configuration
type Config struct {
	Port        string
	StorageType string
	MongoURI    string
	RedisURI    string
	EnableCache bool
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	return &Config{
		Port:        getEnv("PORT", "8080"),
		StorageType: getEnv("STORAGE_TYPE", "memory"),
		MongoURI:    getEnv("MONGO_URI", "mongodb://localhost:27017"),
		RedisURI:    getEnv("REDIS_URI", "redis://localhost:6379"),
		EnableCache: getEnvAsBool("ENABLE_CACHE", false),
	}
}

// getEnv gets an environment variable with default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsBool gets an environment variable as boolean
func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if result, err := strconv.ParseBool(value); err == nil {
			return result
		}
	}
	return defaultValue
}
