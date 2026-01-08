package routes

import (
	"os"

	"github.com/Giankrp/AlcatrazBack/handlers"
	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo, authHandler *handlers.AuthHandler, vaultHandler *handlers.VaultHandler) {
	// API Group
	api := e.Group("/api")

	// Auth routes
	auth := api.Group("/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)

	// Protected routes
	//protected := api.Group("")
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "secret" // Fallback para dev, idealmente fatal error en prod
	}

	// MODO DESARROLLO: Middleware JWT desactivado temporalmente para pruebas fáciles
	// protected.Use(echojwt.WithConfig(echojwt.Config{
	// 	SigningKey: []byte(jwtSecret),
	// }))

	// Vault routes
	vault := e.Group("/vault")
	vault.POST("/items", vaultHandler.CreateItem)
	vault.GET("/items", vaultHandler.GetItems)
	vault.GET("/items/:id", vaultHandler.GetItem)
	vault.PUT("/items/:id", vaultHandler.UpdateItem)
	vault.DELETE("/items/:id", vaultHandler.DeleteItem)
}
