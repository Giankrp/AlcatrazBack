package models

import "time"

type VaultItem struct {
	ID       string  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID   string  `gorm:"index;not null"`
	FolderID *string `gorm:"index"`

	ItemType string `gorm:"not null"`
	Title    string `gorm:"not null"`

	Ciphertext string `gorm:"not null"`
	IV         string `gorm:"not null"`
	Salt       string `gorm:"not null"`

	CreatedAt time.Time
}
