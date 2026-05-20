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

// UserProfileHandler manages HTTP requests for user profile operations and account management.
type UserProfileHandler struct {
	userProfileService services.UserService
}

// NewUserProfileHandler creates a new instance of the user profile controller.
func NewUserProfileHandler(userProfileService services.UserService) *UserProfileHandler {
	return &UserProfileHandler{userProfileService: userProfileService}
}

// GetProfile returns the authenticated user's profile data, including 2FA status.
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


// UpdateProfile updates editable profile fields (name, avatar URL, language).
// Only non-empty fields in the request body are updated (partial update).
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

	if err := u.userProfileService.UpdateProfile(&profile); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to update profile"})
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "profile updated successfully"})

}


// DeleteAccount permanently deletes the authenticated user's account and all associated data
// (vault items, secrets, folders, profile). It also expires the auth cookie on success.
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
