// Package handlers contains the HTTP controllers (handlers) for the Alcatraz API.
// Each handler is responsible for parsing, validating, and delegating requests
// to the appropriate service layer, then formatting the HTTP response.
package handlers

import (
	"net/http"
	"os"
	"time"

	"github.com/Giankrp/AlcatrazBack/internal/dto"
	"github.com/Giankrp/AlcatrazBack/internal/services"
	"github.com/Giankrp/AlcatrazBack/internal/validator"
	"github.com/charmbracelet/log"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// AuthHandler manages HTTP requests related to authentication and user management.
type AuthHandler struct {
	authService services.AuthService
}

// NewAuthHandler creates a new instance of the authentication controller.
func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register handles the creation of new users.
func (h *AuthHandler) Register(c echo.Context) error {
	var register dto.RegisterDTO
	if err := c.Bind(&register); err != nil {
		log.Debug("Register: invalid request format", "error", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request format"})
	}

	if err := validator.Validate.Struct(&register); err != nil {
		log.Debug("Register: validation error", "error", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": validator.ValidationErrors(err)})
	}

	log.Debug("Register attempt", "email", register.Email)
	if err := h.authService.Register(register); err != nil {
		log.Warn("Registration failed", "error", err, "email", register.Email)
		if err.Error() == "email already registered" {
			return c.JSON(http.StatusConflict, echo.Map{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "internal server error"})
	}

	log.Info("User registered successfully", "email", register.Email)
	return c.JSON(http.StatusCreated, echo.Map{"message": "user created successfully"})
}

func (h *AuthHandler) Login(c echo.Context) error {
	var login dto.LoginDTO
	if err := c.Bind(&login); err != nil {
		log.Debug("Login: invalid request format", "error", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request format"})
	}

	if err := validator.Validate.Struct(&login); err != nil {
		log.Debug("Login: validation error", "error", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": validator.ValidationErrors(err)})
	}

	log.Debug("Login attempt", "email", login.Email)
	loginRes, err := h.authService.Login(login)
	if err != nil {
		log.Warn("Login failed", "error", err, "email", login.Email)
		if err.Error() == "invalid credentials" {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid email or password"})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "internal server error"})
	}

	if loginRes.Require2FA {
		log.Info("Login requires 2FA", "email", login.Email)
		return c.JSON(http.StatusOK, echo.Map{
			"require_2fa": true,
			"temp_token":  loginRes.Token,
		})
	}

	// Set JWT as HttpOnly cookie
	isProd := os.Getenv("ENV") == "production"
	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    loginRes.Token,
		Path:     "/",
		HttpOnly: true,
		Secure:   isProd,
		SameSite: http.SameSiteNoneMode,
		MaxAge:   60 * 60 * 12, // 12 hours
	}
	c.SetCookie(cookie)

	log.Info("Login successful", "email", login.Email)
	return c.JSON(http.StatusOK, echo.Map{
		"message":              "login successful",
		"protected_master_key": loginRes.ProtectedMasterKey,
		"master_key_iv":        loginRes.MasterKeyIV,
		"master_key_salt":      loginRes.MasterKeySalt,
	})
}

func (h *AuthHandler) ChangeMasterPassword(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userIDStr := claims["user_id"].(string)
	userID, _ := uuid.Parse(userIDStr)

	var input dto.ChangeMasterPasswordDTO
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request format"})
	}

	if err := validator.Validate.Struct(&input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": validator.ValidationErrors(err)})
	}

	if err := h.authService.ChangeMasterPassword(userID, input); err != nil {
		if err.Error() == "invalid current password" {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "master password updated successfully"})
}

func (h *AuthHandler) Setup2FA(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userIDStr := claims["user_id"].(string)
	userID, _ := uuid.Parse(userIDStr)

	res, err := h.authService.Generate2FASecret(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to generate 2FA secret"})
	}

	return c.JSON(http.StatusOK, res)
}

func (h *AuthHandler) Enable2FA(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userIDStr := claims["user_id"].(string)
	userID, _ := uuid.Parse(userIDStr)

	var enableDTO dto.Enable2FADTO
	if err := c.Bind(&enableDTO); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request format"})
	}

	backupCodes, err := h.authService.Enable2FA(userID, enableDTO)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message":      "2FA enabled successfully",
		"backup_codes": backupCodes,
	})
}

func (h *AuthHandler) Verify2FALogin(c echo.Context) error {
	var verify dto.Verify2FADTO
	if err := c.Bind(&verify); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request format"})
	}

	tempToken := c.Request().Header.Get("X-Temp-Token")
	if tempToken == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "missing temporary token"})
	}

	// Parse temp token
	token, err := jwt.Parse(tempToken, func(token *jwt.Token) (any, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid or expired temporary token"})
	}

	claims := token.Claims.(jwt.MapClaims)
	if claims["pending_2fa"] != true {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid token type"})
	}

	userIDStr := claims["user_id"].(string)
	userID, _ := uuid.Parse(userIDStr)

	loginRes, err := h.authService.Verify2FALogin(userID, verify.Code)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": err.Error()})
	}

	// Set final JWT as HttpOnly cookie
	isProd := os.Getenv("ENV") == "production"
	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    loginRes.Token,
		Path:     "/",
		HttpOnly: true,
		Secure:   isProd,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   60 * 60 * 12, // 12 hours
	}
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, echo.Map{
		"message":              "login successful",
		"protected_master_key": loginRes.ProtectedMasterKey,
		"master_key_iv":        loginRes.MasterKeyIV,
		"master_key_salt":      loginRes.MasterKeySalt,
	})
}

func (h *AuthHandler) Logout(c echo.Context) error {
	log.Debug("Logout attempt")
	// Expire the auth cookie
	isProd := os.Getenv("ENV") == "production"
	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   isProd,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
	}
	c.SetCookie(cookie)

	log.Info("Logout successful")
	return c.JSON(http.StatusOK, echo.Map{"message": "logged out successfully"})
}

func (h *AuthHandler) UserExists(c echo.Context) error {
	email := c.QueryParam("email")
	if email == "" {
		log.Debug("UserExists: email query parameter is required")
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "email query parameter is required"})
	}

	log.Debug("Checking if user exists", "email", email)
	exists, err := h.authService.UserExists(email)
	if err != nil {
		log.Error("Error checking if user exists", "error", err, "email", email)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, echo.Map{"exists": exists})
}

func (h *AuthHandler) FetchRecoveryData(c echo.Context) error {
	var input dto.FetchRecoveryDTO
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request format"})
	}

	if err := validator.Validate.Struct(&input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": validator.ValidationErrors(err)})
	}

	user, err := h.authService.FetchRecoveryData(input.Email)
	if err != nil {
		if err.Error() == "user not found" {
			// To prevent email enumeration during recovery, we could return 200 with fake data,
			return c.JSON(http.StatusNotFound, echo.Map{"error": "user not found"})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"recovery_protected_master_key": user.RecoveryProtectedMasterKey,
		"recovery_key_iv":               user.RecoveryKeyIV,
		"recovery_key_salt":             user.RecoveryKeySalt,
	})
}

func (h *AuthHandler) ResetPassword(c echo.Context) error {
	var input dto.ResetPasswordDTO
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request format"})
	}

	if err := validator.Validate.Struct(&input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": validator.ValidationErrors(err)})
	}

	if err := h.authService.ResetPasswordWithRecoveryKey(input); err != nil {
		if err.Error() == "user not found" || err.Error() == "invalid recovery key" {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid recovery key or email"})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "password reset successfully"})
}
