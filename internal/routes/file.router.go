package routes

import (
	"goCal/internal/controllers"
	"goCal/internal/middleware"
	"goCal/internal/services"

	"github.com/gin-gonic/gin"
	storage_go "github.com/supabase-community/storage-go"
)

func FileRoutes(router *gin.RouterGroup, storageClient *storage_go.Client) {
	fileService := services.NewFileService()
	userService := services.NewUserService()
	newFileStorageService := services.NewFileStorageService(storageClient)
	fileController := controllers.NewFileController(fileService, userService, newFileStorageService)
	router.GET("/", fileController.GetAllFiles)
	router.GET("/:id", fileController.GetFile)
	protectedRoutes := router.Group("/")
	protectedRoutes.Use(middleware.AuthMiddleware())

	protectedRoutes.POST("/", fileController.CreateFile)
	protectedRoutes.DELETE("/file/:id", fileController.DeleteFile)
	protectedRoutes.PATCH("/file/:id", fileController.UpdateFile)
}
