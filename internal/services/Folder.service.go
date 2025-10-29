package services

import (
	"goCal/internal/db"
	"goCal/internal/logger"
	"goCal/internal/schema"
)

type FolderService struct {
}

func NewFolderService() *FolderService {
	return &FolderService{}
}

func (fo *FolderService) GetFolders() ([]*schema.Folder, error) {
	var folders []*schema.Folder
	result := db.DB.Find(&folders)
	if result.Error != nil {
		logger.Error(`Failed to get Folders %w`, result.Error)
		return nil, result.Error
	}
	return folders, nil
}

func (fo *FolderService) GetFolder(folderId string) (*schema.Folder, error) {
	var folder *schema.Folder
	result := db.DB.Where("id = ?", folderId).First(&folder)
	if result.Error != nil {
		logger.Error(`Failed to get Folder %w`, result.Error)
		return nil, result.Error
	}
	return folder, nil
}

func (fo *FolderService) GetFolderByName(folderName string) (*schema.Folder, error) {
	var folder *schema.Folder
	result := db.DB.Where("folder_name = ?", folderName).First(&folder)
	if result.Error != nil {
		logger.Error(`Failed to get Folder %w`, result.Error)
		return nil, result.Error
	}
	return folder, nil
}

func (fo *FolderService) GetFolderByNameForUser(folderName string, userId string) (*schema.Folder, error) {
	var folder *schema.Folder
	result := db.DB.Where("folder_name = ? AND created_by = ?", folderName, userId).First(&folder)
	if result.Error != nil {
		logger.Error(`Failed to get Folder %w`, result.Error)
		return nil, result.Error
	}
	return folder, nil
}

func (fo *FolderService) CreateFolder(userId string, folderData *schema.Folder) (*schema.Folder, error) {
	var existingFolder *schema.Folder
	errFoundFolder := db.DB.Where("folder_name = ? AND created_by = ?", folderData.FolderName, userId).First(&existingFolder).Error

	if errFoundFolder != nil {
		return nil, errFoundFolder
	}

	result := db.DB.Create(&existingFolder)
	if result.Error != nil {
		return nil, result.Error
	}
	return fo.GetFolderByNameForUser(existingFolder.FolderName, userId)
}

func (fo *FolderService) DeleteFolder(userId string, folderId string) (string, error) {
	var existingFolder *schema.Folder
	errFoundFolder := db.DB.Where("id = ? AND created_by = ?", folderId, userId).First(&existingFolder).Error

	if errFoundFolder != nil {
		return "Failed to Delete Folder. Folder Not Found", errFoundFolder
	}

	result := db.DB.Delete(&existingFolder)
	if result.Error != nil {
		return "Failed to delete Folder.", result.Error
	}
	return "Folder Deleted Successfully", nil
}
