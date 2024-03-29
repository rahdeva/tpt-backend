package controllers

import (
	"net/http"
	"strconv"
	"tpt_backend/models"

	"github.com/labstack/echo/v4"
)

func GetAllProductVariants(c echo.Context) error {
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

	categoryID, err := strconv.Atoi(c.QueryParam("category_id"))
	if err != nil {
		categoryID = 0
	}

	typeName := "product_variant" // Set the type name based on your struct

	result, err := models.GetAllProductVariants(typeName, page, pageSize, keyword, categoryID)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	return c.JSON(http.StatusOK, result)
}

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

	sort := c.QueryParam("sort") // Get sorting parameter

	categoryID, err := strconv.Atoi(c.QueryParam("category_id"))
	if err != nil {
		categoryID = 0
	}

	typeName := "product" // Set the type name based on your struct

	result, err := models.GetAllProducts(typeName, page, pageSize, keyword, categoryID, sort)
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
	var createRequest models.CreateProductRequest

	// Parse the request body to populate the new struct
	if err := c.Bind(&createRequest); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{
				"error": "Invalid request body",
			},
		)
	}

	// Call CreateProduct function from models package
	result, err := models.CreateProduct(
		createRequest.ProductName,
		createRequest.CategoryID,
		createRequest.Unit,
		createRequest.Stock,
		createRequest.Brand,
		createRequest.Variants,
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

// func UpdateProduct(c echo.Context) error {
// 	// Parse the request body to get the update data
// 	var updateFields map[string]interface{}
// 	if err := json.NewDecoder(c.Request().Body).Decode(&updateFields); err != nil {
// 		return c.JSON(
// 			http.StatusBadRequest,
// 			map[string]string{"message": "Invalid request body"},
// 		)
// 	}

// 	// Extract the ID from the update data
// 	product_id, ok := updateFields["product_id"].(float64)
// 	if !ok {
// 		return c.JSON(
// 			http.StatusBadRequest,
// 			map[string]string{"message": "Invalid product_id format"},
// 		)
// 	}

// 	// Convert id to integer
// 	convID := int(product_id)

// 	// Remove id from the updateFields map before passing it to the model
// 	delete(updateFields, "product_id")

// 	// Call the UpdateProduct function from the models package
// 	result, err := models.UpdateProduct(convID, updateFields)
// 	if err != nil {
// 		return c.JSON(
// 			http.StatusInternalServerError,
// 			map[string]string{"message": err.Error()},
// 		)
// 	}

// 	return c.JSON(http.StatusOK, result)
// }

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
