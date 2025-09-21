package db

import (
	"fmt"
	"goCal/internal/logger"
	"os"
)

func DBConnect() {
	DATABASE_URL := os.Getenv("DATABASE_URL")
	if DATABASE_URL == "" {
		logger.Error(`Failed to get the database url`)
		fmt.Printf(`Failed to get the database url`)
	}
}
