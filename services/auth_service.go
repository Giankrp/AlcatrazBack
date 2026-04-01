package services

import (
	"errors"
	"os"
	"strings"
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
	UserExists(email string) (bool, error)
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
	name := strings.Split(registerDTO.Email, "@")[0]

	user := &models.User{
		Email:        registerDTO.Email,
		PasswordHash: hashedPassword,
		CreatedAt:    time.Now(),
	}

	if err := s.userRepo.Create(user); err != nil {
		return err
	}

	profile := &models.UserProfile{
		UserID: user.ID,
		Name:   name,
	}
	return s.userRepo.CreateProfile(profile)
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
		"exp":     time.Now().Add(time.Hour * 12).Unix(),
	})

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("JWT_SECRET environment variable is not set")
	}

	return token.SignedString([]byte(secret))
}

func (s *authService) UserExists(email string) (bool, error) {
	_, err := s.userRepo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
