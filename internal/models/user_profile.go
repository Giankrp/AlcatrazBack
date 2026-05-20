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
