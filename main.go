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

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// 1. Initialize Echo
	e := echo.New()
	if err := godotenv.Load(); err != nil {
		e.Logger.Warn("Error loading .env file")
	}

	// 2. Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	allowedOriginsEnv := os.Getenv("ALLOWED_ORIGINS")
	var allowedOrigins []string
	if allowedOriginsEnv != "" {
		allowedOrigins = strings.Split(allowedOriginsEnv, ",")
	} else {
		allowedOrigins = []string{"http://localhost:3000"} // Fallback as warning or strictly to be defined
		e.Logger.Warn("ALLOWED_ORIGINS is not set. Defaulting to '*' which is insecure for production.")
	}

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     allowedOrigins,
		AllowCredentials: true,
	}))

	// 3. Custom Error Handler
	e.HTTPErrorHandler = customHTTPErrorHandler

	// 4. Database Connection
	database, err := db.NewConnection()
	if err != nil {
		e.Logger.Fatal("Error connecting to database: ", err)
	}

	// 5. Database Migration
	if err := db.AutoMigrate(database); err != nil {
		e.Logger.Fatal("Error migrating database: ", err)
	}

	// 6. Dependency Injection (Wiring)
	// Repositories
	userRepo := repositories.NewUserRepository(database)
	vaultRepo := repositories.NewVaultRepository(database)

	// Services
	authService := services.NewAuthService(userRepo)
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
