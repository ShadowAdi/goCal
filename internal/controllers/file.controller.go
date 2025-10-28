package controllers

import (
	"fmt"
	"goCal/internal/logger"
	"goCal/internal/schema"
	"goCal/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type FileController struct {
	FileService        *services.FileService
	UserService        *services.UserService
	FileStorageService *services.FileStorageService
}

func NewFileController(fileService *services.FileService, userService *services.UserService, fileStorageService *services.FileStorageService) *FileController {
	return &FileController{
		FileService:        fileService,
		UserService:        userService,
		FileStorageService: fileStorageService,
	}
}

func (fc *FileController) GetAllFiles(ctx *gin.Context) {
	files, err := fc.FileService.GetFiles()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
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
		return
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
		return
	}

	if errParseForm := ctx.Request.ParseMultipartForm(100 << 20); errParseForm != nil {
		logger.Error("Failed to parse multipart form: " + errParseForm.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to parse form data",
			"success": false,
		})
		return
	}

	form, errFileUpload := ctx.MultipartForm()
	if errFileUpload != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid form data: " + errFileUpload.Error(),
		})
		return
	}

	files := form.File["files"]

	if len(files) == 0 {
		singleFile, err := ctx.FormFile("file")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "No file uploaded",
			})
			return
		}
		files = append(files, singleFile)
	}

	var createdFiles []schema.File
	var uploadErrors []string

	for _, fileHeader := range files {
		src, errFileOpen := fileHeader.Open()
		if errFileOpen != nil {
			logger.Error("Failed to open uploaded file: %v", errFileOpen)
			continue
		}
		fileType := fileHeader.Header.Get("Content-Type")

		fileUrl, uploadError := fc.FileStorageService.UploadFile(userIdStr, fileHeader.Filename, src, fileType)
		src.Close()

		if uploadError != nil {
			logger.Error("Failed to upload file %s: %v", fileHeader.Filename, uploadError)
			uploadErrors = append(uploadErrors, fmt.Sprintf("Failed to upload %s: %v", fileHeader.Filename, uploadError))
			continue
		}

		newFile := &schema.File{
			FileName:     fileHeader.Filename,
			FileUrl:      fileUrl,
			FileSize:     fileHeader.Size,
			FileType:     fileType,
			UploadedById: uuid.MustParse(userIdStr),
		}

		createdFile, errFileCreate := fc.FileService.CreateFile(newFile, userIdStr)
		if errFileCreate != nil {
			logger.Error("Error creating file record for %s: %v\n", fileHeader.Filename, errFileCreate)
			uploadErrors = append(uploadErrors, fmt.Sprintf("Failed to save %s to database: %v", fileHeader.Filename, errFileCreate))
			// TODO: Consider deleting the uploaded file from storage here
			continue
		}

		createdFiles = append(createdFiles, *createdFile)
	}

	if len(createdFiles) == 0 {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to upload any files",
			"details": uploadErrors,
		})
		return
	}

	response := gin.H{
		"success": true,
		"message": fmt.Sprintf("Successfully uploaded %d file(s)", len(createdFiles)),
		"files":   createdFiles,
	}

	if len(uploadErrors) > 0 {
		response["partial_errors"] = uploadErrors
		response["message"] = fmt.Sprintf("Uploaded %d file(s) with %d error(s)", len(createdFiles), len(uploadErrors))
	}

	ctx.JSON(http.StatusOK, response)
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
