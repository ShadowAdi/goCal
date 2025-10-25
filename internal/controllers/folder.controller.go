package controllers

import (
	"goCal/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FolderController struct {
	FileService   *services.FileService
	UserService   *services.UserService
	FolderService *services.FolderService
}

func NewFolderController(folderService *services.FolderService, userService *services.UserService, fileService *services.FileService) *FolderController {
	return &FolderController{
		FileService:   fileService,
		UserService:   userService,
		FolderService: folderService,
	}
}

func (fo *FolderController) GetAllFolders(ctx *gin.Context) {
	folders, err := fo.FolderService.GetFolders()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"folders": folders,
	})
	return
}
