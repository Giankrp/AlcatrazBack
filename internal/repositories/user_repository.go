package repositories

import (
	"github.com/Giankrp/AlcatrazBack/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
	FindByEmail(email string) (*models.User, error)
	CreateProfile(profile *models.UserProfile) error
	FindProfileByUserID(userID uuid.UUID) (*models.UserProfile, error)
	UpdateProfile(profile *models.UserProfile) error
	FindByID(id uuid.UUID) (*models.User, error)
	Update(user *models.User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) CreateProfile(profile *models.UserProfile) error {
	return r.db.Create(profile).Error
}

func (r *userRepository) FindProfileByUserID(userID uuid.UUID) (*models.UserProfile, error) {
	var profile models.UserProfile
	if err := r.db.Preload("User").Where("user_id = ?", userID).First(&profile).Error; err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *userRepository) UpdateProfile(profile *models.UserProfile) error {
	return r.db.Save(profile).Error
}

func (r *userRepository) FindByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}
