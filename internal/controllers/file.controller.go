package controllers

import (
	"goCal/internal/logger"
	"goCal/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FileController struct {
	FileService *services.FileService
}

func NewFileController(fileService *services.FileService) *FileController {
	return &FileController{
		FileService: fileService,
	}
}

func (fc *FileController) GetAllFiles(ctx *gin.Context) {
	files, err := fc.FileService.GetFiles()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"files":   files,
	})
	return
}

func (fc *FileController) GetFile(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		logger.Error("Failed to get the id in the request ", id)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Failed to get the id of the request",
		})
	}

	file, error := fc.FileService.GetFile(id)

	if error != nil {
		logger.Error("Failed to find the find %s", error.Error())
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Failed to get the file ",
			"error":   error.Error(),
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": file,
	})

}
