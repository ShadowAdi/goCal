package config

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var mainRouter *gin.Engine

func InitRouter() *gin.Engine {
	mainRouter = gin.Default()
	healthRouter := mainRouter.Group("/health")
	{
		healthRouter.GET("/", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{
				"message": "Server is working",

				"success": true,
			})
		})
	}

	userRouter := mainRouter.Group("/user")
	{
		userRouter.GET("/", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{
				"message": "User is coming",
				"success": true,
			})
		})
	}

	return mainRouter
}
