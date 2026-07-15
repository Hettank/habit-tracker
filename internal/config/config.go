package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all application configuration loaded from environment variables.
type Config struct {
	Port       string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	JWTSecret       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

// Load reads environment variables from a .env file and returns a populated Config.
func Load() *Config {
	err := godotenv.Load()
	accessTTL, err := time.ParseDuration(os.Getenv("ACCESS_TOKEN_TTL"))
	refreshTTL, err := time.ParseDuration(os.Getenv("REFRESH_TOKEN_TTL"))

	if err != nil {
		log.Fatal("invalid ACCESS_TOKEN_TTL")
	}

	if err != nil {
		log.Fatal("invalid REFRESH_TOKEN_TTL")
	}

	if err != nil {
		log.Println(".env file not found")
	}

	return &Config{
		Port: os.Getenv("APP_PORT"),

		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		DBSSLMode:  os.Getenv("DB_SSLMODE"),

		JWTSecret:       os.Getenv("JWT_SECRET"),
		AccessTokenTTL:  accessTTL,
		RefreshTokenTTL: refreshTTL,
	}
}
