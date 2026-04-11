package repositories

import (
	"github.com/Giankrp/AlcatrazBack/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type VaultRepository interface {
	// Item methods
	Create(item *models.VaultItem) error
	FindByID(id uuid.UUID, userID uuid.UUID) (*models.VaultItem, error)
	FindAllByUserID(userID uuid.UUID) ([]models.VaultItem, error)
	Update(item *models.VaultItem) error
	Delete(id uuid.UUID, userID uuid.UUID) error

	// Folder methods
	CreateFolder(folder *models.VaultFolder) error
	FindFoldersByUserID(userID uuid.UUID) ([]models.VaultFolder, error)
	FindFolderByID(id uuid.UUID, userID uuid.UUID) (*models.VaultFolder, error)
	FindDefaultFolder(userID uuid.UUID) (*models.VaultFolder, error)
	UpdateFolder(folder *models.VaultFolder) error
	DeleteFolder(id uuid.UUID, userID uuid.UUID, defaultFolderID uuid.UUID) error
}

type vaultRepository struct {
	db *gorm.DB
}

func NewVaultRepository(db *gorm.DB) VaultRepository {
	return &vaultRepository{db: db}
}

func (r *vaultRepository) Create(item *models.VaultItem) error {
	return r.db.Create(item).Error
}

func (r *vaultRepository) FindByID(id uuid.UUID, userID uuid.UUID) (*models.VaultItem, error) {
	var item models.VaultItem
	// Preload Secret for details view
	if err := r.db.Preload("Secret").Where("id = ? AND user_id = ?", id, userID).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *vaultRepository) FindAllByUserID(userID uuid.UUID) ([]models.VaultItem, error) {
	var items []models.VaultItem
	// Do NOT Preload Secret for list view (Performance Optimization)
	if err := r.db.Where("user_id = ? AND deleted_at IS NULL", userID).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *vaultRepository) Update(item *models.VaultItem) error {
	// Save essentially does a full update (upsert) on the model and its associations if configured
	// But to be safe and explicit with GORM associations, usually Session.Save works for the primary model.
	// For associations, we might need to be careful. However, with FullSaveAssociations or manual handling:
	// Let's use a transaction to ensure both are updated if we pass the whole object.
	return r.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(item).Error
}

func (r *vaultRepository) Delete(id uuid.UUID, userID uuid.UUID) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.VaultItem{}).Error
}

// Folder Methods Implementation

func (r *vaultRepository) CreateFolder(folder *models.VaultFolder) error {
	return r.db.Create(folder).Error
}

func (r *vaultRepository) FindFoldersByUserID(userID uuid.UUID) ([]models.VaultFolder, error) {
	var folders []models.VaultFolder
	if err := r.db.Where("user_id = ?", userID).Find(&folders).Error; err != nil {
		return nil, err
	}
	return folders, nil
}

func (r *vaultRepository) FindFolderByID(id uuid.UUID, userID uuid.UUID) (*models.VaultFolder, error) {
	var folder models.VaultFolder
	if err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&folder).Error; err != nil {
		return nil, err
	}
	return &folder, nil
}

func (r *vaultRepository) FindDefaultFolder(userID uuid.UUID) (*models.VaultFolder, error) {
	var folder models.VaultFolder
	if err := r.db.Where("user_id = ? AND is_default = ?", userID, true).First(&folder).Error; err != nil {
		return nil, err
	}
	return &folder, nil
}

func (r *vaultRepository) UpdateFolder(folder *models.VaultFolder) error {
	return r.db.Save(folder).Error
}

func (r *vaultRepository) DeleteFolder(id uuid.UUID, userID uuid.UUID, defaultFolderID uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Reassign items to default folder
		if err := tx.Model(&models.VaultItem{}).
			Where("user_id = ? AND folder_id = ?", userID, id).
			Update("folder_id", defaultFolderID).Error; err != nil {
			return err
		}

		// 2. Delete the folder
		if err := tx.Where("id = ? AND user_id = ? AND is_default = ?", id, userID, false).
			Delete(&models.VaultFolder{}).Error; err != nil {
			return err
		}

		return nil
	})
}
