package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"tpt_backend/models"

	"github.com/labstack/echo/v4"
)

func GetAllRoles(c echo.Context) error {
	// Get query parameters for pagination
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.QueryParam("pageSize"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	typeName := "role" // Set the type name based on your struct

	result, err := models.GetAllRoles(typeName, page, pageSize)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	return c.JSON(http.StatusOK, result)
}

func GetRoleDetail(c echo.Context) error {
	roleID, err := strconv.Atoi(c.Param("role_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid role_id"})
	}

	roleDetail, err := models.GetRoleDetail(roleID)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	return c.JSON(http.StatusOK, roleDetail)
}

func CreateRole(c echo.Context) error {
	var role models.Role

	// Parse the request body to populate the product struct
	if err := c.Bind(&role); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{
				"error": "Invalid request body",
			},
		)
	}

	// Call the CreateRole function from the models package
	result, err := models.CreateRole(
		role.RoleName,
	)

	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	return c.JSON(http.StatusOK, result)
}

func UpdateRole(c echo.Context) error {
	// Parse the request body to get the update data
	var updateFields map[string]interface{}
	if err := json.NewDecoder(c.Request().Body).Decode(&updateFields); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{"message": "Invalid request body"},
		)
	}

	// Extract the ID from the update data
	roleID, ok := updateFields["role_id"].(float64)
	if !ok {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{"message": "Invalid role_id format"},
		)
	}

	// Convert id to integer
	convID := int(roleID)

	// Remove id from the updateFields map before passing it to the model
	delete(updateFields, "role_id")

	// Call the UpdateCategory function from the models package
	result, err := models.UpdateRole(convID, updateFields)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	return c.JSON(http.StatusOK, result)
}

func DeleteRole(c echo.Context) error {
	roleID := c.Param("role_id")

	conv_id, err := strconv.Atoi(roleID)

	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	result, err := models.DeleteRole(conv_id)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	return c.JSON(http.StatusOK, result)
}
