package services

import (
	"errors"

	"github.com/Giankrp/AlcatrazBack/dto"
	"github.com/Giankrp/AlcatrazBack/models"
	"github.com/Giankrp/AlcatrazBack/repositories"
)

type VaultService interface {
	CreateItem(userID string, input dto.CreateVaultItemDTO) (*models.VaultItem, error)
	GetItems(userID string) ([]models.VaultItem, error)
	GetItem(userID string, itemID string) (*models.VaultItem, error)
	UpdateItem(userID string, itemID string, input dto.UpdateVaultItemDTO) (*models.VaultItem, error)
	DeleteItem(userID string, itemID string) error
}

type vaultService struct {
	repo repositories.VaultRepository
}

func NewVaultService(repo repositories.VaultRepository) VaultService {
	return &vaultService{repo: repo}
}

func (s *vaultService) CreateItem(userID string, input dto.CreateVaultItemDTO) (*models.VaultItem, error) {
	item := &models.VaultItem{
		UserID:   userID,
		FolderID: input.FolderID,
		ItemType: models.VaultItemType(input.ItemType),
		Title:    input.Title,
		Icon:     input.Icon,
		Secret: &models.VaultSecret{
			EncryptedData: input.Secret.Data,
			IV:            input.Secret.Iv,
			Salt:          input.Secret.Salt,
		},
	}

	if err := s.repo.Create(item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *vaultService) GetItems(userID string) ([]models.VaultItem, error) {
	return s.repo.FindAllByUserID(userID)
}

func (s *vaultService) GetItem(userID string, itemID string) (*models.VaultItem, error) {
	return s.repo.FindByID(itemID, userID)
}

func (s *vaultService) UpdateItem(userID string, itemID string, input dto.UpdateVaultItemDTO) (*models.VaultItem, error) {
	item, err := s.repo.FindByID(itemID, userID)
	if err != nil {
		return nil, err
	}

	if input.FolderID != nil {
		item.FolderID = input.FolderID
	}
	if input.ItemType != "" {
		item.ItemType = models.VaultItemType(input.ItemType)
	}
	if input.Title != "" {
		item.Title = input.Title
	}
	if input.Icon != "" {
		item.Icon = input.Icon
	}
	// Update Secret Data
	if input.Secret.Data != "" || input.Secret.Iv != "" || input.Secret.Salt != "" {
		if item.Secret == nil {
			item.Secret = &models.VaultSecret{VaultItemID: item.ID}
		}
		if input.Secret.Data != "" {
			item.Secret.EncryptedData = input.Secret.Data
		}
		if input.Secret.Iv != "" {
			item.Secret.IV = input.Secret.Iv
		}
		if input.Secret.Salt != "" {
			item.Secret.Salt = input.Secret.Salt
		}
	}

	if input.Trashed != nil {
		item.Trashed = *input.Trashed
	}

	if err := s.repo.Update(item); err != nil {
		return nil, err
	}

	return item, nil
}

func (s *vaultService) DeleteItem(userID string, itemID string) error {
	// Verificar que existe y pertenece al usuario
	_, err := s.repo.FindByID(itemID, userID)
	if err != nil {
		return errors.New("item not found or unauthorized")
	}
	return s.repo.Delete(itemID, userID)
}
