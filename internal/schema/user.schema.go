package schema

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID         uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Username   string         `gorm:"uniqueIndex;not null;size:100" json:"username" validate:"required,min=3,max=50"`
	Email      string         `gorm:"uniqueIndex;not null;size:100" json:"email" validate:"required,email"`
	Password   string         `gorm:"not null" json:"-"`
	ProfileUrl string         `gorm:"size:500" json:"profile_url,omitempty"`
	CustomLink *string        `gorm:"size:255" json:"custom_link,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// UpdateUserRequest defines which fields can be updated
type UpdateUserRequest struct {
	Username   *string `json:"username,omitempty" validate:"omitempty,min=3,max=50"`
	ProfileUrl *string `json:"profile_url,omitempty"`
	CustomLink *string `json:"custom_link,omitempty"`
}

func (User) TableName() string {
	return "users"
}
