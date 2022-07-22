package database

import (
	"github.com/joho/godotenv"
	"github.com/stonear/go-template/helper"
	"testing"
)

func TestNew(t *testing.T) {
	err := godotenv.Load("../.env")
	helper.Panic(err)
	New()
}
