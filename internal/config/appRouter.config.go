package config

import (
	"goCal/internal/routes"

	"github.com/gin-gonic/gin"
)

var mainRouter *gin.Engine

func InitRouter() *gin.Engine {
	mainRouter = gin.Default()

	healthRouter := mainRouter.Group("/api/health")
	routes.RegisterHealthRoute(healthRouter)

	userRouter := mainRouter.Group("/api/user")
	routes.UserRoutes(userRouter)

	return mainRouter
}
