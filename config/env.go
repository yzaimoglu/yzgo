package config

import (
	"os"
	"strconv"

	"github.com/gookit/slog"
	"github.com/joho/godotenv"
)

func Load() {
	err := godotenv.Load()
	if err != nil {
		slog.Errorf("Error loading .env file: %v", err)
		slog.Errorf("Using default environment variables instead.")
	}
}

func GetString(key string) string {
	return os.Getenv(key)
}

func GetInteger(key string) int {
	intValue, err := strconv.ParseInt(os.Getenv(key), 10, 64)
	if err != nil {
		slog.Errorf("Could not parse the integer value for key: %v", err)
		return 0
	}

	return int(intValue)
}

// GetBoolean returns a boolean value from the environment
func GetBoolean(key string) bool {
	boolValue, err := strconv.ParseBool(os.Getenv(key))
	if err != nil {
		slog.Errorf("Could not parse the boolean value for key: %v", err)
		return false
	}

	return boolValue
}
