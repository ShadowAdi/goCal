package repository

import (
	"goCal/internal/schema"

	"gorm.io/gorm"
)

// UserRepository defines the interface for user database operations
type UserRepository interface {
	Create(user *schema.User) error
	GetByID(id uint) (*schema.User, error)
	GetByEmail(email string) (*schema.User, error)
	GetByUsername(username string) (*schema.User, error)
	Update(user *schema.User) error
	Delete(id uint) error
	List(limit, offset int) ([]schema.User, error)
	Count() (int64, error)
}

// userRepository implements UserRepository
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// Create creates a new user
func (r *userRepository) Create(user *schema.User) error {
	return r.db.Create(user).Error
}

// GetByID retrieves a user by ID
func (r *userRepository) GetByID(id uint) (*schema.User, error) {
	var user schema.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *userRepository) GetByEmail(email string) (*schema.User, error) {
	var user schema.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByUsername retrieves a user by username
func (r *userRepository) GetByUsername(username string) (*schema.User, error) {
	var user schema.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update updates a user
func (r *userRepository) Update(user *schema.User) error {
	return r.db.Save(user).Error
}

// Delete soft deletes a user by ID
func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&schema.User{}, id).Error
}

// List retrieves users with pagination
func (r *userRepository) List(limit, offset int) ([]schema.User, error) {
	var users []schema.User
	err := r.db.Limit(limit).Offset(offset).Find(&users).Error
	return users, err
}

// Count returns the total number of users
func (r *userRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&schema.User{}).Count(&count).Error
	return count, err
}