package services

import (
	"goCal/internal/db"
	"goCal/internal/logger"
	"goCal/internal/schema"
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) GetUsers() ([]*schema.User, error) {
	var users []*schema.User
	result := db.DB.Find(&users)
	if result.Error != nil {
		logger.Error(`Failed to get Users %w`, result.Error)
		return nil, result.Error
	}
	return users, nil
}

func (s *UserService) GetUser(id string) (*schema.User, error) {
	var user *schema.User
	result := db.DB.Where("id = ?", id).First(&user)
	if result.Error != nil {
		logger.Error(`Failed to get User  %w withd id %d`, result.Error, id)
		return nil, result.Error
	}
	return user, nil
}

func (s *UserService) GetUserByEmail(email string) (*schema.User, error) {
	var user *schema.User
	result := db.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (s *UserService) CreateUser(newUser *schema.User) (*schema.User, error) {
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

func (s *UserService) UpdateUser(id string, updateUser *schema.User) (*schema.User, error) {
	if updateErr := db.DB.Where("id = ?", id).Updates(updateUser); updateErr.Error != nil {
		return nil, updateErr.Error
	}
	return updateUser, nil
}
