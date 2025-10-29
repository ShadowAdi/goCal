package controllers

import (
	"fmt"
	"goCal/internal/logger"
	"goCal/internal/schema"
	"goCal/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FolderController struct {
	UserService   *services.UserService
	FolderService *services.FolderService
}

func NewFolderController(userService *services.UserService, folderService *services.FolderService) *FolderController {
	return &FolderController{
		UserService:   userService,
		FolderService: folderService,
	}
}

func (fo FolderController) GetAllFolders(ctx *gin.Context) {
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

func (fo FolderController) GetFolder(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Failed to get the id",
		})
		return
	}

	folders, err := fo.FolderService.GetFolder(id)
	if err != nil {
		logger.Error("Failed to get folder %v ", err)
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

func (fo FolderController) CreateFolder(ctx *gin.Context) {
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

	folder, error := fo.FolderService.CreateFolder(userIdStr, newFolder)

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

func (fo FolderController) UpdateFolder(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Failed to get the id",
		})
		return
	}

	userId, exists := ctx.Get("userId")

	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "User Id not found in context",
		})
		return
	}

	userIdStr, ok := userId.(string)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Invalid User Id type in context",
		})
		return
	}

	loggedInUserFound, loggedInUserError := fo.UserService.GetUser(userIdStr)
	if loggedInUserError != nil {
		fmt.Printf("Error finding logged-in user: %v\n", loggedInUserError)
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

	var existingFolder *schema.Folder
	if err := ctx.ShouldBindJSON(&existingFolder); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Failed to parse the folder",
		})
		return
	}

	if existingFolder == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Failed to get the folder",
		})
	}

}

func (fo FolderController) DeleteFolder(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Failed to get the id",
		})
		return
	}

	userId, exists := ctx.Get("userId")

	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "User Id not found in context",
		})
		return
	}

	userIdStr, ok := userId.(string)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Invalid User Id type in context",
		})
		return
	}

	loggedInUserFound, loggedInUserError := fo.UserService.GetUser(userIdStr)
	if loggedInUserError != nil {
		fmt.Printf("Error finding logged-in user: %v\n", loggedInUserError)
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

	var existingFolder *schema.Folder
	if err := ctx.ShouldBindJSON(&existingFolder); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Failed to parse the folder",
		})
		return
	}

	if existingFolder == nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Failed to get the folder",
		})
	}

}
