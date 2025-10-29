package schema

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Folder struct {
	ID                uuid.UUID      `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	FolderName        string         `gorm:"uniqueIndex;not null;size:200" json:"folder_name" validate:"required,min=3,max=50"`
	FolderDescription string         `gorm:"size:500" json:"folder_description"`
	FolderTags        pq.StringArray `gorm:"type:text[]" json:"folder_tags"`
	CreatedById       uuid.UUID      `gorm:"not null" json:"created_by_id"`
	CreatedBy         User           `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
}

type UpdateFolderRequest struct {
	FolderName        *string   `json:"folder_name,omitempty" validate:"omitempty,min=3,max=50"`
	FolderDescription *string   `json:"folder_description"`
	FolderTags        []*string `json:"folder_tags"`
}

func (Folder) TableName() string {
	return "folders"
}

func (fo *Folder) BeforeCreate(tx *gorm.DB) (err error) {
	fo.ID = uuid.New()
	return nil
}
