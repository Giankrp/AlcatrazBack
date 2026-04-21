// Package models contains the data models for the alcatraz database.
package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type VaultItemType string

const (
	ItemTypePassword VaultItemType = "password"
	ItemTypeNote     VaultItemType = "note"
	ItemTypeCard     VaultItemType = "card"
	ItemTypeIdentity VaultItemType = "identity"
)

type VaultItem struct {
	ID       uuid.UUID     `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID   uuid.UUID     `gorm:"type:uuid;index;not null" json:"user_id"`
	User     User          `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"-"`
	FolderID *uuid.UUID    `gorm:"type:uuid;index" json:"folder_id"`
	ItemType VaultItemType `gorm:"column:item_type;not null;index" json:"item_type"`

	// BaseVaultItem (Visible/Metadata)
	Title   string `gorm:"not null" json:"title"`
	Icon    string `gorm:"default:'default_icon'" json:"icon"`
	Trashed bool   `gorm:"default:false;index" json:"trashed"`

	// Relation
	Secret *VaultSecret `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"secret,omitempty"` // HasOne relation

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"-"` // Soft delete hidden from JSON usually
}

// VaultSecret holds the heavy encrypted data.
// It is separated to allow list queries without loading this blob.
type VaultSecret struct {
	VaultItemID uuid.UUID `gorm:"type:uuid;primaryKey" json:"vault_item_id"` // Uses the same ID as VaultItem for 1:1

	// Specific Data (Encrypted)
	// Serialized JSON encrypted with the user's Master Key.
	EncryptedData string `gorm:"not null" json:"encrypted_data"`

	// Encryption Metadata
	IV   string `gorm:"not null" json:"iv"`   // Unique IV per item.
	Salt string `gorm:"not null" json:"salt"` // Optional Salt used to derive per-item keys from the Master Key.
}

// VaultItemMeta is an auxiliary structure to store non-encrypted additional data (e.g., tags).
type VaultItemMeta struct {
	Tags []string `json:"tags"`
}

// VaultItemPublicData represents data as queryable JSONB (NOT encrypted).
type VaultItemPublicData struct {
	Data datatypes.JSON `gorm:"type:jsonb"`
}
