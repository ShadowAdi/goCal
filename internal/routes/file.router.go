package routes

import (
	"goCal/internal/controllers"
	"goCal/internal/middleware"
	"goCal/internal/services"

	"github.com/gin-gonic/gin"
)

func FileRoutes(router *gin.RouterGroup) {
	fileService := services.NewFileService()
	userService := services.NewUserService()
	fileController := controllers.NewFileController(fileService, userService)
	router.GET("/", fileController.GetAllFiles)
	router.GET("/:id", fileController.GetFile)
	protectedRoutes := router.Group("/")
	protectedRoutes.Use(middleware.AuthMiddleware())

	router.POST("/", fileController.CreateFile)
	router.DELETE("/file/:id", fileController.DeleteFile)
	router.PATCH("/file/:id", fileController.UpdateFile)
}
