package services

import (
	"fmt"
	"goCal/internal/db"
	"goCal/internal/logger"
	"goCal/internal/schema"

	"gorm.io/gorm"
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) GetUsers() ([]*schema.User, error) {
	var users []*schema.User
	// This automatically excludes soft-deleted records due to GORM's default behavior
	result := db.DB.Find(&users)
	if result.Error != nil {
		logger.Error(`Failed to get Users %w`, result.Error)
		return nil, result.Error
	}
	return users, nil
}

func (s *UserService) GetUser(id string) (*schema.User, error) {
	var user *schema.User
	// This automatically excludes soft-deleted records due to GORM's default behavior
	result := db.DB.Where("id = ?", id).First(&user)
	if result.Error != nil {
		logger.Error(`Failed to get User  %w with id %s`, result.Error, id)
		return nil, result.Error
	}
	return user, nil
}

func (s *UserService) GetUserByEmail(email string) (*schema.User, error) {
	var user *schema.User
	// This automatically excludes soft-deleted records
	result := db.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

// GetUserByEmailIncludingDeleted gets user by email including soft-deleted users
func (s *UserService) GetUserByEmailIncludingDeleted(email string) (*schema.User, error) {
	var user *schema.User
	result := db.DB.Unscoped().Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

// GetSoftDeletedUsers returns all soft-deleted users
func (s *UserService) GetSoftDeletedUsers() ([]*schema.User, error) {
	var users []*schema.User
	result := db.DB.Unscoped().Where("deleted_at IS NOT NULL").Find(&users)
	if result.Error != nil {
		logger.Error(`Failed to get soft-deleted users %w`, result.Error)
		return nil, result.Error
	}
	return users, nil
}

// RestoreUser restores a soft-deleted user
func (s *UserService) RestoreUser(id string) (*schema.User, error) {
	var user *schema.User
	// Find the soft-deleted user
	result := db.DB.Unscoped().Where("id = ? AND deleted_at IS NOT NULL", id).First(&user)
	if result.Error != nil {
		return nil, fmt.Errorf("soft-deleted user not found: %w", result.Error)
	}

	// Restore by setting deleted_at to NULL
	user.DeletedAt = gorm.DeletedAt{}
	if err := db.DB.Save(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to restore user: %w", err)
	}

	return user, nil
}

// PermanentlyDeleteUser permanently deletes a user (hard delete)
func (s *UserService) PermanentlyDeleteUser(id string) error {
	var user *schema.User
	// Use Unscoped to find even soft-deleted users
	result := db.DB.Unscoped().Where("id = ?", id).First(&user)
	if result.Error != nil {
		return fmt.Errorf("user not found: %w", result.Error)
	}

	// Permanently delete
	if err := db.DB.Unscoped().Delete(&user).Error; err != nil {
		return fmt.Errorf("failed to permanently delete user: %w", err)
	}

	return nil
}

func (s *UserService) CreateUser(newUser *schema.User) (*schema.User, error) {
	// Check if user with this email exists (including soft-deleted)
	var existingUser *schema.User
	err := db.DB.Unscoped().Where("email = ?", newUser.Email).First(&existingUser).Error

	if err == nil {
		// User exists
		if existingUser.DeletedAt.Valid {
			// User is soft-deleted, restore and update
			existingUser.DeletedAt = gorm.DeletedAt{}
			existingUser.Username = newUser.Username
			existingUser.Password = newUser.Password
			existingUser.Country = newUser.Country
			existingUser.ProfileUrl = newUser.ProfileUrl
			existingUser.WelcomeMessage = newUser.WelcomeMessage
			existingUser.Timezone = newUser.Timezone
			existingUser.Pronouns = newUser.Pronouns
			existingUser.DateFormat = newUser.DateFormat
			existingUser.TimeFormat = newUser.TimeFormat
			existingUser.CustomLink = newUser.CustomLink

			if err := db.DB.Save(&existingUser).Error; err != nil {
				return nil, err
			}
			return existingUser, nil
		} else {
			// User exists and is not deleted - return error
			return nil, fmt.Errorf("user with email %s already exists", newUser.Email)
		}
	}

	// User doesn't exist, create new one
	result := db.DB.Create(newUser)
	if result.Error != nil {
		return nil, result.Error
	}
	return newUser, nil
}

func (s *UserService) DeleteUser(id string) (string, error) {
	userFound, err := s.GetUser(id)
	if err != nil {
		return "User Not Found", err
	}
	if err := db.DB.Delete(&userFound).Error; err != nil {
		return "Failed to delete user", err
	}
	return "User deleted successfully", nil
}

func (s *UserService) UpdateUser(id string, updateRequest *schema.UpdateUserRequest) (*schema.User, error) {
	// Create a map of only the non-nil fields to update
	updateFields := make(map[string]interface{})

	if updateRequest.Username != nil {
		updateFields["username"] = *updateRequest.Username
	}
	if updateRequest.Country != nil {
		updateFields["country"] = *updateRequest.Country
	}
	if updateRequest.ProfileUrl != nil {
		updateFields["profile_url"] = *updateRequest.ProfileUrl
	}
	if updateRequest.WelcomeMessage != nil {
		updateFields["welcome_message"] = *updateRequest.WelcomeMessage
	}
	if updateRequest.Timezone != nil {
		updateFields["timezone"] = *updateRequest.Timezone
	}
	if updateRequest.Pronouns != nil {
		updateFields["pronouns"] = *updateRequest.Pronouns
	}
	if updateRequest.DateFormat != nil {
		updateFields["date_format"] = *updateRequest.DateFormat
	}
	if updateRequest.TimeFormat != nil {
		updateFields["time_format"] = *updateRequest.TimeFormat
	}
	if updateRequest.CustomLink != nil {
		updateFields["custom_link"] = *updateRequest.CustomLink
	}

	// Update only the specified fields
	if len(updateFields) > 0 {
		if err := db.DB.Model(&schema.User{}).Where("id = ?", id).Updates(updateFields).Error; err != nil {
			return nil, err
		}
	}

	// Fetch and return the updated user
	return s.GetUser(id)
}
