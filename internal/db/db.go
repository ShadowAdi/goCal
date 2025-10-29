package db

import (
	"fmt"
	"goCal/internal/logger"
	"goCal/internal/schema"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func DBConnect() {
	DATABASE_URL := os.Getenv("DATABASE_URL")
	if DATABASE_URL == "" {
		logger.Error(`Failed to get the database url`)
		fmt.Printf(`Failed to get the database url`)
	}

	db, err := gorm.Open(postgres.Open(DATABASE_URL), &gorm.Config{})
	if err != nil {
		fmt.Println("Failed to connect to DB: %w", err)
		logger.Error("Failed to connect to DB: %w", err)
	}

	sqlDb, err := db.DB()
	if err != nil {
		logger.Error("Failed to get sql.DB: %w", err)
		panic(fmt.Errorf("Failed to get sql.DB: %w", err))
	}

	sqlDb.SetMaxOpenConns(25)
	sqlDb.SetMaxIdleConns(5)
	sqlDb.SetConnMaxLifetime(5 * time.Minute)
	sqlDb.SetConnMaxIdleTime(1 * time.Minute)

	if err := sqlDb.Ping(); err != nil {
		logger.Error("Failed to ping DB: %w", err)
		panic(fmt.Errorf("Failed to ping DB: %w", err))
	}

	DB = db

	if err := DB.AutoMigrate(&schema.User{}); err != nil {
		logger.Error("Failed to auto-migrate tables: %w", err)
		panic(fmt.Errorf("Failed to auto-migrate tables: %w", err))
	}

	fmt.Println("Connection established")
	logger.Info("Database connected")
}
