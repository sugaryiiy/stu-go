package common

import (
	"os"
	"time"
)

// Config holds runtime configuration values for the application.
type Config struct {
	Port         string
	MySQLDSN     string
	RedisAddr    string
	RedisPass    string
	RedisDB      int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	JWTSecret    string
	AccessTTL    time.Duration
	RefreshTTL   time.Duration
}

// LoadConfig builds a Config instance from environment variables with sensible defaults.
func LoadConfig() Config {
	return Config{
		Port:         getEnv("PORT", "8080"),
		MySQLDSN:     getEnv("MYSQL_DSN", "root:Gu1106..@tcp(121.36.61.64:3306)/dream?parseTime=true&loc=Local"),
		RedisAddr:    getEnv("REDIS_ADDR", "121.36.61.64:6379"),
		RedisPass:    os.Getenv("REDIS_PASSWORD"),
		RedisDB:      0,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		JWTSecret:    getEnv("JWT_SECRET", "dev-secret"),
		AccessTTL:    getEnvDuration("ACCESS_TOKEN_TTL", 15*time.Minute),
		RefreshTTL:   getEnvDuration("REFRESH_TOKEN_TTL", 24*time.Hour),
	}
}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	if val := os.Getenv(key); val != "" {
		if d, err := time.ParseDuration(val); err == nil {
			return d
		}
	}
	return fallback
}
