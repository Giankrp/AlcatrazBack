package main

import (
	"github.com/Giankrp/AlcatrazBack/db"

	"net/http"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	err := godotenv.Load()

	e := echo.New()
	if err != nil {
		e.Logger.Fatal("Error loading .env file")
	}

	err = db.DbConnection()
	if err != nil {
		e.Logger.Fatal("Error connecting to database")
	}

	err = db.AutoMigrate()
	if err != nil {
		e.Logger.Fatal("Error migrating database")
	}

	e.Use(middleware.Logger())
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":8080"))
}
