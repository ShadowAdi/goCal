package routes

import (
	"goCal/internal/controllers"
	"goCal/internal/services"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.RouterGroup) {
	userService := services.NewUserService()
	userController := controllers.NewUserController(userService)

	router.GET("/", userController.GetUsers)
	router.GET("/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		userController.GetUser(id, ctx)
	})
	router.POST("/", func(ctx *gin.Context) {
		userController.CreateUser(ctx)
	})
	router.PATCH("/", func(ctx *gin.Context) {
		userController.CreateUser(ctx)
	})
	router.DELETE("/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		userController.DeleteUser(id, ctx)
	})

}
