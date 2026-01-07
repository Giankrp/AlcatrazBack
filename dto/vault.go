package dto

type VaultItemType string

const (
	ItemTypePassword VaultItemType = "password"
	ItemTypeNote     VaultItemType = "note"
	ItemTypeCard     VaultItemType = "card"
	ItemTypeIdentity VaultItemType = "identity"
)

type CreateVaultItemDTO struct {
	FolderID *string       `json:"folder_id"`
	Type     VaultItemType `json:"type" validate:"required,oneof=password note card identity"`
	Title    string        `json:"title" validate:"required"`
	Icon     string        `json:"icon"`

	// Encrypted Blob (JSON stringificado y cifrado)
	EncryptedData string `json:"encrypted_data" validate:"required"`

	// Encryption Params
	IV   string `json:"iv" validate:"required"`
	Salt string `json:"salt" validate:"required"`
}

type UpdateVaultItemDTO struct {
	FolderID      *string       `json:"folder_id"`
	Type          VaultItemType `json:"type" validate:"omitempty,oneof=password note card identity"`
	Title         string        `json:"title"`
	Icon          string        `json:"icon"`
	Trashed       *bool         `json:"trashed"`
	EncryptedData string        `json:"encrypted_data"`
	IV            string        `json:"iv"`
	Salt          string        `json:"salt"`
}

type CreateVaultFolderDTO struct {
	Name string `json:"name" validate:"required"`
}
