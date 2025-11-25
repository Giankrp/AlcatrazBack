package models

import "time"

type User struct {
	ID           string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Email        string `gorm:"unique;not null"`
	PasswordHash string `gorm:"not null"`
	CreatedAt    time.Time
}
