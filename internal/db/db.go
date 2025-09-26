package db

import (
	"context"
	"fmt"
	"goCal/internal/logger"
	"os"

	"github.com/go-pg/pg"
)

func DBConnect() {
	DATABASE_URL := os.Getenv("DATABASE_URL")
	if DATABASE_URL == "" {
		logger.Error(`Failed to get the database url`)
		fmt.Printf(`Failed to get the database url`)
	}
	ctx := context.Background()

	opt, err := pg.ParseURL(DATABASE_URL)
	if err != nil {
		logger.Error("Failed to connect to the database " + err.Error())
		fmt.Printf("Failed to connect to the database %s ", err)
		os.Exit(1)
	}
	db := pg.Connect(opt)

	_, err = db.ExecContext(ctx, "SELECT 1")
	if err != nil {
		logger.Error("Failed to ping the database " + err.Error())
		fmt.Printf("Failed to ping the database %s ", err)
		os.Exit(1)
	}

	fmt.Println("Connection established")
	logger.Info("Database in Connected")
}
