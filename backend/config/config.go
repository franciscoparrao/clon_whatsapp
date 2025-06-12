package config

import (
	"os"
	"strconv"
)

type Config struct {
	ServerPort      string
	DatabaseURL     string
	JWTSecret       string
	RedisURL        string
	AllowedOrigins  []string
	MaxMessageSize  int64
	UploadDir       string
}

func LoadConfig() *Config {
	return &Config{
		ServerPort:      getEnv("PORT", "5000"),
		DatabaseURL:     getEnv("DATABASE_URL", "postgres://user:password@localhost/whatsapp_clone"),
		JWTSecret:       getEnv("JWT_SECRET", "your-secret-key-whatsapp-clone-2024"),
		RedisURL:        getEnv("REDIS_URL", "redis://localhost:6379"),
		AllowedOrigins:  []string{"http://localhost:3000"},
		MaxMessageSize:  getEnvAsInt("MAX_MESSAGE_SIZE", 1024*1024), // 1MB
		UploadDir:       getEnv("UPLOAD_DIR", "./uploads"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int64) int64 {
	valueStr := os.Getenv(key)
	if value, err := strconv.ParseInt(valueStr, 10, 64); err == nil {
		return value
	}
	return defaultValue
}