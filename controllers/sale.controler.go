package controllers

import (
	"net/http"
	"strconv"
	"tpt_backend/models"

	"github.com/labstack/echo/v4"
)

func GetAllSales(c echo.Context) error {
	// Get query parameters for pagination
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.QueryParam("pageSize"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	typeName := "sale" // Set the type name based on your struct

	result, err := models.GetAllSales(typeName, page, pageSize)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	return c.JSON(http.StatusOK, result)
}

func GetSalesDetail(c echo.Context) error {
	// Get query parameters for pagination
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.QueryParam("pageSize"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	saleID, err := strconv.Atoi(c.Param("sale_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid sale_id"})
	}

	typeName := "sale_detail" // Set the type name based on your struct

	result, err := models.GetSalesDetail(
		typeName,
		page,
		pageSize,
		saleID,
	)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	return c.JSON(http.StatusOK, result)
}

func GetSaleByID(c echo.Context) error {
	saleID, err := strconv.Atoi(c.Param("sale_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid sale_id"})
	}

	supplierDetail, err := models.GetSaleByID(saleID)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	return c.JSON(http.StatusOK, supplierDetail)
}

func CreateSale(c echo.Context) error {
	var createRequest models.CreateSaleRequest

	// Parse the request body to populate the new struct
	if err := c.Bind(&createRequest); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{
				"error": "Invalid request body",
			},
		)
	}

	// 	// Call the CreateSale function from the models package
	result, err := models.CreateSale(
		createRequest.UserID,
		createRequest.SaleDate,
		createRequest.TotalItem,
		createRequest.TotalPrice,
		createRequest.SalesDetail,
	)

	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	return c.JSON(http.StatusOK, result)
}

// func UpdateSale(c echo.Context) error {
// 	// Parse the request body to get the update data
// 	var updateRequest models.UpdateSaleRequest
// 	if err := json.NewDecoder(c.Request().Body).Decode(&updateRequest); err != nil {
// 		return c.JSON(
// 			http.StatusBadRequest,
// 			map[string]string{"message": "Invalid request body"},
// 		)
// 	}

// 	// Call the UpdateSale function from the models package
// 	result, err := models.UpdateSale(
// 		updateRequest.SaleID,
// 		updateRequest.UserID,
// 		updateRequest.SaleDate,
// 		updateRequest.TotalItem,
// 		updateRequest.TotalPrice,
// 		updateRequest.SalesDetail,
// 	)
// 	if err != nil {
// 		return c.JSON(
// 			http.StatusInternalServerError,
// 			map[string]string{"message": err.Error()},
// 		)
// 	}

// 	return c.JSON(http.StatusOK, result)
// }

func DeleteSale(c echo.Context) error {
	saleID := c.Param("sale_id")

	conv_id, err := strconv.Atoi(saleID)

	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	result, err := models.DeleteSale(conv_id)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	return c.JSON(http.StatusOK, result)
}
