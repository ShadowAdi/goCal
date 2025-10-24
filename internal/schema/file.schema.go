package schema

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FileVisibility string

const (
	Private FileVisibility = "private"
	Shared  FileVisibility = "shared"
	Public  FileVisibility = "public"
)

type File struct {
	Id       uuid.UUID  `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	FolderId *uuid.UUID `gorm:"type:uuid" json:"folder_id,omitempty"`
	Folder   *Folder    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"folder,omitempty"`

	FileName string `gorm:"not null;size:255" json:"file_name"`
	FileType string `gorm:"size:100" json:"file_type"`
	FileSize int64  `json:"file_size"`
	FileUrl  string `gorm:"not null" json:"file_url"`

	Visibility FileVisibility `gorm:"type:varchar(20);default:'private'" json:"visibility"`

	UploadedById uuid.UUID `gorm:"type:uuid;not null" json:"uploaded_by"`
	UploadedBy   User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	AccessList []FileAccess `gorm:"foreignKey:FileId;constraint:OnDelete:CASCADE;" json:"access_list,omitempty"`
}

type UpdateFileRequest struct {
	FileName *string `json:"file_name,omitempty" validate:"omitempty,min=3,max=50"`
	FileType *string `json:"file_type"`
	FileSize *int64  `json:"file_size"`
}

func (File) TableName() string {
	return "files"
}

func (fc *File) BeforeCreate(tx *gorm.DB) (err error) {
	fc.Id = uuid.New()
	return nil
}
