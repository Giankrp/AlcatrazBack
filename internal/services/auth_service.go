package services

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/Giankrp/AlcatrazBack/internal/dto"
	"github.com/Giankrp/AlcatrazBack/internal/models"
	"github.com/Giankrp/AlcatrazBack/internal/repositories"
	"github.com/Giankrp/AlcatrazBack/internal/security"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/pquerna/otp/totp"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// AuthService manages authentication, registration, and user security (2FA/MasterKey).
type AuthService interface {
	// Register creates a new user with their encrypted MasterKey.
	Register(registerDTO dto.RegisterDTO) error

	// Login verifies credentials and returns the MasterKey if successful.
	Login(loginDTO dto.LoginDTO) (*dto.LoginResponseDTO, error)

	// UserExists checks if an email is already registered in the system.
	UserExists(email string) (bool, error)

	// Generate2FASecret generates a new seed for two-factor authentication.
	Generate2FASecret(userID uuid.UUID) (*dto.Setup2FAResponseDTO, error)

	// Enable2FA activates two-factor authentication for a user.
	Enable2FA(userID uuid.UUID, enableDTO dto.Enable2FADTO) ([]string, error)

	// Verify2FALogin validates the 2FA code during login.
	Verify2FALogin(userID uuid.UUID, code string) (string, error)

	// ChangeMasterPassword atomically updates the password and the protected MasterKey.
	ChangeMasterPassword(userID uuid.UUID, input dto.ChangeMasterPasswordDTO) error
}

type authService struct {
	userRepo  repositories.UserRepository
	vaultRepo repositories.VaultRepository
}

func NewAuthService(userRepo repositories.UserRepository, vaultRepo repositories.VaultRepository) AuthService {
	return &authService{userRepo: userRepo, vaultRepo: vaultRepo}
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
		Email:              registerDTO.Email,
		PasswordHash:       hashedPassword,
		ProtectedMasterKey: registerDTO.ProtectedMasterKey,
		MasterKeyIV:        registerDTO.MasterKeyIV,
		MasterKeySalt:      registerDTO.MasterKeySalt,
		CreatedAt:          time.Now(),
	}

	if err := s.userRepo.Create(user); err != nil {
		return err
	}

	profile := &models.UserProfile{
		UserID: user.ID,
		Name:   name,
	}
	if err := s.userRepo.CreateProfile(profile); err != nil {
		return err
	}

	// Create default "Personal" folder
	defaultFolder := &models.VaultFolder{
		UserID:    user.ID,
		Name:      "Personal",
		IsDefault: true,
	}
	return s.vaultRepo.CreateFolder(defaultFolder)
}

func (s *authService) Login(loginDTO dto.LoginDTO) (*dto.LoginResponseDTO, error) {
	user, err := s.userRepo.FindByEmail(loginDTO.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}

	match, err := security.VerifyPassword(loginDTO.Password, user.PasswordHash)
	if err != nil {
		return nil, err
	}
	if !match {
		return nil, errors.New("invalid credentials")
	}

	if user.TwoFactorEnabled {
		// Generate a short-lived "pending 2FA" token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id":     user.ID.String(),
			"pending_2fa": true,
			"exp":         time.Now().Add(time.Minute * 10).Unix(),
		})

		secret := os.Getenv("JWT_SECRET")
		signedToken, err := token.SignedString([]byte(secret))
		if err != nil {
			return nil, err
		}

		return &dto.LoginResponseDTO{
			Require2FA: true,
			Token:      signedToken,
		}, nil
	}

	// Generate normal session JWT
	tokenString, err := s.generateSessionToken(user)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponseDTO{
		Require2FA:         false,
		Token:              tokenString,
		ProtectedMasterKey: user.ProtectedMasterKey,
		MasterKeyIV:        user.MasterKeyIV,
		MasterKeySalt:      user.MasterKeySalt,
	}, nil
}

func (s *authService) generateSessionToken(user *models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID.String(),
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 12).Unix(),
	})

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("JWT_SECRET environment variable is not set")
	}

	return token.SignedString([]byte(secret))
}

func (s *authService) Generate2FASecret(userID uuid.UUID) (*dto.Setup2FAResponseDTO, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Alcatraz",
		AccountName: user.Email,
	})
	if err != nil {
		return nil, err
	}

	return &dto.Setup2FAResponseDTO{
		Secret: key.Secret(),
		QRURI:  key.URL(),
	}, nil
}

func (s *authService) Enable2FA(userID uuid.UUID, enableDTO dto.Enable2FADTO) ([]string, error) {
	valid := totp.Validate(enableDTO.Code, enableDTO.Secret)
	if !valid {
		return nil, errors.New("invalid verification code")
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	// Generate 8 backup codes
	backupCodes := make([]string, 8)
	for i := range 8 {
		code, err := security.GenerateRandomString(8)
		if err != nil {
			return nil, err
		}
		backupCodes[i] = code
	}

	// Persist changes
	user.TwoFactorEnabled = true
	user.TwoFactorSecret = enableDTO.Secret
	backupCodesBytes, _ := json.Marshal(backupCodes)
	user.TwoFactorBackupCodes = datatypes.JSON(backupCodesBytes)

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return backupCodes, nil
}

func (s *authService) Verify2FALogin(userID uuid.UUID, code string) (string, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return "", err
	}

	if !user.TwoFactorEnabled {
		return "", errors.New("2FA is not enabled for this user")
	}

	// Try TOTP first
	valid := totp.Validate(code, user.TwoFactorSecret)
	if !valid {
		// Try backup codes
		var backupCodes []string
		_ = json.Unmarshal(user.TwoFactorBackupCodes, &backupCodes)

		found := -1
		for i, bc := range backupCodes {
			if bc == code {
				found = i
				break
			}
		}

		if found == -1 {
			return "", errors.New("invalid verification code")
		}

		// Remove the used backup code
		backupCodes = append(backupCodes[:found], backupCodes[found+1:]...)
		backupCodesBytes, _ := json.Marshal(backupCodes)
		user.TwoFactorBackupCodes = datatypes.JSON(backupCodesBytes)
		if err := s.userRepo.Update(user); err != nil {
			return "", err
		}
	}

	return s.generateSessionToken(user)
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

func (s *authService) ChangeMasterPassword(userID uuid.UUID, input dto.ChangeMasterPasswordDTO) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	// 1. Verify old password
	match, err := security.VerifyPassword(input.OldPassword, user.PasswordHash)
	if err != nil {
		return err
	}
	if !match {
		return errors.New("invalid current password")
	}

	// 2. Hash new password
	newPasswordHash, err := security.HashPassword(input.NewPassword)
	if err != nil {
		return err
	}

	// 3. Update all security metadata
	user.PasswordHash = newPasswordHash
	user.ProtectedMasterKey = input.ProtectedMasterKey
	user.MasterKeyIV = input.MasterKeyIV
	user.MasterKeySalt = input.MasterKeySalt

	return s.userRepo.Update(user)
}
