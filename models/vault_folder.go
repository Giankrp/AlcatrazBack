package models

import "time"

type VaultFolder struct {
	ID        string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    string `gorm:"index;not null"`
	Name      string `gorm:"not null"`
	CreatedAt time.Time
}
