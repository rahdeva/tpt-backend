package controllers

import (
	"net/http"
	"tpt_backend/models"

	"github.com/labstack/echo/v4"
)

func GetProductByCategory(c echo.Context) error {
	productByCategory, err := models.GetProductByCategory()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, productByCategory)
}

func GetTopProductStock(c echo.Context) error {
	topProductStock, err := models.GetTopProductStock()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, topProductStock)
}

func GetTopExpensiveProducts(c echo.Context) error {
	topExpensiveProduct, err := models.GetTopExpensiveProducts()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, topExpensiveProduct)
}
