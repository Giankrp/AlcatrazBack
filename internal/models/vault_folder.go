package models

import (
	"time"

	"github.com/google/uuid"
)

type VaultFolder struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_user_folder_name;not null" json:"user_id"`
	Name      string    `gorm:"not null;uniqueIndex:idx_user_folder_name" json:"name"`
	IsDefault bool      `gorm:"default:false" json:"is_default"`
	CreatedAt time.Time `json:"created_at"`
}
