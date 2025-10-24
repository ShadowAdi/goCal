package controllers

import (
	"goCal/internal/logger"
	"goCal/internal/schema"
	"goCal/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FileController struct {
	FileService *services.FileService
	UserService *services.UserService
}

func NewFileController(fileService *services.FileService, userService *services.UserService) *FileController {
	return &FileController{
		FileService: fileService,
		UserService: userService,
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

func (fc *FileController) CreateFile(ctx *gin.Context) {
	userId, exists := ctx.Get("userId")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Not Authorized",
		})
	}

	userIdStr, ok := userId.(string)
	if !ok || userIdStr == "" {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Invalid User Id type in context",
		})
		return
	}

	var newFile *schema.File
	if err := ctx.ShouldBindJSON(&newFile); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	_, loggedInUserError := fc.UserService.GetUser(userIdStr)
	if loggedInUserError != nil {
		logger.Error("Error finding logged-in user: %v\n", loggedInUserError)
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   loggedInUserError.Error(),
		})
		return
	}

	file, error := fc.FileService.CreateFile(newFile, userIdStr)

	if error != nil {
		logger.Error("Error creating file %v\n", error.Error())
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   error.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "File Created",
		"file":    file,
	})

}
