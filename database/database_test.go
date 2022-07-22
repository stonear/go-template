package database

import (
	"testing"

	"github.com/joho/godotenv"
	"github.com/stonear/go-template/helper"
)

func TestNew(t *testing.T) {
	err := godotenv.Load("../.env")
	helper.Panic(err)
	New()
}
