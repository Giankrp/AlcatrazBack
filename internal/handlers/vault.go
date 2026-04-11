package handlers

import (
	"net/http"

	"github.com/Giankrp/AlcatrazBack/internal/dto"
	"github.com/Giankrp/AlcatrazBack/internal/services"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
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
	if userID == uuid.Nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	var input dto.CreateVaultItemDTO
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	if err := h.validator.Struct(input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	item, err := h.service.CreateItem(userID, input)
	if err != nil {
		if err.Error() == "folder not found or unauthorized" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create item"})
	}

	return c.JSON(http.StatusCreated, item)
}

func (h *VaultHandler) GetItems(c echo.Context) error {
	userID := getUserIDFromToken(c)
	if userID == uuid.Nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}
	items, err := h.service.GetItems(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch items"})
	}
	return c.JSON(http.StatusOK, items)
}

func (h *VaultHandler) GetItem(c echo.Context) error {
	userID := getUserIDFromToken(c)
	if userID == uuid.Nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}
	itemID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid item ID"})
	}
	item, err := h.service.GetItem(userID, itemID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Item not found"})
	}
	return c.JSON(http.StatusOK, item)
}

func (h *VaultHandler) UpdateItem(c echo.Context) error {
	userID := getUserIDFromToken(c)
	if userID == uuid.Nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}
	itemID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid item ID"})
	}

	var input dto.UpdateVaultItemDTO
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	if err := h.validator.Struct(input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	item, err := h.service.UpdateItem(userID, itemID, input)
	if err != nil {
		if err.Error() == "folder not found or unauthorized" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update item"})
	}

	return c.JSON(http.StatusOK, item)
}

func (h *VaultHandler) DeleteItem(c echo.Context) error {
	userID := getUserIDFromToken(c)
	if userID == uuid.Nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}
	itemID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid item ID"})
	}

	if err := h.service.DeleteItem(userID, itemID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete item"})
	}

	return c.NoContent(http.StatusNoContent)
}

// Folder Handlers

func (h *VaultHandler) CreateFolder(c echo.Context) error {
	userID := getUserIDFromToken(c)
	if userID == uuid.Nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	var input dto.CreateVaultFolderDTO
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	if err := h.validator.Struct(input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	folder, err := h.service.CreateFolder(userID, input)
	if err != nil {
		if err.Error() == "a folder with this name already exists" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create folder"})
	}

	return c.JSON(http.StatusCreated, folder)
}

func (h *VaultHandler) GetFolders(c echo.Context) error {
	userID := getUserIDFromToken(c)
	if userID == uuid.Nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	folders, err := h.service.GetFolders(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch folders"})
	}

	return c.JSON(http.StatusOK, folders)
}

func (h *VaultHandler) UpdateFolder(c echo.Context) error {
	userID := getUserIDFromToken(c)
	if userID == uuid.Nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid folder ID"})
	}

	var input dto.UpdateVaultFolderDTO
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	if err := h.validator.Struct(input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	folder, err := h.service.UpdateFolder(userID, id, input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update folder"})
	}

	return c.JSON(http.StatusOK, folder)
}

func (h *VaultHandler) DeleteFolder(c echo.Context) error {
	userID := getUserIDFromToken(c)
	if userID == uuid.Nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid folder ID"})
	}

	if err := h.service.DeleteFolder(userID, id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

// Helper para extraer ID del token JWT
func getUserIDFromToken(c echo.Context) uuid.UUID {
	// Verificar si el contexto tiene el usuario
	userToken, ok := c.Get("user").(*jwt.Token)
	if !ok || userToken == nil {
		return uuid.Nil
	}

	claims, ok := userToken.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return uuid.Nil
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil
	}

	return userID
}
