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
	router.GET("/deleted", userController.GetSoftDeletedUsers) // Get all soft-deleted users
	router.GET("/:id", userController.GetUser)
	router.POST("/", userController.CreateUser)
	router.PATCH("/:id", userController.UpdateUser)
	router.DELETE("/:id", userController.DeleteUser)                      // Soft delete
	router.POST("/:id/restore", userController.RestoreUser)               // Restore soft-deleted user
	router.DELETE("/:id/permanent", userController.PermanentlyDeleteUser) // Hard delete

}
