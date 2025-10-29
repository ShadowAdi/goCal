// @title goCal API Documentation
// @version 1.0
// @description API documentation for goCal backend
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com

// @host localhost:8080
// @BasePath /api
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
