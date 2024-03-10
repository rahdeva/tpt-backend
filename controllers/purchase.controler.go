package controllers

import (
	"net/http"
	"strconv"
	"tpt_backend/models"

	"github.com/labstack/echo/v4"
)

func GetAllPurchases(c echo.Context) error {
	// Get query parameters for pagination
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.QueryParam("pageSize"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	typeName := "purchase" // Set the type name based on your struct

	result, err := models.GetAllPurchases(typeName, page, pageSize)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	return c.JSON(http.StatusOK, result)
}

func GetPurchasesDetail(c echo.Context) error {
	// Get query parameters for pagination
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.QueryParam("pageSize"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	purchaseID, err := strconv.Atoi(c.Param("purchase_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid purchase_id"})
	}

	typeName := "purchase_detail" // Set the type name based on your struct

	result, err := models.GetPurchasesDetail(
		typeName,
		page,
		pageSize,
		purchaseID,
	)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	return c.JSON(http.StatusOK, result)
}

func GetPurchasebyID(c echo.Context) error {
	purchaseID, err := strconv.Atoi(c.Param("purchase_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid purchase_id"})
	}

	supplierDetail, err := models.GetPurchasebyID(purchaseID)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	return c.JSON(http.StatusOK, supplierDetail)
}

// func CreatePurchase(c echo.Context) error {
// 	var createRequest models.CreatePurchaseRequest

// 	// Parse the request body to populate the new struct
// 	if err := c.Bind(&createRequest); err != nil {
// 		return c.JSON(
// 			http.StatusBadRequest,
// 			map[string]string{
// 				"error": "Invalid request body",
// 			},
// 		)
// 	}

// 	// Call the CreatePurchase function from the models package
// 	result, err := models.CreatePurchase(
// 		createRequest.UserID,
// 		createRequest.SupplierID,
// 		createRequest.PurchaseDate,
// 		createRequest.TotalItem,
// 		createRequest.TotalPrice,
// 		createRequest.PurchasesDetail,
// 	)

// 	if err != nil {
// 		return c.JSON(
// 			http.StatusInternalServerError,
// 			map[string]string{"message": err.Error()},
// 		)
// 	}

// 	return c.JSON(http.StatusOK, result)
// }

// func UpdatePurchase(c echo.Context) error {
// 	// Parse the request body to get the update data
// 	var updateRequest models.UpdatePurchaseRequest
// 	if err := json.NewDecoder(c.Request().Body).Decode(&updateRequest); err != nil {
// 		return c.JSON(
// 			http.StatusBadRequest,
// 			map[string]string{"message": "Invalid request body"},
// 		)
// 	}

// 	// Call the UpdatePurchase function from the models package
// 	result, err := models.UpdatePurchase(
// 		updateRequest.PurchaseID,
// 		updateRequest.SupplierID,
// 		updateRequest.PurchaseDate,
// 		updateRequest.TotalItem,
// 		updateRequest.TotalPrice,
// 		updateRequest.PurchasesDetail,
// 	)
// 	if err != nil {
// 		return c.JSON(
// 			http.StatusInternalServerError,
// 			map[string]string{"message": err.Error()},
// 		)
// 	}

// 	return c.JSON(http.StatusOK, result)
// }

// func DeletePurchase(c echo.Context) error {
// 	purchaseID := c.Param("purchase_id")

// 	conv_id, err := strconv.Atoi(purchaseID)

// 	if err != nil {
// 		return c.JSON(
// 			http.StatusInternalServerError,
// 			map[string]string{"message": err.Error()},
// 		)
// 	}

// 	result, err := models.DeletePurchase(conv_id)
// 	if err != nil {
// 		return c.JSON(
// 			http.StatusInternalServerError,
// 			map[string]string{"message": err.Error()},
// 		)
// 	}

// 	return c.JSON(http.StatusOK, result)
// }
