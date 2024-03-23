package controllers

import (
	"net/http"
	"tpt_backend/models"

	"github.com/labstack/echo/v4"
)

func GetCashFlow(c echo.Context) error {
	cashFlow, err := models.GetCashFlow()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, cashFlow)
}

func GetFinancialByType(c echo.Context) error {
	financialTransaction, err := models.GetFinancialByType()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, financialTransaction)
}

func GetCashIn(c echo.Context) error {
	cashIn, err := models.GetCashIn()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, cashIn)
}

func GetCashOut(c echo.Context) error {
	cashOut, err := models.GetCashOut()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, cashOut)
}
