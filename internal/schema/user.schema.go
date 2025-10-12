package schema

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID                uuid.UUID      `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	Username          string         `gorm:"uniqueIndex;not null;size:100" json:"username" validate:"required,min=3,max=50"`
	Email             string         `gorm:"uniqueIndex;not null;size:100" json:"email" validate:"required,email"`
	Password          string         `gorm:"not null" json:"-"`
	ProfileUrl        string         `gorm:"size:500" json:"profile_url,omitempty"`
	CustomLink        *string        `gorm:"size:255" json:"custom_link,omitempty"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	StorageUsed       int64          `gorm:"default:0" json:"storage_used"`
	StorageLimit      int64          `gorm:"default:524288000" json:"storage_limit"`
	Role              string         `gorm:"default:user" json:"role"` // e.g. "user" | "admin"
	EmailVerified     bool           `gorm:"default:false" json:"email_verified"`
	VerificationToken *string        `gorm:"size:255" json:"-"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
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
