package controllers

import (
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
