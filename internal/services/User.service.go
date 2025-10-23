package services

import (
	"fmt"
	"goCal/internal/db"
	"goCal/internal/logger"
	"goCal/internal/schema"
	"time"

	"gorm.io/gorm"
)

type UserService struct {
	emailService *EmailService
}

func NewUserService() *UserService {
	emailService, err := NewEmailServices()
	if err != nil {
		logger.Error("Failed to initialize email service: %v", err)
		// You can decide whether to fail completely or continue without email service
		// For now, we'll log the error and continue
	}
	return &UserService{
		emailService: emailService,
	}
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
			existingUser.ProfileUrl = newUser.ProfileUrl
			existingUser.CustomLink = newUser.CustomLink

			if err := db.DB.Save(&existingUser).Error; err != nil {
				return nil, err
			}

			// Send verification email for restored user
			if s.emailService != nil {
				go func() {
					if _, err := s.emailService.SendVerificationEmail(existingUser); err != nil {
						logger.Error("Failed to send verification email to restored user %s: %v", existingUser.Email, err)
					}
				}()
			}

			return existingUser, nil
		} else {
			// User exists and is not deleted - return error
			return nil, fmt.Errorf("user with email %s already exists", newUser.Email)
		}
	}

	// Create new user
	result := db.DB.Create(newUser)
	if result.Error != nil {
		return nil, result.Error
	}

	// Send verification email for new user
	if s.emailService != nil {
		go func() {
			if emailResponse, err := s.emailService.SendVerificationEmail(newUser); err != nil {
				logger.Error("Failed to send verification email to new user %s: %v", newUser.Email, err)
			} else {
				logger.Info("Verification email queued for user %s: %s", newUser.Email, emailResponse.Message)
			}
		}()
	} else {
		logger.Warn("Email service not available - verification email not sent for user: %s", newUser.Email)
	}

	return newUser, nil
}

func (s *UserService) DeleteUser(id string) (string, error) {
	fmt.Printf("Attempting to delete user with ID: %s\n", id)
	userFound, err := s.GetUser(id)
	if err != nil {
		fmt.Printf("User not found for deletion: %v\n", err)
		return "User Not Found", err
	}
	fmt.Printf("User found for deletion: %+v\n", userFound)
	if err := db.DB.Delete(&userFound).Error; err != nil {
		fmt.Printf("Error deleting user: %v\n", err)
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
	if updateRequest.ProfileUrl != nil {
		updateFields["profile_url"] = *updateRequest.ProfileUrl
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

// ResendVerificationEmail resends verification email to a user
func (s *UserService) ResendVerificationEmail(email string) (*EmailResponse, error) {
	if s.emailService == nil {
		return &EmailResponse{
			Success: false,
			Message: "Email Service Not Available",
			Error:   "Email Service Not Available",
		}, fmt.Errorf("Email Service Not Available")
	}

	user, err := s.GetUserByEmail(email)

	if err != nil {
		return &EmailResponse{
			Success: false,
			Message: "User Not Found",
			Error:   err.Error(),
		}, err
	}

	if user.IsVerified {
		return &EmailResponse{
			Success: false,
			Message: "Already Verified",
			Error:   "Already Verified",
		}, fmt.Errorf("Already Verified")
	}

	user.VerifyCode = fmt.Sprintf("%40d")
	user.CodeExpiry = time.Now().Add(15 * time.Minute)

	if err := db.DB.Save(user).Error; err != nil {
		return &EmailResponse{
			Success: false,
			Message: "Failed to update verification code",
			Error:   err.Error(),
		}, err
	}

	return s.emailService.SendVerificationEmail(user)

}

// VerifyUser verifies a user with the provided verification code
func (s *UserService) VerifyUser(email, verificationCode string) (*schema.User, error) {
	user, err := s.GetUserByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("User Not Found")
	}
	if user.IsVerified {
		return user, nil
	}

	if user.VerifyCode != verificationCode {
		return nil, fmt.Errorf("Invalid Verification Code")
	}

	if time.Now().After(user.CodeExpiry) {
		return nil, fmt.Errorf("verification code has expired")
	}

	user.IsVerified = true
	user.VerifyCode = ""

	if err := db.DB.Save(user).Error; err != nil {
		return nil, fmt.Errorf("failed to update user verification status: %w", err)
	}

	logger.Info("User %s successfully verified", user.Email)
	return user, nil
}
