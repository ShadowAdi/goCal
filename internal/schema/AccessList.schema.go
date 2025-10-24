package schema

import "github.com/google/uuid"

type AccessType string

const (
	View AccessType = "view"
	Edit AccessType = "edit"
)

type FileAccess struct {
	Id         uuid.UUID  `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	FileID     uuid.UUID  `gorm:"type:uuid;not null;index" json:"file_id"`
	UserId     uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	AccessType AccessType `gorm:"size:50;default:'view'" json:"access_type"`

	File File `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	User User `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
}
