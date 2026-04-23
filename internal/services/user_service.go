package services

import (
	"github.com/Giankrp/AlcatrazBack/internal/models"
	"github.com/Giankrp/AlcatrazBack/internal/repositories"
	"github.com/google/uuid"
)

type UserService interface {
	GetProfile(userID uuid.UUID) (*models.UserProfile, error)
	UpdateProfile(profile *models.UserProfile) error
	DeleteAccount(userID uuid.UUID) error
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) GetProfile(userID uuid.UUID) (*models.UserProfile, error) {
	return s.userRepo.FindProfileByUserID(userID)
}

func (s *userService) UpdateProfile(profile *models.UserProfile) error {

	return s.userRepo.UpdateProfile(profile)
}

func (s *userService) DeleteAccount(userID uuid.UUID) error {
	return s.userRepo.Delete(userID)
}
