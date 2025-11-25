package models

import "time"

type Session struct {
	ID        string  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    string  `gorm:"index;not null"`
	DeviceID  *string `gorm:"index"`
	IP        string
	UserAgent string
	ExpiresAt time.Time `gorm:"index"`
	CreatedAt time.Time
}
