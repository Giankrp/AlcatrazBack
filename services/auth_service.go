package services

import (
	"errors"
	"os"
	"time"

	"github.com/Giankrp/AlcatrazBack/dto"
	"github.com/Giankrp/AlcatrazBack/models"
	"github.com/Giankrp/AlcatrazBack/repositories"
	"github.com/Giankrp/AlcatrazBack/security"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type AuthService interface {
	Register(registerDTO dto.RegisterDTO) error
	Login(loginDTO dto.LoginDTO) (string, error)
}

type authService struct {
	userRepo repositories.UserRepository
}

func NewAuthService(userRepo repositories.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

func (s *authService) Register(registerDTO dto.RegisterDTO) error {
	// Check if user exists
	_, err := s.userRepo.FindByEmail(registerDTO.Email)
	if err == nil {
		return errors.New("email already registered")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// Hash password
	hashedPassword, err := security.HashPassword(registerDTO.Password)
	if err != nil {
		return err
	}

	user := &models.User{
		Email:        registerDTO.Email,
		PasswordHash: hashedPassword,
		CreatedAt:    time.Now(),
	}

	return s.userRepo.Create(user)
}

func (s *authService) Login(loginDTO dto.LoginDTO) (string, error) {
	user, err := s.userRepo.FindByEmail(loginDTO.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("invalid credentials")
		}
		return "", err
	}

	match, err := security.VerifyPassword(loginDTO.Password, user.PasswordHash)
	if err != nil {
		return "", err
	}
	if !match {
		return "", errors.New("invalid credentials")
	}

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "default_secret_change_me" // Fallback only for dev
	}

	return token.SignedString([]byte(secret))
}
