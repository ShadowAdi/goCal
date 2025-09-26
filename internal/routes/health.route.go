package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterHealthRoute(router *gin.RouterGroup) {
	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Server is working",
			"success": true,
		})
	})
}
