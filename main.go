package main

import (
	"net/http"
	"os"
	"strings"

	"github.com/Giankrp/AlcatrazBack/db"
	"github.com/Giankrp/AlcatrazBack/handlers"
	"github.com/Giankrp/AlcatrazBack/repositories"
	"github.com/Giankrp/AlcatrazBack/routes"
	"github.com/Giankrp/AlcatrazBack/services"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func init() {
	// Setup charmbracelet/log for nice debugging output
	log.SetLevel(log.DebugLevel)
	log.SetReportCaller(true)
}

func main() {
	// 1. Initialize Echo
	e := echo.New()
	if err := godotenv.Load(); err != nil {
		log.Warn("Error loading .env file, relying on environment variables")
	}

	// 2. Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	allowedOriginsEnv := os.Getenv("ALLOWED_ORIGINS")
	var allowedOrigins []string
	if allowedOriginsEnv != "" {
		allowedOrigins = strings.Split(allowedOriginsEnv, ",")
		log.Debug("Allowed origins configured", "origins", allowedOrigins)
	} else {
		allowedOrigins = []string{"http://localhost:3000"} // Fallback as warning or strictly to be defined
		log.Warn("ALLOWED_ORIGINS is not set. Defaulting to '*' which is insecure for production.")
	}

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     allowedOrigins,
		AllowCredentials: true,
	}))

	// 3. Custom Error Handler
	e.HTTPErrorHandler = customHTTPErrorHandler

	// 4. Database Connection
	log.Info("Connecting to database...")
	database, err := db.NewConnection()
	if err != nil {
		log.Fatal("Error connecting to database", "error", err)
	}
	log.Info("Successfully connected to database")

	// 5. Database Migration
	log.Info("Running database migrations...")
	if err := db.AutoMigrate(database); err != nil {
		log.Fatal("Error migrating database", "error", err)
	}
	log.Info("Database migrations completed")

	// 6. Dependency Injection (Wiring)
	// Repositories
	userRepo := repositories.NewUserRepository(database)
	vaultRepo := repositories.NewVaultRepository(database)

	// Services
	authService := services.NewAuthService(userRepo, vaultRepo)
	vaultService := services.NewVaultService(vaultRepo)
	userService := services.NewUserService(userRepo)

	// Handlers
	authHandler := handlers.NewAuthHandler(authService)
	vaultHandler := handlers.NewVaultHandler(vaultService)
	userProfileHandler := handlers.NewUserProfileHandler(userService)

	// 7. Routes
	routes.SetupRoutes(e, authHandler, vaultHandler, userProfileHandler)

	// 8. Start Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Infof("Starting server on port %s", port)
	log.Fatal("Server stopped", "error", e.Start(":"+port))
}

func customHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	message := "Internal Server Error"

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		message = he.Message.(string)
	}

	// Log the actual error for debugging
	log.Error("HTTP Error", "method", c.Request().Method, "path", c.Request().URL.Path, "status", code, "error", err.Error())

	// Send JSON response
	if !c.Response().Committed {
		if c.Request().Method == http.MethodHead {
			err = c.NoContent(code)
		} else {
			err = c.JSON(code, echo.Map{"error": message})
		}
		if err != nil {
			log.Error("Error sending error response", "error", err)
		}
	}
}
