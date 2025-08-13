package config

import (
	"os"
	"strconv"
)

// Config contiene toda la configuración de la aplicación
type Config struct {
	Port        string
	StorageType string
	MongoURI    string
	RedisURI    string
	EnableCache bool
}

// LoadConfig carga la configuración desde variables de entorno
func LoadConfig() *Config {
	return &Config{
		Port:        getEnv("PORT", "8080"),
		StorageType: getEnv("STORAGE_TYPE", "memory"),
		MongoURI:    getEnv("MONGO_URI", "mongodb://localhost:27017"),
		RedisURI:    getEnv("REDIS_URI", "redis://localhost:6379"),
		EnableCache: getEnvAsBool("ENABLE_CACHE", false),
	}
}

// getEnv obtiene una variable de entorno con valor por defecto
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsBool obtiene una variable de entorno como boolean
func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if result, err := strconv.ParseBool(value); err == nil {
			return result
		}
	}
	return defaultValue
}
