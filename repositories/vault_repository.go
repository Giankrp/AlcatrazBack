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
	// Buscamos por ID y UserID para asegurar que el item pertenece al usuario
	if err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *vaultRepository) FindAllByUserID(userID string) ([]models.VaultItem, error) {
	var items []models.VaultItem
	if err := r.db.Where("user_id = ? AND deleted_at IS NULL", userID).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *vaultRepository) Update(item *models.VaultItem) error {
	return r.db.Save(item).Error
}

func (r *vaultRepository) Delete(id string, userID string) error {
	// Soft delete (o hard delete si prefieres)
	// Aquí usamos Delete de GORM que manejará Soft Delete si DeletedAt existe en el modelo
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.VaultItem{}).Error
}
