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
	router.GET("/:id", userController.GetUser)
	router.POST("/", userController.CreateUser)
	router.PATCH("/:id", userController.UpdateUser)
	router.DELETE("/:id", userController.DeleteUser)

}
