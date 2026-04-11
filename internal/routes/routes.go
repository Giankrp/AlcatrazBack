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
	user.POST("/2fa/setup", authHandler.Setup2FA)
	user.POST("/2fa/enable", authHandler.Enable2FA)

	// Vault routes
	vault := protected.Group("/vault")
	vault.POST("/items", vaultHandler.CreateItem)
	vault.GET("/items", vaultHandler.GetItems)
	vault.GET("/items/:id", vaultHandler.GetItem)
	vault.PUT("/items/:id", vaultHandler.UpdateItem)
	vault.DELETE("/items/:id", vaultHandler.DeleteItem)

	// Folder routes
	vault.POST("/folders", vaultHandler.CreateFolder)
	vault.GET("/folders", vaultHandler.GetFolders)
	vault.PUT("/folders/:id", vaultHandler.UpdateFolder)
	vault.DELETE("/folders/:id", vaultHandler.DeleteFolder)
}
