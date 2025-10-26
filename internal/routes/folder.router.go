package routes

import (
	"goCal/internal/controllers"
	"goCal/internal/middleware"
	"goCal/internal/services"

	"github.com/gin-gonic/gin"
)

func FolderRoutes(router *gin.RouterGroup) {
	fileService := services.NewFileService()
	userService := services.NewUserService()
	folderService := services.NewFolderService()
	folderController := controllers.NewFolderController(folderService, userService, fileService)

	router.GET("/", folderController.GetAllFolders)
	router.GET("/:id", folderController.GetFolder)

	protectedRoutes := router.Group("/")
	protectedRoutes.Use(middleware.AuthMiddleware())

	protectedRoutes.POST("/", folderController.CreateFolder)
	protectedRoutes.PATCH("/folder/:id", folderController.UpdateFolder)
	protectedRoutes.DELETE("/folder/:id", folderController.DeleteFolder)

}
