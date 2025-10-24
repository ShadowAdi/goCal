package routes

import (
	"goCal/internal/controllers"
	"goCal/internal/services"

	"github.com/gin-gonic/gin"
)

func FileRoutes(router *gin.RouterGroup) {
	fileService := services.NewFileService()
	userService := services.NewUserService()
	fileController := controllers.NewFileController(fileService, userService)
	router.GET("/", fileController.GetAllFiles)
	router.GET("/:id", fileController.GetFile)
	router.POST("/", fileController.CreateFile)
	router.DELETE("/:id", fileController.DeleteFile)

}
