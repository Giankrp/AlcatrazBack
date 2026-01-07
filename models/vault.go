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
	Type     VaultItemType `gorm:"not null;index"`

	// BaseVaultItem (Visible/Metadata)
	Title   string `gorm:"not null"`
	Icon    string `gorm:"default:'default_icon'"`
	Trashed bool   `gorm:"default:false;index"`

	// Specific Data (Encrypted)
	// En el frontend:
	// - PasswordItem: { username, password, url }
	// - NoteItem: { note }
	// - CardItem: { holder, number, expiry, cvv }
	// - IdentityItem: { firstName, lastName, email, phone, address }
	//
	// Todo este objeto específico se serializa a JSON, se cifra, y se guarda aquí.
	EncryptedData string `gorm:"not null"`

	// Encryption Metadata
	IV   string `gorm:"not null"`
	Salt string `gorm:"not null"` // Si usas KDF único por item

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"` // Soft delete para "papelera" real si se desea
}

// Estructura auxiliar para guardar datos adicionales NO cifrados si fuera necesario (ej. tags)
type VaultItemMeta struct {
	Tags []string `json:"tags"`
}

// Si en el futuro quisieras guardar parte de la data como JSONB consultable (NO cifrado):
type VaultItemPublicData struct {
	Data datatypes.JSON `gorm:"type:jsonb"`
}
