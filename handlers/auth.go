package handlers

import (
	"net/http"
	"time"

	"github.com/Giankrp/AlcatrazBack/dto"
	"github.com/Giankrp/AlcatrazBack/services"
	"github.com/Giankrp/AlcatrazBack/validator"
	"github.com/charmbracelet/log"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

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
		// En un caso real, chequear tipo de error para devolver 409 Conflict si ya existe, etc.
		// Por simplicidad, 400 o 500 según corresponda.
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
	token, err := h.authService.Login(login)
	if err != nil {
		log.Warn("Login failed", "error", err, "email", login.Email)
		if err.Error() == "invalid credentials" {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid email or password"})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "internal server error"})
	}

	// Set JWT as HttpOnly cookie
	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   60 * 60 * 12, // 12 hours
	}
	c.SetCookie(cookie)

	log.Info("Login successful", "email", login.Email)
	return c.JSON(http.StatusOK, echo.Map{"message": "login successful"})
}

func (h *AuthHandler) Logout(c echo.Context) error {
	log.Debug("Logout attempt")
	// Expire the auth cookie
	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
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
