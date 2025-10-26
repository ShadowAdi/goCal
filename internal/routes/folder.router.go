package routes

import (
	"goCal/internal/controllers"
	"goCal/internal/services"

	"github.com/gin-gonic/gin"
)

func FolderRoutes(router *gin.RouterGroup) {
	fileService := services.NewFileService()
	userService := services.NewUserService()
	folderService := services.NewFolderService()
	folderController := controllers.NewFolderController(folderService, userService, fileService)

	router.GET("/", folderController.GetAllFolders)

}
