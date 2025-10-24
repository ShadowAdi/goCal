package schema

import "github.com/google/uuid"

type Folder struct {
	ID                uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	FolderName        string    `gorm:"uniqueIndex;not null;size:200" json:"folder_name" validate:"required,min=3,max=50"`
	FolderDescription string    `gorm:"not null;size:500" json:"folder_description"`
	FolderTags        []string  `gorm:"type:text[];default:[]" json:"folder_tags"`
	CreatedById       uuid.UUID `gorm:"type:uuid;not null" json:"created_by"`
	createdBy         User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:OnDelete;" json:"-"`

	Files []File `gorm:"foreignKey:FolderId" json:"files"`
}
