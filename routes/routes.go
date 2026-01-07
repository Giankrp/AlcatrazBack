package routes

import (
	"github.com/Giankrp/AlcatrazBack/handlers"
	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo, authHandler *handlers.AuthHandler) {
	// API Group
	api := e.Group("/api")

	// Auth routes
	auth := api.Group("/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)
}
