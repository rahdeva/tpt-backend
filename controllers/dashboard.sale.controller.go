package controllers

import (
	"net/http"
	"tpt_backend/models"

	"github.com/labstack/echo/v4"
)

func GetSalesTransactions(c echo.Context) error {
	totalSalesTransactions, err := models.GetSalesTransactions()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, totalSalesTransactions)
}

func GetTotalSales(c echo.Context) error {
	totalSales, err := models.GetTotalSales()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, totalSales)
}

func GetTotalItemsSold(c echo.Context) error {
	totalItemsSold, err := models.GetTotalItemsSold()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, totalItemsSold)
}

func GetTotalProfit(c echo.Context) error {
	totalProfit, err := models.GetTotalProfit()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, totalProfit)
}

func GetTopBestSellingProducts(c echo.Context) error {
	topSaleProduct, err := models.GetTopBestSellingProducts()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, topSaleProduct)
}
