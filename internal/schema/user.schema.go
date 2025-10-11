package schema

import (
	"time"

	"gorm.io/gorm"
)

// PronounType defines the enum for pronouns
type PronounType string

const (
	PronounHeHim    PronounType = "He/Him"
	PronounSheHer   PronounType = "She/Her"
	PronounTheyThem PronounType = "They/Them"
	PronounOther    PronounType = "Other"
)

// DateFormatType defines the enum for date formats
type DateFormatType string

const (
	DateFormatDDMMYYYY DateFormatType = "DD/MM/YYYY"
	DateFormatMMDDYYYY DateFormatType = "MM/DD/YYYY"
	DateFormatYYYYMMDD DateFormatType = "YYYY-MM-DD"
)

// TimeFormatType defines the enum for time formats
type TimeFormatType string

const (
	TimeFormat12h TimeFormatType = "12h"
	TimeFormat24h TimeFormatType = "24h"
)

type User struct {
	ID             uint            `gorm:"primaryKey;autoIncrement" json:"id"`
	Username       string          `gorm:"uniqueIndex;not null;size:100" json:"username" validate:"required,min=3,max=50"`
	Email          string          `gorm:"uniqueIndex;not null;size:100" json:"email" validate:"required,email"`
	Password       string          `gorm:"not null" json:"-"`
	Country        string          `gorm:"size:100" json:"country"`
	ProfileUrl     string          `gorm:"size:500" json:"profile_url,omitempty"`
	WelcomeMessage string          `gorm:"size:500" json:"welcome_message,omitempty"`
	Timezone       *string         `gorm:"size:50;default:'UTC'" json:"timezone,omitempty"`
	Pronouns       *PronounType    `gorm:"type:varchar(20)" json:"pronouns,omitempty"`
	DateFormat     *DateFormatType `gorm:"type:varchar(20);default:'YYYY-MM-DD'" json:"date_format,omitempty"`
	TimeFormat     *TimeFormatType `gorm:"type:varchar(5);default:'24h'" json:"time_format,omitempty"`
	CustomLink     *string         `gorm:"size:255" json:"custom_link,omitempty"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
	DeletedAt      gorm.DeletedAt  `gorm:"index" json:"-"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.Timezone == nil {
		defaultTz := "UTC"
		u.Timezone = &defaultTz
	}
	if u.DateFormat == nil {
		defaultDf := DateFormatYYYYMMDD
		u.DateFormat = &defaultDf
	}
	if u.TimeFormat == nil {
		defaultTf := TimeFormat24h
		u.TimeFormat = &defaultTf
	}
	return nil
}
