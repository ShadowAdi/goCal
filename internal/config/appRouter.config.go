package config

import (
	"goCal/internal/routes"

	"github.com/gin-gonic/gin"
)

var mainRouter *gin.Engine

func InitRouter() *gin.Engine {
	mainRouter = gin.Default()

	storageClient := GetStorageClient()

	healthRouter := mainRouter.Group("/api/health")
	routes.RegisterHealthRoute(healthRouter)

	userRouter := mainRouter.Group("/api/user")
	routes.UserRoutes(userRouter)

	fileRouter := mainRouter.Group("/api/file")
	routes.FileRoutes(fileRouter, storageClient)

	return mainRouter
}
