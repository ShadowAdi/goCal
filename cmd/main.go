package main

import (
	"goCal/internal/config"
	"goCal/internal/db"
	"goCal/internal/logger"
)

func main() {
	logger.InitLogger()
	config.GetLoadEnvVars()
	config.StorageInit()
	db.DBConnect()

	r := config.InitRouter()
	r.Run(":8080")
}
