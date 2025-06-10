package config

import (
	"log"
	"os"
)

type Config struct {
	MongoURI  string
	MongoDB   string
	Port      string
	JWTSecret string
}

func LoadConfig() *Config {
	cfg := &Config{
		MongoURI: mustGetEnv("MONGO_URI"),
		MongoDB:  mustGetEnv("MONGO_DB_NAME"),
		Port:     getenv("PORT", "8080"),
		JWTSecret: mustGetEnv("JWT_SECRET"),
	}
	return cfg
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	log.Printf("⚠️  Env var %s not set, using fallback '%s'\n", key, fallback)
	return fallback
}

func mustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("❌ Environment variable %s must be set", key)
	}
	return value
}
