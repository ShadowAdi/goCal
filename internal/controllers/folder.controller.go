package controllers

import (
	"goCal/internal/logger"
	"goCal/internal/schema"
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
		logger.Error("Failed to get all folders %v ", err)
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

func (fo *FolderController) GetFolder(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		logger.Error("Folder Id Not Provided")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Folder Id Not Provided",
		})
	}
	folderFound, err := fo.FolderService.GetFolder(id)
	if err != nil {
		logger.Error("Failed to get the folder\n")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err,
			"message": "Failed to get the folder",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"folder":  folderFound,
	})
	return
}

func (fo *FolderController) CreateFolder(ctx *gin.Context) {
	userId, exists := ctx.Get("userId")
	if !exists {
		logger.Error("Error getting userId")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Not Authorized",
		})
	}

	userIdStr, ok := userId.(string)
	if !ok || userIdStr == "" {
		logger.Error("Error parsing userId\n")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Invalid User Id type in context",
		})
		return
	}

	var newFolder *schema.Folder
	if err := ctx.ShouldBindJSON(&newFolder); err != nil {
		logger.Error("Error parsing body: %v\n", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	loggedInUserFound, loggedInUserError := fo.UserService.GetUser(userIdStr)
	if loggedInUserError != nil {
		logger.Error("Error finding logged-in user: %v\n", loggedInUserError)
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   loggedInUserError.Error(),
		})
		return
	}

	if loggedInUserFound.IsVerified == false {
		ctx.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "User Is Not Verified",
		})
	}

	folder, error := fo.FolderService.CreateFolder(newFolder, userIdStr)

	if error != nil {
		logger.Error("Error creating folder %v\n", error.Error())
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   error.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Folder Created",
		"file":    folder,
	})

	return

}

func (fo *FolderController) DeleteFolder(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		logger.Error("Failed to get the id in the request ", id)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Failed to get the id of the request",
		})
	}

	userId, exists := ctx.Get("userId")
	if !exists {
		logger.Error("Error getting userId")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Not Authorized",
		})
	}

	userIdStr, ok := userId.(string)
	if !ok || userIdStr == "" {
		logger.Error("Error parsing userId\n")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Invalid User Id type in context",
		})
		return
	}

	loggedInUserFound, loggedInUserError := fo.UserService.GetUser(userIdStr)
	if loggedInUserError != nil {
		logger.Error("Error finding logged-in user: %v\n", loggedInUserError)
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   loggedInUserError.Error(),
		})
		return
	}

	if loggedInUserFound.IsVerified == false {
		ctx.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "User Is Not Verified",
		})
	}

	message, error := fo.FolderService.DeleteFolder(id, userIdStr)

	if error != nil {
		logger.Error("Error deleting folder %v\n", error.Error())
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   error.Error(),
			"message": message,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": message,
	})
	return

}
