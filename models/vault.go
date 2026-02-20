package models

import (
	"time"

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
	ID       string        `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID   string        `gorm:"index;not null"`
	FolderID *string       `gorm:"index"`
	ItemType VaultItemType `gorm:"column:item_type;not null;index"`

	// BaseVaultItem (Visible/Metadata)
	Title   string `gorm:"not null"`
	Icon    string `gorm:"default:'default_icon'"`
	Trashed bool   `gorm:"default:false;index"`

	// Relation
	Secret *VaultSecret `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // HasOne relation

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"` // Soft delete
}

// VaultSecret holds the heavy encrypted data.
// It is separated to allow list queries without loading this blob.
type VaultSecret struct {
	VaultItemID string `gorm:"primaryKey"` // Uses the same ID as VaultItem for 1:1

	// Specific Data (Encrypted)
	// Serialized JSON of specific item data (PasswordItem, NoteItem, etc.)
	EncryptedData string `gorm:"not null"`

	// Encryption Metadata
	IV   string `gorm:"not null"`
	Salt string `gorm:"not null"`
}

// Estructura auxiliar para guardar datos adicionales NO cifrados si fuera necesario (ej. tags)
type VaultItemMeta struct {
	Tags []string `json:"tags"`
}

// Si en el futuro quisieras guardar parte de la data como JSONB consultable (NO cifrado):
type VaultItemPublicData struct {
	Data datatypes.JSON `gorm:"type:jsonb"`
}
