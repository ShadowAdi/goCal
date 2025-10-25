package services

import (
	"goCal/internal/db"
	"goCal/internal/logger"
	"goCal/internal/schema"
)

type FolderService struct{}

func NewFolderService() *FolderService {
	return &FolderService{}
}

func (fo *FolderService) GetFolders() ([]*schema.Folder, error) {
	var folders []*schema.Folder

	result := db.DB.Find(&folders)
	if result.Error != nil {
		logger.Error("Failed to get all the folder %s ", result.Error)
		return nil, result.Error
	}
	return folders, nil
}
