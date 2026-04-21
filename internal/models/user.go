package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// User represents the main user entity in the system, including security metadata.
type User struct {
	ID                   uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Email                string         `gorm:"unique;not null" json:"email"`
	PasswordHash         string         `gorm:"not null" json:"-"`
	ProtectedMasterKey   string         `gorm:"not null" json:"protected_master_key"` // MK encrypted with the KEK (password)
	MasterKeyIV          string         `gorm:"not null" json:"master_key_iv"`        // IV for decrypting the MK
	MasterKeySalt        string         `gorm:"not null" json:"master_key_salt"`      // Salt for deriving the KEK
	TwoFactorEnabled     bool           `gorm:"default:false" json:"two_factor_enabled"`
	TwoFactorSecret      string         `json:"-"`
	TwoFactorBackupCodes datatypes.JSON `json:"-"`
	CreatedAt            time.Time      `json:"created_at"`
}
