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

package dto

import (
	"github.com/google/uuid"
)

// VaultItemType mirrors the model type to avoid cross-layer imports.
// The validate tag uses string literals directly, so these constants are
// informational only; they are not required by the validator itself.
type VaultItemType string

const (
	ItemTypePassword VaultItemType = "password"
	ItemTypeNote     VaultItemType = "note"
	ItemTypeCard     VaultItemType = "card"
	ItemTypeIdentity VaultItemType = "identity"
)

// CreateVaultItemDTO defines the payload for creating a new vault item.
// The Secret field must contain client-side encrypted data (Zero Knowledge).
type CreateVaultItemDTO struct {
	FolderID      *uuid.UUID    `json:"folder_id"`
	ItemType      VaultItemType `json:"type" validate:"required,oneof=password note card identity"`
	Title         string        `json:"title" validate:"required"`
	Icon          string        `json:"icon"`
	SecurityScore *int          `json:"security_score" validate:"omitempty,min=0,max=100"`

	// Encrypted Blob (JSON stringificado y cifrado)
	Secret Secret `json:"secret" validate:"required"`
}

// Secret holds the client-encrypted blob and its decryption metadata.
type Secret struct {
	Data string `json:"data" validate:"required"`
	Iv   string `json:"iv" validate:"required"`
	Salt string `json:"salt" validate:"required"`
}
// UpdateVaultItemDTO defines the payload for a partial update of a vault item.
// All fields except Secret are optional. Only non-zero values should be applied.
type UpdateVaultItemDTO struct {
	FolderID      *uuid.UUID    `json:"folder_id"`
	ItemType      VaultItemType `json:"type" validate:"omitempty,oneof=password note card identity"`
	Title         string        `json:"title"`
	Icon          string        `json:"icon"`
	Trashed       *bool         `json:"trashed"`
	SecurityScore *int          `json:"security_score" validate:"omitempty,min=0,max=100"`
	Secret        Secret        `json:"secret" validate:"required"`
}

// CreateVaultFolderDTO defines the payload for creating a new vault folder.
type CreateVaultFolderDTO struct {
	Name      string `json:"name" validate:"required"`
	IsDefault bool   `json:"is_default"` // Usually false when created via API, but good to have
}

// UpdateVaultFolderDTO defines the payload for updating a vault folder.
type UpdateVaultFolderDTO struct {
	Name string `json:"name" validate:"required"`
}
