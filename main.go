package main

import (
	"net/http"
	"os"

	"github.com/Giankrp/AlcatrazBack/db"
	"github.com/Giankrp/AlcatrazBack/handlers"
	"github.com/Giankrp/AlcatrazBack/repositories"
	"github.com/Giankrp/AlcatrazBack/routes"
	"github.com/Giankrp/AlcatrazBack/services"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// 1. Load configuration
	if err := godotenv.Load(); err != nil {
		// Log but continue (env vars might be set in system)
		// e.Logger.Warn("Error loading .env file")
	}

	// 2. Initialize Echo
	e := echo.New()

	// 3. Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// 4. Custom Error Handler
	e.HTTPErrorHandler = customHTTPErrorHandler

	// 5. Database Connection
	database, err := db.NewConnection()
	if err != nil {
		e.Logger.Fatal("Error connecting to database: ", err)
	}

	// 6. Database Migration
	if err := db.AutoMigrate(database); err != nil {
		e.Logger.Fatal("Error migrating database: ", err)
	}

	// 7. Dependency Injection (Wiring)
	// Repositories
	userRepo := repositories.NewUserRepository(database)
	vaultRepo := repositories.NewVaultRepository(database)

	// Services
	authService := services.NewAuthService(userRepo)
	vaultService := services.NewVaultService(vaultRepo)

	// Handlers
	authHandler := handlers.NewAuthHandler(authService)
	vaultHandler := handlers.NewVaultHandler(vaultService)

	// 8. Routes
	routes.SetupRoutes(e, authHandler, vaultHandler)

	// 9. Start Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	e.Logger.Fatal(e.Start(":" + port))
}

func customHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	message := "Internal Server Error"

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		message = he.Message.(string)
	}

	// Send JSON response
	if !c.Response().Committed {
		if c.Request().Method == http.MethodHead {
			err = c.NoContent(code)
		} else {
			err = c.JSON(code, echo.Map{"error": message})
		}
		if err != nil {
			c.Logger().Error(err)
		}
	}
}
