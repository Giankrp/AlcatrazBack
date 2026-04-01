package models

import "time"

type UserProfile struct {
	UserID    string `gorm:"type:uuid;primaryKey;not null"`
	User      User   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Name      string `gorm:"default:''"`
	AvatarURL string `gorm:"default:''"`
	Language  string `gorm:"default:'es'"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
