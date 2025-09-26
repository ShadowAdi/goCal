package schema

import "time"

type User struct {
	ID             string    `json:"id" db:"id"`
	Username       string    `json:"username" binding:"required"`
	Email          string    `json:"email" binding:"required"`
	Password       string    `json:"password" binding:"required"`
	ProfileUrl     string    `json:"profileUrl,omitempty"`
	Country        string    `json:"country" binding:"required"`
	WelcomeMessage string    `json:"welcome_message,omitempty"`
	Timezone       string    `json:"timezone,omitempty"`
	Pronouns       string    `json:"pronouns,omitempty"`
	IsVerified     bool      `json:"isverified" db:"isverified"`
	DateFormat     string    `json:"date_format,omitempty"`
	TimeFormat     string    `json:"time_format,omitempty"`
	CustomLink     string    `json:"custom_link,omitempty"`
	CreatedAt      time.Time `json:"created_at,omitempty"`
}
