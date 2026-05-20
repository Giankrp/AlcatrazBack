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

	"github.com/Giankrp/AlcatrazBack/internal/dto"
	"github.com/Giankrp/AlcatrazBack/internal/services"
	"github.com/Giankrp/AlcatrazBack/internal/validator"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// VaultHandler manages HTTP requests for vault items and folders.
type VaultHandler struct {
	service services.VaultService
}

// NewVaultHandler creates a new instance of the vault controller.
func NewVaultHandler(service services.VaultService) *VaultHandler {
	return &VaultHandler{service: service}
}

// CreateItem creates a new encrypted vault item for the authenticated user.
// The encrypted payload is accepted as-is and stored without inspection (Zero Knowledge).
func (h *VaultHandler) CreateItem(c echo.Context) error {
	userID := getUserIDFromToken(c)
	if userID == uuid.Nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}
	var input dto.CreateVaultItemDTO
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	if err := validator.Validate.Struct(input); err != nil {
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

// GetItems returns all active (non-trashed) vault items for the authenticated user.
// Encrypted secrets are NOT included in list responses for performance reasons.
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

// GetTrash returns all vault items currently in the trash for the authenticated user.
func (h *VaultHandler) GetTrash(c echo.Context) error {
	userID := getUserIDFromToken(c)
	if userID == uuid.Nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}
	items, err := h.service.GetTrashedItems(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch trash items"})
	}
	return c.JSON(http.StatusOK, items)
}

// GetItem returns a specific vault item by ID, including its encrypted secret payload.
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

// UpdateItem updates a vault item's metadata and/or its encrypted secret payload.
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

	if err := validator.Validate.Struct(input); err != nil {
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

// MoveToTrash performs a soft-delete on a vault item, moving it to the trash.
func (h *VaultHandler) MoveToTrash(c echo.Context) error {
	userID := getUserIDFromToken(c)
	if userID == uuid.Nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}
	itemID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid item ID"})
	}

	if err := h.service.MoveToTrash(userID, itemID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to move item to trash"})
	}

	return c.NoContent(http.StatusNoContent)
}

// RestoreItem recovers a previously trashed vault item, making it active again.
func (h *VaultHandler) RestoreItem(c echo.Context) error {
	userID := getUserIDFromToken(c)
	if userID == uuid.Nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}
	itemID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid item ID"})
	}

	if err := h.service.RestoreFromTrash(userID, itemID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to restore item"})
	}

	return c.NoContent(http.StatusNoContent)
}

// DeleteItem permanently removes a vault item and its associated secret from the database.
// This operation is irreversible. The item must be in the trash before it can be deleted permanently.
func (h *VaultHandler) DeleteItem(c echo.Context) error {
	userID := getUserIDFromToken(c)
	if userID == uuid.Nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}
	itemID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid item ID"})
	}

	if err := h.service.PermanentlyDelete(userID, itemID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete item permanently"})
	}

	return c.NoContent(http.StatusNoContent)
}

// Folder Handlers

// CreateFolder creates a new folder for organizing vault items.
func (h *VaultHandler) CreateFolder(c echo.Context) error {
	userID := getUserIDFromToken(c)
	if userID == uuid.Nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	var input dto.CreateVaultFolderDTO
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	if err := validator.Validate.Struct(input); err != nil {
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

// GetFolders returns all folders belonging to the authenticated user.
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

// UpdateFolder modifies an existing folder's metadata (e.g., name).
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

	if err := validator.Validate.Struct(input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	folder, err := h.service.UpdateFolder(userID, id, input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update folder"})
	}

	return c.JSON(http.StatusOK, folder)
}

// DeleteFolder deletes a folder. Its items are automatically reassigned to the default (Personal) folder.
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

// getUserIDFromToken is a helper to extract the user ID from the JWT token.
func getUserIDFromToken(c echo.Context) uuid.UUID {
	// Check if the context has the user
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
