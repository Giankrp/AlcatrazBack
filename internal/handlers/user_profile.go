package handlers

import (
	"net/http"
	"os"
	"time"

	"github.com/Giankrp/AlcatrazBack/internal/dto"
	"github.com/Giankrp/AlcatrazBack/internal/models"
	"github.com/Giankrp/AlcatrazBack/internal/services"
	"github.com/Giankrp/AlcatrazBack/internal/validator"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type UserProfileHandler struct {
	userProfileService services.UserService
}

func NewUserProfileHandler(userProfileService services.UserService) *UserProfileHandler {
	return &UserProfileHandler{userProfileService: userProfileService}
}

func (u *UserProfileHandler) GetProfile(c echo.Context) error {
	userID := getUserIDFromToken(c)
	if userID == uuid.Nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Unauthorized"})
	}

	profile, err := u.userProfileService.GetProfile(userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Profile not found"})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"user_id":            profile.UserID,
		"name":               profile.Name,
		"avatar_url":         profile.AvatarURL,
		"language":           profile.Language,
		"created_at":         profile.CreatedAt,
		"updated_at":         profile.UpdatedAt,
		"two_factor_enabled": profile.User.TwoFactorEnabled,
	})
}

func (u *UserProfileHandler) UpdateProfile(c echo.Context) error {
	var update dto.UpdateUserProfileDTO

	userID := getUserIDFromToken(c)

	if userID == uuid.Nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Unauthorized"})
	}
	if err := c.Bind(&update); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request format"})
	}

	if err := validator.Validate.Struct(update); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": validator.ValidationErrors(err)})
	}

	profile := models.UserProfile{UserID: userID}

	if update.Name != "" {
		profile.Name = update.Name
	}

	if update.AvatarURL != "" {
		profile.AvatarURL = update.AvatarURL
	}

	if update.Language != "" {
		profile.Language = update.Language
	}

	u.userProfileService.UpdateProfile(&profile)
	return c.JSON(http.StatusOK, echo.Map{"Message": "User updated"})

}

func (u *UserProfileHandler) DeleteAccount(c echo.Context) error {
	userID := getUserIDFromToken(c)
	if userID == uuid.Nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Unauthorized"})
	}

	if err := u.userProfileService.DeleteAccount(userID); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to delete account"})
	}

	// Delete cookie (logout)
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

	return c.JSON(http.StatusOK, echo.Map{"message": "Account successfully deleted"})
}
