package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port         string
	MaxImageSize string
}

func GetConfig() Config {
	err := godotenv.Load()
	if err != nil {
		slog.Error("Warning: .env file not found or could not be loaded. Relying on system environment variables: ", "error", err)
	} else {
		slog.Error(".env file loaded successfully.")
	}
	portStr := os.Getenv("PORT")
	if portStr == "" {
		slog.Warn("No port found Defaulting to 8080")
		portStr = "8080"
	}

	maxImageSizeStr := os.Getenv("MAXIMAGESIZE")
	if maxImageSizeStr == "" {
		slog.Warn("No MAXIMAGESIZE found Defaulting to 10 MB")
		maxImageSizeStr = "10"
	}

	Conf := Config{Port: portStr, MaxImageSize: maxImageSizeStr}

	return Conf
}
