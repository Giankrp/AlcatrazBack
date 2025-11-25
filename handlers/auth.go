package handlers

import (
	"net/http"

	"github.com/Giankrp/AlcatrazBack/dto"
	"github.com/Giankrp/AlcatrazBack/validator"

	"github.com/labstack/echo/v4"
)

func Register(c echo.Context) error {
	var register dto.RegisterDTO
	if err := c.Bind(&register); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request"})
	}

	if err := validator.Validate.Struct(&register); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": validator.ValidationErrors(err)})
	}

	return nil
}
