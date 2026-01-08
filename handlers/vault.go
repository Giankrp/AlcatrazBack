package handlers

import (
	"net/http"

	"github.com/Giankrp/AlcatrazBack/dto"
	"github.com/Giankrp/AlcatrazBack/services"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type VaultHandler struct {
	service   services.VaultService
	validator *validator.Validate
}

func NewVaultHandler(service services.VaultService) *VaultHandler {
	return &VaultHandler{
		service:   service,
		validator: validator.New(),
	}
}

func (h *VaultHandler) CreateItem(c echo.Context) error {
	userID := getUserIDFromToken(c)
	var input dto.CreateVaultItemDTO
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	if err := h.validator.Struct(input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	item, err := h.service.CreateItem(userID, input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create item"})
	}

	return c.JSON(http.StatusCreated, item)
}

func (h *VaultHandler) GetItems(c echo.Context) error {
	userID := getUserIDFromToken(c)
	items, err := h.service.GetItems(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch items"})
	}
	return c.JSON(http.StatusOK, items)
}

func (h *VaultHandler) GetItem(c echo.Context) error {
	userID := getUserIDFromToken(c)
	itemID := c.Param("id")
	item, err := h.service.GetItem(userID, itemID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Item not found"})
	}
	return c.JSON(http.StatusOK, item)
}

func (h *VaultHandler) UpdateItem(c echo.Context) error {
	userID := getUserIDFromToken(c)
	itemID := c.Param("id")

	var input dto.UpdateVaultItemDTO
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	if err := h.validator.Struct(input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	item, err := h.service.UpdateItem(userID, itemID, input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update item"})
	}

	return c.JSON(http.StatusOK, item)
}

func (h *VaultHandler) DeleteItem(c echo.Context) error {
	userID := getUserIDFromToken(c)
	itemID := c.Param("id")

	if err := h.service.DeleteItem(userID, itemID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete item"})
	}

	return c.NoContent(http.StatusNoContent)
}

// Helper para extraer ID del token JWT
func getUserIDFromToken(c echo.Context) string {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	return claims["user_id"].(string)
}
