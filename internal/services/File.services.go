package services

import (
	"goCal/internal/db"
	"goCal/internal/logger"
	"goCal/internal/schema"
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

func (f *FileService) CreateFile(file *schema.File, userId string) (*schema.File, error) {
	var existingFile *schema.File
	result := db.DB.Find(&existingFile).Where("file_name = ? AND uploaded_by = ?", file.FileName, userId)
	if result.Error == nil {
		logger.Error("File With same name already Exists. Choose another One %s ", result.Error)
		return nil, result.Error
	}

	resultFileCreation := db.DB.Create(file)
	if resultFileCreation.Error != nil {
		logger.Error("Failed to create file %s", resultFileCreation.Error)
	}

	return file, nil
}
