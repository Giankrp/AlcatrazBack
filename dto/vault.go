package dto

type CreateVaultItemDTO struct {
	FolderID   *string `json:"folder_id"`
	ItemType   string  `json:"item_type" validate:"required"`
	Title      string  `json:"title" validate:"required"`
	Ciphertext string  `json:"ciphertext" validate:"required"`
	IV         string  `json:"iv" validate:"required"`
	Salt       string  `json:"salt" validate:"required"`
}
type CreateVaultFolderDTO struct {
	Name string `json:"name" validate:"required"`
}
