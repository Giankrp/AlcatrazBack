package dto

import (
	"github.com/google/uuid"
)

type VaultItemType string

const (
	ItemTypePassword VaultItemType = "password"
	ItemTypeNote     VaultItemType = "note"
	ItemTypeCard     VaultItemType = "card"
	ItemTypeIdentity VaultItemType = "identity"
)

type CreateVaultItemDTO struct {
	FolderID *uuid.UUID    `json:"folder_id"`
	ItemType VaultItemType `json:"type" validate:"required,oneof=password note card identity"`
	Title    string        `json:"title" validate:"required"`
	Icon     string        `json:"icon"`

	// Encrypted Blob (JSON stringificado y cifrado)
	Secret Secret `json:"secret" validate:"required"`
}

type Secret struct {
	Data string `json:"data" validate:"required"`
	Iv   string `json:"iv" validate:"required"`
	Salt string `json:"salt" validate:"required"`
}
type UpdateVaultItemDTO struct {
	FolderID *uuid.UUID    `json:"folder_id"`
	ItemType VaultItemType `json:"type" validate:"omitempty,oneof=password note card identity"`
	Title    string        `json:"title"`
	Icon     string        `json:"icon"`
	Trashed  *bool         `json:"trashed"`
	Secret   Secret        `json:"secret" validate:"required"`
}

type CreateVaultFolderDTO struct {
	Name      string `json:"name" validate:"required"`
	IsDefault bool   `json:"is_default"` // Usually false when created via API, but good to have
}

type UpdateVaultFolderDTO struct {
	Name string `json:"name" validate:"required"`
}
