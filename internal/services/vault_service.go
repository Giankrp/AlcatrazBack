package services

import (
	"errors"
	"strings"

	"github.com/Giankrp/AlcatrazBack/internal/dto"
	"github.com/Giankrp/AlcatrazBack/internal/models"
	"github.com/Giankrp/AlcatrazBack/internal/repositories"
	"github.com/google/uuid"
)

type VaultService interface {
	// Items
	CreateItem(userID uuid.UUID, input dto.CreateVaultItemDTO) (*models.VaultItem, error)
	GetItems(userID uuid.UUID) ([]models.VaultItem, error)
	GetItem(userID uuid.UUID, itemID uuid.UUID) (*models.VaultItem, error)
	UpdateItem(userID uuid.UUID, itemID uuid.UUID, input dto.UpdateVaultItemDTO) (*models.VaultItem, error)
	DeleteItem(userID uuid.UUID, itemID uuid.UUID) error

	// Folders
	CreateFolder(userID uuid.UUID, input dto.CreateVaultFolderDTO) (*models.VaultFolder, error)
	GetFolders(userID uuid.UUID) ([]models.VaultFolder, error)
	UpdateFolder(userID uuid.UUID, folderID uuid.UUID, input dto.UpdateVaultFolderDTO) (*models.VaultFolder, error)
	DeleteFolder(userID uuid.UUID, folderID uuid.UUID) error
}

type vaultService struct {
	repo repositories.VaultRepository
}

func NewVaultService(repo repositories.VaultRepository) VaultService {
	return &vaultService{repo: repo}
}

func (s *vaultService) CreateItem(userID uuid.UUID, input dto.CreateVaultItemDTO) (*models.VaultItem, error) {
	// Validate or Assign Folder
	var folderID *uuid.UUID
	if input.FolderID == nil || *input.FolderID == uuid.Nil {
		defaultFolder, err := s.repo.FindDefaultFolder(userID)
		if err != nil {
			return nil, errors.New("failed to find default folder")
		}
		folderID = &defaultFolder.ID
	} else {
		// Validate that folder belongs to user
		_, err := s.repo.FindFolderByID(*input.FolderID, userID)
		if err != nil {
			return nil, errors.New("folder not found or unauthorized")
		}
		folderID = input.FolderID
	}

	item := &models.VaultItem{
		UserID:   userID,
		FolderID: folderID,
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

func (s *vaultService) GetItems(userID uuid.UUID) ([]models.VaultItem, error) {
	return s.repo.FindAllByUserID(userID)
}

func (s *vaultService) GetItem(userID uuid.UUID, itemID uuid.UUID) (*models.VaultItem, error) {
	return s.repo.FindByID(itemID, userID)
}

func (s *vaultService) UpdateItem(userID uuid.UUID, itemID uuid.UUID, input dto.UpdateVaultItemDTO) (*models.VaultItem, error) {
	item, err := s.repo.FindByID(itemID, userID)
	if err != nil {
		return nil, err
	}

	if input.FolderID != nil && *input.FolderID != uuid.Nil {
		// Validate that folder belongs to user
		_, err := s.repo.FindFolderByID(*input.FolderID, userID)
		if err != nil {
			return nil, errors.New("folder not found or unauthorized")
		}
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

func (s *vaultService) DeleteItem(userID uuid.UUID, itemID uuid.UUID) error {
	// Verificar que existe y pertenece al usuario
	_, err := s.repo.FindByID(itemID, userID)
	if err != nil {
		return errors.New("item not found or unauthorized")
	}
	return s.repo.Delete(itemID, userID)
}

// Folder Methods Implementation

func (s *vaultService) CreateFolder(userID uuid.UUID, input dto.CreateVaultFolderDTO) (*models.VaultFolder, error) {
	folder := &models.VaultFolder{
		UserID:    userID,
		Name:      input.Name,
		IsDefault: false, // Created via API are never default
	}

	if err := s.repo.CreateFolder(folder); err != nil {
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique constraint") {
			return nil, errors.New("a folder with this name already exists")
		}
		return nil, err
	}
	return folder, nil
}

func (s *vaultService) GetFolders(userID uuid.UUID) ([]models.VaultFolder, error) {
	return s.repo.FindFoldersByUserID(userID)
}

func (s *vaultService) UpdateFolder(userID uuid.UUID, folderID uuid.UUID, input dto.UpdateVaultFolderDTO) (*models.VaultFolder, error) {
	folder, err := s.repo.FindFolderByID(folderID, userID)
	if err != nil {
		return nil, errors.New("folder not found or unauthorized")
	}

	if input.Name != "" {
		folder.Name = input.Name
	}

	if err := s.repo.UpdateFolder(folder); err != nil {
		return nil, err
	}
	return folder, nil
}

func (s *vaultService) DeleteFolder(userID uuid.UUID, folderID uuid.UUID) error {
	// 1. Verify folder exists and belongs to user
	folder, err := s.repo.FindFolderByID(folderID, userID)
	if err != nil {
		return errors.New("folder not found or unauthorized")
	}

	// 2. Ensure it's not the default folder
	if folder.IsDefault {
		return errors.New("cannot delete default folder")
	}

	// 3. Find the default folder for re-assignment
	defaultFolder, err := s.repo.FindDefaultFolder(userID)
	if err != nil {
		return errors.New("failed to find default folder for re-assignment")
	}

	// 4. Perform transactional delete
	return s.repo.DeleteFolder(folderID, userID, defaultFolder.ID)
}
