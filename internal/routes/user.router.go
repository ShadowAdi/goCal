package routes

import (
	"goCal/internal/controllers"
	"goCal/internal/middleware"
	"goCal/internal/services"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.RouterGroup) {
	userService := services.NewUserService()
	userController := controllers.NewUserController(userService)

	router.GET("/", userController.GetUsers)
	router.GET("/:id", userController.GetUser)
	router.POST("/", userController.CreateUser)
	router.POST("/login", userController.LoginUser)

	protectedRoutes := router.Group("/")

	protectedRoutes.Use(middleware.AuthMiddleware())

	protectedRoutes.PATCH("/:id", userController.UpdateUser)
	protectedRoutes.DELETE("/:id", userController.DeleteUser)

	protectedRoutes.GET("/deleted", userController.GetSoftDeletedUsers)
	protectedRoutes.POST("/:id/restore", userController.RestoreUser)               // Restore soft-deleted user
	protectedRoutes.DELETE("/:id/permanent", userController.PermanentlyDeleteUser) // Hard delete

}
