package config

import (
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func Load(log *zap.Logger) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
