package controllers

import (
	"net/http"
	"tpt_backend/models"

	"github.com/labstack/echo/v4"
)

func GetPurchaseTransactions(c echo.Context) error {
	purchaseTransactions, err := models.GetPurchaseTransactions()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, purchaseTransactions)
}

func GetTotalPurchase(c echo.Context) error {
	totalPurchase, err := models.GetTotalPurchase()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, totalPurchase)
}

func GetTotalItemsPurchased(c echo.Context) error {
	totalItemPurchased, err := models.GetTotalItemsPurchased()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, totalItemPurchased)
}

func GetTopSuppliersByTransaction(c echo.Context) error {
	topSupplier, err := models.GetTopSuppliersByTransaction()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, topSupplier)
}

func GetTopSuppliersByTotalItem(c echo.Context) error {
	topSupplier, err := models.GetTopSuppliersByTotalItem()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, topSupplier)
}
