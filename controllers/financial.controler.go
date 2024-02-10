package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"tpt_backend/models"

	"github.com/labstack/echo/v4"
)

func GetAllFinancials(c echo.Context) error {
	// Get query parameters for pagination
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.QueryParam("pageSize"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	keyword := c.QueryParam("keyword")

	typeName := "financial" // Set the type name based on your struct

	result, err := models.GetAllFinancials(typeName, page, pageSize, keyword)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	return c.JSON(http.StatusOK, result)
}

func GetFinancialDetail(c echo.Context) error {
	financialID, err := strconv.Atoi(c.Param("financial_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid financial_id"})
	}

	financialDetail, err := models.GetFinancialDetail(financialID)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	return c.JSON(http.StatusOK, financialDetail)
}

func GetFinancialBalance(c echo.Context) error {
	financialDetail, err := models.GetFinancialBalance()
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	return c.JSON(http.StatusOK, financialDetail)
}

func CreateFinancial(c echo.Context) error {
	var financial models.Financial

	// Parse the request body to populate the product struct
	if err := c.Bind(&financial); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{
				"error": "Invalid request body",
			},
		)
	}

	// Call the CreateCategory function from the models package
	result, err := models.CreateFinancial(
		financial.UserID,
		financial.Type,
		financial.FinancialDate,
		financial.Information,
		financial.CashIn,
		financial.CashOut,
		financial.Balance,
	)

	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	return c.JSON(http.StatusOK, result)
}

func UpdateFinancial(c echo.Context) error {
	// Parse the request body to get the update data
	var updateFields map[string]interface{}
	if err := json.NewDecoder(c.Request().Body).Decode(&updateFields); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{"message": "Invalid request body"},
		)
	}

	// Extract the ID from the update data
	financialID, ok := updateFields["financial_id"].(float64)
	if !ok {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{"message": "Invalid financial_id format"},
		)
	}

	// Convert id to integer
	convID := int(financialID)

	// Remove id from the updateFields map before passing it to the model
	delete(updateFields, "financial_id")

	// Call the UpdateCategory function from the models package
	result, err := models.UpdateFinancial(convID, updateFields)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	return c.JSON(http.StatusOK, result)
}

func DeleteFinancial(c echo.Context) error {
	financialID := c.Param("financial_id")

	conv_id, err := strconv.Atoi(financialID)

	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	result, err := models.DeleteFinancial(conv_id)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	return c.JSON(http.StatusOK, result)
}
