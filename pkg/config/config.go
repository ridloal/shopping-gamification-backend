package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost        string
	DBUser        string
	DBPassword    string
	DBName        string
	DBPort        string
	JWTSecret     string
	RedisAddress  string
	RedisUser     string
	RedisPassword string
	RedisDB       string
}

func LoadConfig() (*Config, error) {
	// Load .env file
	godotenv.Load()

	config := &Config{
		DBHost:        os.Getenv("DB_HOST"),
		DBUser:        os.Getenv("DB_USER"),
		DBPassword:    os.Getenv("DB_PASSWORD"),
		DBName:        os.Getenv("DB_NAME"),
		DBPort:        os.Getenv("DB_PORT"),
		JWTSecret:     os.Getenv("JWT_SECRET"),
		RedisAddress:  os.Getenv("REDIS_ADDRESS"),
		RedisUser:     os.Getenv("REDIS_USER"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		RedisDB:       os.Getenv("REDIS_DB"),
	}

	// Validate required environment variables
	required := []struct {
		value, name string
	}{
		{config.DBHost, "DB_HOST"},
		{config.DBUser, "DB_USER"},
		{config.DBPassword, "DB_PASSWORD"},
		{config.DBName, "DB_NAME"},
	}

	for _, r := range required {
		if r.value == "" {
			return nil, fmt.Errorf("%s environment variable is required", r.name)
		}
	}

	return config, nil
}
