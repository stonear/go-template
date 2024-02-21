package config

import (
	"github.com/joho/godotenv"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
)

func Load(log *otelzap.Logger) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
