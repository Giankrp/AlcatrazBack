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

// Package routes wires up all HTTP routes and applies JWT middleware to protected groups.
package routes

import (
	"log"
	"os"

	"github.com/Giankrp/AlcatrazBack/internal/handlers"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo, authHandler *handlers.AuthHandler, vaultHandler *handlers.VaultHandler, userProfileHandler *handlers.UserProfileHandler) {
	// API Group
	api := e.Group("/api")

	// Auth routes
	auth := api.Group("/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)
	auth.POST("/logout", authHandler.Logout)
	auth.GET("/exists", authHandler.UserExists)
	auth.POST("/2fa/verify", authHandler.Verify2FALogin)
	auth.POST("/recovery/fetch", authHandler.FetchRecoveryData)
	auth.POST("/recovery/reset", authHandler.ResetPassword)

	// Protected routes
	protected := api.Group("")
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}

	// Apply JWT middleware — read token from HttpOnly cookie
	protected.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey:  []byte(jwtSecret),
		TokenLookup: "cookie:auth_token",
	}))

	// User profile routes
	user := protected.Group("/user")
	user.GET("/profile", userProfileHandler.GetProfile)
	user.PUT("/profile", userProfileHandler.UpdateProfile)
	user.POST("/change-password", authHandler.ChangeMasterPassword)
	user.POST("/2fa/setup", authHandler.Setup2FA)
	user.POST("/2fa/enable", authHandler.Enable2FA)
	user.DELETE("/account", userProfileHandler.DeleteAccount)

	// Vault routes
	vault := protected.Group("/vault")
	vault.POST("/items", vaultHandler.CreateItem)
	vault.GET("/items", vaultHandler.GetItems)
	vault.GET("/trash", vaultHandler.GetTrash)
	vault.GET("/items/:id", vaultHandler.GetItem)
	vault.PUT("/items/:id", vaultHandler.UpdateItem)
	vault.DELETE("/items/:id", vaultHandler.MoveToTrash)
	vault.POST("/items/:id/restore", vaultHandler.RestoreItem)
	vault.DELETE("/items/:id/permanent", vaultHandler.DeleteItem)

	// Folder routes
	vault.POST("/folders", vaultHandler.CreateFolder)
	vault.GET("/folders", vaultHandler.GetFolders)
	vault.PUT("/folders/:id", vaultHandler.UpdateFolder)
	vault.DELETE("/folders/:id", vaultHandler.DeleteFolder)
}
