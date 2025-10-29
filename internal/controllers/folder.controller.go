package controllers

import (
	"goCal/internal/schema"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FolderController struct {
}

func NewFolderController() *FolderController {
	return &FolderController{}
}

func (fo FolderController) GetAllFolders(ctx *gin.Context) {

}

func (fo FolderController) GetFolder(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Failed to get the id",
		})
		return
	}

}

func (fo FolderController) CreateFolder(ctx *gin.Context) {
	var newFolder *schema.Folder
	if err := ctx.ShouldBindJSON(&newFolder); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Failed to parse the folder",
		})
		return
	}

}
