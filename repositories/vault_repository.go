package repositories

import (
	"github.com/Giankrp/AlcatrazBack/models"
	"gorm.io/gorm"
)

type VaultRepository interface {
	Create(item *models.VaultItem) error
	FindByID(id string, userID string) (*models.VaultItem, error)
	FindAllByUserID(userID string) ([]models.VaultItem, error)
	Update(item *models.VaultItem) error
	Delete(id string, userID string) error
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

func (r *vaultRepository) FindByID(id string, userID string) (*models.VaultItem, error) {
	var item models.VaultItem
	// Preload Secret for details view
	if err := r.db.Preload("Secret").Where("id = ? AND user_id = ?", id, userID).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *vaultRepository) FindAllByUserID(userID string) ([]models.VaultItem, error) {
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

func (r *vaultRepository) Delete(id string, userID string) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.VaultItem{}).Error
}
