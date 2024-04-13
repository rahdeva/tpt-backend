package controllers

import (
	"net/http"
	"tpt_backend/models"

	"github.com/labstack/echo/v4"
)

func GetSaleForecast(c echo.Context) error {
	saleForecast, err := models.GetSaleForecast()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, saleForecast)
}
