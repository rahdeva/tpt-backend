package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"tpt_backend/models"

	"github.com/labstack/echo/v4"
)

func GetAllFinancialTypes(c echo.Context) error {
	// Get query parameters for pagination
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.QueryParam("pageSize"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	typeName := "financialtype" // Set the type name based on your struct

	result, err := models.GetAllFinancialTypes(typeName, page, pageSize)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	return c.JSON(http.StatusOK, result)
}

func GetFinancialTypeDetail(c echo.Context) error {
	financialTypeID, err := strconv.Atoi(c.Param("financial_type_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid financial_type_id"})
	}

	financialTypeDetail, err := models.GetFinancialTypeDetail(financialTypeID)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	return c.JSON(http.StatusOK, financialTypeDetail)
}

func CreateFinancialType(c echo.Context) error {
	var financialType models.FinancialType

	// Parse the request body to populate the product struct
	if err := c.Bind(&financialType); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{
				"error": "Invalid request body",
			},
		)
	}

	// Call the CreatefinancialType function from the models package
	result, err := models.CreateFinancialType(
		financialType.FinancialTypeName,
	)

	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	return c.JSON(http.StatusOK, result)
}

func UpdateFinancialType(c echo.Context) error {
	// Parse the request body to get the update data
	var updateFields map[string]interface{}
	if err := json.NewDecoder(c.Request().Body).Decode(&updateFields); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{"message": "Invalid request body"},
		)
	}

	// Extract the ID from the update data
	financialTypeID, ok := updateFields["financial_type_id"].(float64)
	if !ok {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{"message": "Invalid financial_type_id format"},
		)
	}

	// Convert id to integer
	convID := int(financialTypeID)

	// Remove id from the updateFields map before passing it to the model
	delete(updateFields, "financial_type_id")

	// Call the UpdateCategory function from the models package
	result, err := models.UpdateFinancialType(convID, updateFields)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	return c.JSON(http.StatusOK, result)
}

func DeleteFinancialType(c echo.Context) error {
	financialTypeID := c.Param("financial_type_id")

	conv_id, err := strconv.Atoi(financialTypeID)

	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	result, err := models.DeleteFinancialType(conv_id)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	return c.JSON(http.StatusOK, result)
}
