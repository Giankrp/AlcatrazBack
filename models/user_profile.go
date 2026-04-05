package models

import (
	"time"

	"github.com/google/uuid"
)

type UserProfile struct {
	UserID    uuid.UUID `gorm:"type:uuid;primaryKey;not null" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	Name      string    `gorm:"default:''" json:"name"`
	AvatarURL string    `gorm:"default:''" json:"avatar_url"`
	Language  string    `gorm:"default:'es'" json:"language"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
