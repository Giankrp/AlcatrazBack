package handlers

import (
	"net/http"
	"time"

	"github.com/Giankrp/AlcatrazBack/dto"
	"github.com/Giankrp/AlcatrazBack/services"
	"github.com/Giankrp/AlcatrazBack/validator"
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
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request format"})
	}

	if err := validator.Validate.Struct(&register); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": validator.ValidationErrors(err)})
	}

	if err := h.authService.Register(register); err != nil {
		// En un caso real, chequear tipo de error para devolver 409 Conflict si ya existe, etc.
		// Por simplicidad, 400 o 500 según corresponda.
		if err.Error() == "email already registered" {
			return c.JSON(http.StatusConflict, echo.Map{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "internal server error"})
	}

	return c.JSON(http.StatusCreated, echo.Map{"message": "user created successfully"})
}

func (h *AuthHandler) Login(c echo.Context) error {
	var login dto.LoginDTO
	if err := c.Bind(&login); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request format"})
	}

	if err := validator.Validate.Struct(&login); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": validator.ValidationErrors(err)})
	}

	token, err := h.authService.Login(login)
	if err != nil {
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

	return c.JSON(http.StatusOK, echo.Map{"message": "login successful"})
}

func (h *AuthHandler) Logout(c echo.Context) error {
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

	return c.JSON(http.StatusOK, echo.Map{"message": "logged out successfully"})
}
