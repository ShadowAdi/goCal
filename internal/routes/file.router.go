package routes

import (
	"goCal/internal/controllers"
	"goCal/internal/services"

	"github.com/gin-gonic/gin"
)

func FileRoutes(router *gin.RouterGroup) {
	fileService := services.NewFileService()
	fileController := controllers.NewFileController(fileService)
	router.GET("/", fileController.GetAllFiles)
	router.GET("/:id", fileController.GetFile)

}
