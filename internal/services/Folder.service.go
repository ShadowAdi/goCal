package services

import (
	"errors"
	"fmt"
	"goCal/internal/db"
	"goCal/internal/logger"
	"goCal/internal/schema"

	"github.com/google/uuid"
	"gorm.io/gorm"
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
	// check if folder already exists
	var existingFolder schema.Folder
	err := db.DB.Where("folder_name = ? AND created_by = ?", folderData.FolderName, userId).First(&existingFolder).Error

	if err == nil {
		// folder already exists
		return nil, fmt.Errorf("folder with the same name already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		// actual DB error
		logger.Error(`DB error while checking existing folder: %v`, err)
		return nil, err
	}

	// create new folder
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, fmt.Errorf("invalid userId format")
	}

	folderData.CreatedById = userUUID
	result := db.DB.Create(folderData)
	if result.Error != nil {
		logger.Error(`Failed to create folder: %v`, result.Error)
		return nil, result.Error
	}

	return folderData, nil
}

func (fo *FolderService) DeleteFolder(userId string, folderId string) (string, error) {
	var existingFolder *schema.Folder
	errFoundFolder := db.DB.Where("id = ? AND created_by = ?", folderId, userId).First(&existingFolder).Error

	if errFoundFolder != nil {
		logger.Error(`Failed to delete folder. Folder Not Found %w`, errFoundFolder.Error())
		return "Failed to Delete Folder. Folder Not Found", errFoundFolder
	}

	result := db.DB.Delete(&existingFolder)
	if result.Error != nil {
		logger.Error(`Failed to delete folder. %w`, result.Error)
		return "Failed to delete Folder.", result.Error
	}
	return "Folder Deleted Successfully", nil
}

func (fo *FolderService) UpdateFolder(userId string, folderId string, folderData *schema.UpdateFolderRequest) (*schema.Folder, error) {
	var existingFolder *schema.Folder
	errFoundFolder := db.DB.Where("id = ? AND created_by = ?", folderId, userId).First(&existingFolder).Error

	if errFoundFolder != nil {
		logger.Error(`Failed to Update folder. Not Found %w`, errFoundFolder.Error())
		return nil, errFoundFolder
	}

	updatedFields := make(map[string]interface{})

	if folderData.FolderName != nil {
		updatedFields["folder_name"] = *folderData.FolderName
	}

	if folderData.FolderDescription != nil {
		updatedFields["folder_description"] = *folderData.FolderDescription
	}

	if folderData.FolderTags != nil {
		updatedFields["folder_tags"] = folderData.FolderTags
	}

	if len(updatedFields) > 0 {
		if err := db.DB.Model(&schema.Folder{}).Where("id = ? AND created_by = ?", folderId, userId).Updates(updatedFields); err != nil {
			logger.Error("Failed to Update Folder %w", err)
			return nil, err.Error
		}
	}

	return existingFolder, nil
}
