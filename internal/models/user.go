package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type User struct {
	ID                   uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Email                string         `gorm:"unique;not null" json:"email"`
	PasswordHash         string         `gorm:"not null" json:"-"`
	TwoFactorEnabled     bool           `gorm:"default:false" json:"two_factor_enabled"`
	TwoFactorSecret      string         `json:"-"`
	TwoFactorBackupCodes datatypes.JSON `json:"-"`
	CreatedAt            time.Time      `json:"created_at"`
}
