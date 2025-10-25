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

func (fo *FolderService) GetFolder(folderId string) (*schema.Folder, error) {
	var folder *schema.Folder

	result := db.DB.Find(&folder).Where("id = ?", folderId)
	if result.Error != nil {
		logger.Error("Failed to get the folder %s ", result.Error)
		return nil, result.Error
	}
	return folder, nil
}

func (fo *FolderService) GetUserFolder(folderId string, userId string) (*schema.Folder, error) {
	var folder *schema.Folder

	result := db.DB.Find(&folder).Where("id = ? AND created_by = ?", folderId, userId)
	if result.Error != nil {
		logger.Error("Failed to get the folder %s ", result.Error)
		return nil, result.Error
	}
	return folder, nil
}

func (fo *FolderService) CreateFolder(folder *schema.Folder, userId string) (*schema.Folder, error) {
	var existingFolder *schema.Folder

	result := db.DB.Find(&existingFolder).Where("folder_name = ? AND created_by = ?", folder.FolderName, userId)
	if result.Error == nil {
		logger.Error("Folder Already Exists %s ", result.Error)
		return nil, result.Error
	}

	resultFolderCreation := db.DB.Create(folder)
	if resultFolderCreation.Error != nil {
		logger.Error("Failed to create folder %s ", resultFolderCreation.Error)
		return nil, result.Error
	}
	return folder, nil
}

func (fo *FolderService) DeleteFolder(folderId string, userId string) (message string, err error) {
	folderFound, err := fo.GetUserFolder(folderId, userId)
	if err != nil {
		logger.Error("Failed to get the folder %s", err)
		return "Failed to get the folder ", err
	}

	if deleteError := db.DB.Delete(folderFound).Error; err != nil {
		logger.Error("Failed to delete the folder %s ", deleteError)
		return "Failed to delete folder", deleteError
	}

	return "Folder Deleted Successfully", err

}
