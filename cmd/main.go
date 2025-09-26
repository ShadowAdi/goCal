package main

import (
	"context"
	"goCal/internal/config"
	"goCal/internal/db"
	"goCal/internal/logger"
	"goCal/internal/migrations"
)

func main() {
	logger.InitLogger()
	config.GetLoadEnvVars()
	db.DBConnect()
	defer func() {
		if db.Conn != nil {
			db.Conn.Close(context.Background())
			logger.Info("Database connection closed")
		}
	}()
	migrations.CreateUserTable()

	r := config.InitRouter()
	r.Run(":8080")
}
