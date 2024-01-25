package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"tpt_backend/models"

	"github.com/labstack/echo/v4"
)

func GetAllUsers(c echo.Context) error {
	// Get query parameters for pagination
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.QueryParam("pageSize"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	typeName := "user" // Set the type name based on your struct

	result, err := models.GetAllUsers(typeName, page, pageSize)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	return c.JSON(http.StatusOK, result)
}

func GetUserDetail(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user_id"})
	}

	userDetail, err := models.GetUserDetail(userID)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	return c.JSON(http.StatusOK, userDetail)
}

func GetUserDetailbyUID(c echo.Context) error {
	uid := c.Param("uid")

	userDetailbyUID, err := models.GetUserDetailByUID(uid)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	return c.JSON(http.StatusOK, userDetailbyUID)
}

func CreateUser(c echo.Context) error {
	var user models.User

	// Parse the request body to populate the product struct
	if err := c.Bind(&user); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{
				"error": "Invalid request body",
			},
		)
	}

	// Call the CreateCategory function from the models package
	result, err := models.CreateUser(
		user.RoleID,
		user.UID,
		user.Name,
		user.Email,
		user.Address,
		user.PhoneNumber,
		user.ProfilePicture,
	)

	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	return c.JSON(http.StatusOK, result)
}

func UpdateUser(c echo.Context) error {
	// Parse the request body to get the update data
	var updateFields map[string]interface{}
	if err := json.NewDecoder(c.Request().Body).Decode(&updateFields); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{"message": "Invalid request body"},
		)
	}

	// Extract the ID from the update data
	userID, ok := updateFields["user_id"].(float64)
	if !ok {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{"message": "Invalid user_id format"},
		)
	}

	// Convert id to integer
	convID := int(userID)

	// Remove id from the updateFields map before passing it to the model
	delete(updateFields, "user_id")

	// Call the UpdateCategory function from the models package
	result, err := models.UpdateUser(convID, updateFields)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	return c.JSON(http.StatusOK, result)
}

func DeleteUser(c echo.Context) error {
	userID := c.Param("user_id")

	conv_id, err := strconv.Atoi(userID)

	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	result, err := models.DeleteUser(conv_id)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"message": err.Error()},
		)
	}

	return c.JSON(http.StatusOK, result)
}
