package routes

import (
	"goCal/internal/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.RouterGroup) {
	userController := controllers.NewUserController()
	router.GET("/", userController.GetUsers)
	router.GET("/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		userController.GetUser(id, ctx)
	})

}
