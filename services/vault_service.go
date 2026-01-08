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
		UserID:        userID,
		FolderID:      input.FolderID,
		Type:          models.VaultItemType(input.Type),
		Title:         input.Title,
		Icon:          input.Icon,
		EncryptedData: input.EncryptedData,
		IV:            input.IV,
		Salt:          input.Salt,
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
	if input.Type != "" {
		item.Type = models.VaultItemType(input.Type)
	}
	if input.Title != "" {
		item.Title = input.Title
	}
	if input.Icon != "" {
		item.Icon = input.Icon
	}
	if input.EncryptedData != "" {
		item.EncryptedData = input.EncryptedData
	}
	if input.IV != "" {
		item.IV = input.IV
	}
	if input.Salt != "" {
		item.Salt = input.Salt
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
