package schema

import (
	"time"

	"github.com/google/uuid"
)

type File struct {
	Id       uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	FolderId uuid.UUID `gorm:"type:uuid;not null" json:"folder_id"`
	Folder   Folder    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`

	FileName string `gorm:"not null;size:255" json:"file_name"`
	FileType string `gorm:"size:100" json:"file_type"`
	FileSize int64  `json:"file_size"`
	FileUrl  string `gorm:"not null" json:"file_url"`

	UploadedById uuid.UUID `gorm:"type:uuid;not null" json:"uploaded_by"`
	UploadedBy   User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
