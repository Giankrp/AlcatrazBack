// Alcatraz - Secure open source Password Manager and Storage System
// Copyright (C) 2026 Gian Carlo Ruiz Patiño
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

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
	PasswordHash               string         `gorm:"not null" json:"-"`
	RecoveryKeyHash            string         `gorm:"not null" json:"-"` // Hash of the RK to verify recovery requests
	ProtectedMasterKey         string         `gorm:"not null" json:"protected_master_key"`
	MasterKeyIV                string         `gorm:"not null" json:"master_key_iv"`
	MasterKeySalt              string         `gorm:"not null" json:"master_key_salt"`
	RecoveryProtectedMasterKey string         `gorm:"not null" json:"recovery_protected_master_key"` // MK encrypted with RK
	RecoveryKeyIV              string         `gorm:"not null" json:"recovery_key_iv"`               // IV for decrypting the MK with RK
	RecoveryKeySalt            string         `gorm:"not null" json:"recovery_key_salt"`             // Salt for deriving the KEK from RK
	TwoFactorEnabled     bool           `gorm:"default:false" json:"two_factor_enabled"`
	TwoFactorSecret      string         `json:"-"`
	TwoFactorBackupCodes datatypes.JSON `json:"-"`
	CreatedAt            time.Time      `json:"created_at"`
}
