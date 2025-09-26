package db

import (
	"context"
	"fmt"
	"goCal/internal/logger"
	"os"

	"github.com/jackc/pgx/v5"
)

var Conn *pgx.Conn

func DBConnect() {
	DATABASE_URL := os.Getenv("DATABASE_URL")
	if DATABASE_URL == "" {
		logger.Error(`Failed to get the database url`)
		fmt.Printf(`Failed to get the database url`)
	}
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, DATABASE_URL)
	if err != nil {
		logger.Error("Failed to connect to the database: " + err.Error())
		fmt.Println("Failed to connect to the database: " + err.Error())
		os.Exit(1)
	}
	var one int
	err = conn.QueryRow(ctx, "SELECT 1").Scan(&one)
	if err != nil {
		logger.Error("Failed to ping the database: " + err.Error())
		fmt.Println("Failed to ping the database:", err)
		os.Exit(1)
	}

	Conn = conn
	fmt.Println("Connection established")
	logger.Info("Database connected")
}
