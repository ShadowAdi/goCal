package services

import (
	"errors"
	"fmt"
	"goCal/internal/db"
	"goCal/internal/logger"
	"goCal/internal/schema"

	"gorm.io/gorm"
)

type FileService struct {
}

func NewFileService() *FileService {
	return &FileService{}
}

func (f *FileService) GetFiles() ([]*schema.File, error) {
	var files []*schema.File
	result := db.DB.Find(&files)
	if result.Error != nil {
		logger.Error("Failed to get all the files %s ", result.Error)
		return nil, result.Error
	}
	return files, nil
}

func (f *FileService) GetFile(id string) (*schema.File, error) {
	var file *schema.File
	result := db.DB.Find(&file).Where("id = ?", id)
	if result.Error != nil {
		logger.Error("Failed to get  the file %s ", result.Error)
		return nil, result.Error
	}
	return file, nil
}

func (f *FileService) GetUserFile(id string, userId string) (*schema.File, error) {
	var file *schema.File
	result := db.DB.Find(&file).Where("id = ? AND uploaded_by = ?", file.FileName, userId)
	if result.Error != nil {
		logger.Error("Failed to get  the file %s ", result.Error)
		return nil, result.Error
	}
	return file, nil
}

func (f *FileService) CreateFile(file *schema.File, userId string) (*schema.File, error) {
	var existingFile schema.File

	// Check if file with same name already exists for user
	result := db.DB.Where("file_name = ? AND uploaded_by_id = ?", file.FileName, userId).First(&existingFile)

	if result.Error == nil {
		// record found => duplicate
		return nil, fmt.Errorf("file with same name already exists")
	}

	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// real DB error
		return nil, result.Error
	}

	// Create new file
	resultFileCreation := db.DB.Create(file)
	if resultFileCreation.Error != nil {
		logger.Error("Failed to create file %s", resultFileCreation.Error)
		return nil, resultFileCreation.Error
	}

	return file, nil
}

func (f *FileService) DeleteFile(fileId string, userId string) (message string, err error) {
	fileFound, err := f.GetUserFile(fileId, userId)
	if err != nil {
		logger.Error("Failed to get the file  with the fileId %s ", err.Error())
		return "Failed to delete file", err
	}

	if err := db.DB.Delete(fileFound).Error; err != nil {
		logger.Error("Failed to delete the file  with the fileId %s ", err.Error())
		return "Failed to delete file", err
	}
	return "File Deleted Successfully", nil
}

func (f *FileService) UpdateFile(fileId string, userId string, updateFile *schema.UpdateFileRequest) (message *schema.File, err error) {
	_, errFile := f.GetUserFile(fileId, userId)
	if errFile != nil {
		logger.Error("Failed to get the file  with the fileId %s ", err.Error())
		return nil, errFile
	}
	updateFields := make(map[string]interface{})

	if updateFile.FileName != nil {
		updateFields["file_name"] = *updateFile.FileName
	}

	if updateFile.FileSize != nil {
		updateFields["file_size"] = *updateFile.FileSize
	}

	if updateFile.FileType != nil {
		updateFields["file_type"] = *updateFile.FileType
	}

	if len(updateFields) > 0 {
		if err := db.DB.Model(&schema.File{}).Where("id = ? AND uploaded_by = ?", fileId, userId).Updates(updateFields).Error; err != nil {
			return nil, err
		}
	}

	return f.GetFile(fileId)
}
