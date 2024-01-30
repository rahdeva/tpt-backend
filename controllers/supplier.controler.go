package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"tpt_backend/models"

	"github.com/labstack/echo/v4"
)

func GetAllSuppliers(c echo.Context) error {
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

	typeName := "supplier" // Set the type name based on your struct

	result, err := models.GetAllSuppliers(typeName, page, pageSize, keyword)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	return c.JSON(http.StatusOK, result)
}

func GetSupplierDetail(c echo.Context) error {
	supplierID, err := strconv.Atoi(c.Param("supplier_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid supplier_id"})
	}

	supplierDetail, err := models.GetSupplierDetail(supplierID)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	return c.JSON(http.StatusOK, supplierDetail)
}

func CreateSupplier(c echo.Context) error {
	var supplier models.Supplier

	// Parse the request body to populate the product struct
	if err := c.Bind(&supplier); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{
				"error": "Invalid request body",
			},
		)
	}

	// Call the CreateCategory function from the models package
	result, err := models.CreateSupplier(
		supplier.SupplierName,
		supplier.PhoneNumber,
		supplier.Address,
	)

	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	return c.JSON(http.StatusOK, result)
}

func UpdateSupplier(c echo.Context) error {
	// Parse the request body to get the update data
	var updateFields map[string]interface{}
	if err := json.NewDecoder(c.Request().Body).Decode(&updateFields); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{"message": "Invalid request body"},
		)
	}

	// Extract the ID from the update data
	supplierID, ok := updateFields["supplier_id"].(float64)
	if !ok {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{"message": "Invalid supplier_id format"},
		)
	}

	// Convert id to integer
	convID := int(supplierID)

	// Remove id from the updateFields map before passing it to the model
	delete(updateFields, "supplier_id")

	// Call the UpdateCategory function from the models package
	result, err := models.UpdateSupplier(convID, updateFields)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	return c.JSON(http.StatusOK, result)
}

func DeleteSupplier(c echo.Context) error {
	supplierID := c.Param("supplier_id")

	conv_id, err := strconv.Atoi(supplierID)

	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	result, err := models.DeleteSupplier(conv_id)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	return c.JSON(http.StatusOK, result)
}
