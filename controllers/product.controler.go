package controllers

import (
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

	typeName := "product" // Set the type name based on your struct

	result, err := models.GetAllProducts(typeName, page, pageSize)

	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{
				"message": err.Error(),
			},
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
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, productDetail)
}

// func CreateBarang(c echo.Context) error {
// 	kode_barang := c.FormValue("kode_barang")
// 	nama_barang := c.FormValue("nama_barang")

// 	fmt.Println(kode_barang)
// 	fmt.Println(nama_barang)

// 	result, err := models.CreateBarang(kode_barang, nama_barang)

// 	fmt.Println(result)

// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
// 	}

// 	return c.JSON(http.StatusOK, result)
// }

// func CreateBarangNew(c echo.Context) error {
// 	json_map := make(map[string]interface{})
// 	err := json.NewDecoder(c.Request().Body).Decode(&json_map)
// 	if err != nil {
// 		return err
// 	}

// 	kode_barang, ok := json_map["kode_barang"].(string)
// 	if !ok {
// 		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid kode_barang format"})
// 	}

// 	nama_barang, ok := json_map["nama_barang"].(string)
// 	if !ok {
// 		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid nama_barang format"})
// 	}

// 	result, err := models.CreateBarang(kode_barang, nama_barang)

// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
// 	}

// 	return c.JSON(http.StatusOK, result)
// }

// func UpdateBarang(c echo.Context) error {
// 	id := c.FormValue("id")
// 	kode_barang := c.FormValue("kode_barang")
// 	nama_barang := c.FormValue("nama_barang")

// 	conv_id, err := strconv.Atoi(id)

// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
// 	}

// 	result, err := models.UpdateBarang(conv_id, kode_barang, nama_barang)

// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
// 	}

// 	return c.JSON(http.StatusOK, result)
// }

// func UpdateBarangNew(c echo.Context) error {
// 	json_map := make(map[string]interface{})
// 	err := json.NewDecoder(c.Request().Body).Decode(&json_map)
// 	if err != nil {
// 		return err
// 	}

// 	id, ok := json_map["id"].(float64)
// 	if !ok {
// 		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid id format"})
// 	}

// 	kode_barang, ok := json_map["kode_barang"].(string)
// 	if !ok {
// 		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid kode_barang format"})
// 	}

// 	nama_barang, ok := json_map["nama_barang"].(string)
// 	if !ok {
// 		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid nama_barang format"})
// 	}

// 	// Convert id to integer
// 	conv_id := int(id)

// 	result, err := models.UpdateBarang(conv_id, kode_barang, nama_barang)
// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
// 	}

// 	return c.JSON(http.StatusOK, result)
// }

// func DeleteBarang(c echo.Context) error {
// 	id := c.Param("id")

// 	conv_id, err := strconv.Atoi(id)

// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
// 	}

// 	result, err := models.DeleteBarang(conv_id)

// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
// 	}

// 	return c.JSON(http.StatusOK, result)
// }
