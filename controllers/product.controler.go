package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"tpt_backend/models"

	"github.com/labstack/echo/v4"
)

func GetAllProducts(c echo.Context) error {
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

	typeName := "product" // Set the type name based on your struct

	result, err := models.GetAllProducts(typeName, page, pageSize, keyword)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	return c.JSON(http.StatusOK, result)
}

func GetProductDetail(c echo.Context) error {
	productID, err := strconv.Atoi(c.Param("product_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid product_id"})
	}

	productDetail, err := models.GetProductDetail(productID)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	return c.JSON(http.StatusOK, productDetail)
}

func CreateProduct(c echo.Context) error {
	var product models.Product

	// Parse the request body to populate the product struct
	if err := c.Bind(&product); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{
				"error": "Invalid request body",
			},
		)
	}

	// Call the CreateProduct function from the models package
	result, err := models.CreateProduct(
		product.ProductCode,
		product.ProductName,
		product.CategoryID,
		product.Brand,
		product.PurchasePrice,
		product.SalePrice,
		product.Stock,
		product.Sold,
		product.Image,
	)

	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	return c.JSON(http.StatusOK, result)
}

func UpdateProduct(c echo.Context) error {
	// Parse the request body to get the update data
	var updateFields map[string]interface{}
	if err := json.NewDecoder(c.Request().Body).Decode(&updateFields); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{"message": "Invalid request body"},
		)
	}

	// Extract the ID from the update data
	product_id, ok := updateFields["product_id"].(float64)
	if !ok {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{"message": "Invalid product_id format"},
		)
	}

	// Convert id to integer
	convID := int(product_id)

	// Remove id from the updateFields map before passing it to the model
	delete(updateFields, "product_id")

	// Call the UpdateProduct function from the models package
	result, err := models.UpdateProduct(convID, updateFields)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	return c.JSON(http.StatusOK, result)
}

func DeleteProduct(c echo.Context) error {
	product_id := c.Param("product_id")

	conv_id, err := strconv.Atoi(product_id)

	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	result, err := models.DeleteProduct(conv_id)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	return c.JSON(http.StatusOK, result)
}
