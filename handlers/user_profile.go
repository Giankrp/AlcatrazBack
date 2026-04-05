package handlers

import (
	"net/http"

	"github.com/Giankrp/AlcatrazBack/dto"
	"github.com/Giankrp/AlcatrazBack/models"
	"github.com/Giankrp/AlcatrazBack/services"
	"github.com/Giankrp/AlcatrazBack/validator"
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

	return c.JSON(http.StatusOK, profile)
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
