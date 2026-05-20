// Alcatraz - Secure open source Password Manager and Storage System
// Copyright (C) 2026 Gian Carlo Ruiz Patiño
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

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
