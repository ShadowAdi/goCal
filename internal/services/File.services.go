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
