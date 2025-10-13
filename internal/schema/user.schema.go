package schema

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID           uuid.UUID      `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	Username     string         `gorm:"uniqueIndex;not null;size:100" json:"username" validate:"required,min=3,max=50"`
	Email        string         `gorm:"uniqueIndex;not null;size:100" json:"email" validate:"required,email"`
	Password     string         `gorm:"not null" json:"-"`
	ProfileUrl   string         `gorm:"size:500" json:"profile_url,omitempty"`
	CustomLink   *string        `gorm:"size:255" json:"custom_link,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	IsVerified   bool           `gorm:"default:false" json:"is_verified"`
	VerifyCode   string         `gorm:"size:4" json:"-"`
	CodeExpiry   time.Time      `json:"-"`
	StorageUsed  int64          `gorm:"default:0" json:"storage_used"`
	StorageLimit int64          `gorm:"default:524288000" json:"storage_limit"`
	Role         string         `gorm:"default:user" json:"role"` // e.g. "user" | "admin"
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

var ADMIN_EMAIL string

func init() {
	ADMIN_EMAIL = os.Getenv("ADMIN_EMAIL")
	if ADMIN_EMAIL == "" {
		fmt.Printf(`Failed to get the ADMIN_EMAIL`)
	}
}

// UpdateUserRequest defines which fields can be updated

// UpdateUserRequest defines which fields can be updated
type UpdateUserRequest struct {
	Username   *string `json:"username,omitempty" validate:"omitempty,min=3,max=50"`
	ProfileUrl *string `json:"profile_url,omitempty"`
	CustomLink *string `json:"custom_link,omitempty"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()

	if strings.ToLower(u.Email) == ADMIN_EMAIL {
		u.Role = "admin"
	} else {
		u.Role = "user"
	}
	return nil
}
