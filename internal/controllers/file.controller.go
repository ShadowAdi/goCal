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

	_, fileError := ctx.FormFile("file")
	if fileError != nil {
		logger.Error("Failed to get the file %v", fileError)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Failed to get the file %s " + fileError.Error(),
		})
	}

	var newFile *schema.File
	if err := ctx.ShouldBindJSON(&newFile); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	loggedInUserFound, loggedInUserError := fc.UserService.GetUser(userIdStr)
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

	return
}

func (fc *FileController) DeleteFile(ctx *gin.Context) {
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

	loggedInUserFound, loggedInUserError := fc.UserService.GetUser(userIdStr)
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

	message, err := fc.FileService.DeleteFile(id, userIdStr)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   err.Error(),
			"message": message,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "File Deleted Successfully",
	})

	return

}

func (fc *FileController) UpdateFile(ctx *gin.Context) {
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

	loggedInUserFound, loggedInUserError := fc.UserService.GetUser(userIdStr)
	if loggedInUserError != nil {
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

	var updateRequest *schema.UpdateFileRequest
	if err := ctx.ShouldBindJSON(&updateRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	updateFile, updateFileError := fc.FileService.UpdateFile(id, userIdStr, updateRequest)
	if updateFileError != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   updateFileError.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "File Updated Successfully",
		"file":    updateFile,
	})

	return

}
