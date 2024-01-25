package models

import (
	"fmt"
	"reflect"
	"time"
	"tpt_backend/db"
)

type User struct {
	UserID         int       `json:"user_id"`
	RoleID         int       `json:"role_id"`
	UID            string    `json:"uid"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	Address        string    `json:"address"`
	PhoneNumber    string    `json:"phone_number"`
	ProfilePicture string    `json:"profile_picture"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func GetAllUsers(typeName string, page, pageSize int) (Response, error) {
	objType := reflect.TypeOf(ResponseData(typeName))
	if objType == nil {
		return Response{}, fmt.Errorf("invalid type: %s", typeName)
	}

	var res Response
	var arrobj reflect.Value
	var meta Meta

	con := db.CreateCon()

	// Count total items in the database
	var totalItems int
	err := con.QueryRow("SELECT COUNT(*) FROM user").Scan(&totalItems)
	if err != nil {
		return res, err
	}

	// Load the UTC+8 time zone
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return res, err
	}

	// Calculate the total number of pages
	totalPages := calculateTotalPages(totalItems, pageSize)

	// Check if the requested page is greater than the total number of pages
	if page > totalPages {
		return res, fmt.Errorf("requested page (%d) exceeds total number of pages (%d)", page, totalPages)
	}

	// Calculate the offset based on the page number and page size
	offset := (page - 1) * pageSize
	sqlStatement := fmt.Sprintf("SELECT * FROM user LIMIT %d OFFSET %d", pageSize, offset)
	rows, err := con.Query(sqlStatement)
	if err != nil {
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		obj := ResponseData(typeName)
		objValue := reflect.ValueOf(obj).Elem() // Dereference the pointer

		if objValue.Kind() != reflect.Struct {
			return res, fmt.Errorf("expected a struct type, got %v", objValue.Kind())
		}

		fields := make([]interface{}, objValue.NumField())
		for i := 0; i < objValue.NumField(); i++ {
			fields[i] = objValue.Field(i).Addr().Interface()
		}

		err := rows.Scan(fields...)
		if err != nil {
			return res, err
		}

		// Convert time fields to UTC+8 (Asia/Shanghai) before including them in the response
		createdAtField, _ := objValue.Type().FieldByName("CreatedAt")
		updatedAtField, _ := objValue.Type().FieldByName("UpdatedAt")

		if createdAtField.Type == reflect.TypeOf(time.Time{}) {
			createdAtFieldIndex := createdAtField.Index[0]
			createdAtValue := objValue.Field(createdAtFieldIndex).Interface().(time.Time)
			objValue.Field(createdAtFieldIndex).Set(reflect.ValueOf(createdAtValue.In(loc)))
		}

		if updatedAtField.Type == reflect.TypeOf(time.Time{}) {
			updatedAtFieldIndex := updatedAtField.Index[0]
			updatedAtValue := objValue.Field(updatedAtFieldIndex).Interface().(time.Time)
			objValue.Field(updatedAtFieldIndex).Set(reflect.ValueOf(updatedAtValue.In(loc)))
		}

		if !arrobj.IsValid() {
			arrobj = reflect.MakeSlice(reflect.SliceOf(objType.Elem()), 0, 0)
		}

		arrobj = reflect.Append(arrobj, objValue)

		meta.Limit = pageSize
		meta.Page = page
		meta.TotalPages = calculateTotalPages(totalItems, pageSize)
		meta.TotalItems = totalItems
	}

	res.Data = map[string]interface{}{
		typeName: arrobj.Interface(),
		"meta":   meta,
	}

	return res, nil
}

func GetUserDetail(userID int) (Response, error) {
	var user User
	var res Response

	con := db.CreateCon()

	sqlStatement := "SELECT * FROM user WHERE user_id = ?"

	row := con.QueryRow(sqlStatement, userID)

	err := row.Scan(
		&user.UserID,
		&user.RoleID,
		&user.UID,
		&user.Name,
		&user.Email,
		&user.Address,
		&user.PhoneNumber,
		&user.ProfilePicture,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return res, err
	}

	// Load the UTC+8 time zone
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return res, err
	}

	// Convert time fields to UTC+8 (Asia/Shanghai) before including them in the response
	user.CreatedAt = user.CreatedAt.In(loc)
	user.UpdatedAt = user.UpdatedAt.In(loc)

	res.Data = map[string]interface{}{
		"user": user,
	}

	return res, nil
}

func GetUserDetailByUID(uid string) (Response, error) {
	var user User
	var res Response

	con := db.CreateCon()

	sqlStatement := "SELECT * FROM user WHERE uid = ?"

	row := con.QueryRow(sqlStatement, uid)

	err := row.Scan(
		&user.UserID,
		&user.RoleID,
		&user.UID,
		&user.Name,
		&user.Email,
		&user.Address,
		&user.PhoneNumber,
		&user.ProfilePicture,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return res, err
	}

	// Load the UTC+8 time zone
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return res, err
	}

	// Convert time fields to UTC+8 (Asia/Shanghai) before including them in the response
	user.CreatedAt = user.CreatedAt.In(loc)
	user.UpdatedAt = user.UpdatedAt.In(loc)

	res.Data = map[string]interface{}{
		"user": user,
	}

	return res, nil
}

func CreateUser(
	roleID int,
	uid string,
	name string,
	email string,
	address string,
	phoneNumber string,
	profilePicture string,
) (Response, error) {
	var res Response

	con := db.CreateCon()

	sqlStatement := "INSERT INTO user (role_id, uid, name, email, address, phone_number, profile_picture, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)"

	stmt, err := con.Prepare(sqlStatement)

	if err != nil {
		return res, err
	}

	// Load the UTC+8 time zone
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return res, err
	}

	created_at := time.Now().In(loc)
	updated_at := time.Now().In(loc)

	result, err := stmt.Exec(
		roleID,
		uid,
		name,
		email,
		address,
		phoneNumber,
		profilePicture,
		created_at,
		updated_at,
	)

	if err != nil {
		return res, err
	}

	getIdLast, err := result.LastInsertId()

	if err != nil {
		return res, err
	}

	res.Data = map[string]interface{}{
		"getIdLast":  getIdLast,
		"created_at": created_at,
	}

	return res, nil
}

func UpdateUser(userID int, updateFields map[string]interface{}) (Response, error) {
	var res Response

	// Load the UTC+8 time zone
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return res, err
	}

	// Add or update the 'updated_at' field in the updateFields map
	updateFields["updated_at"] = time.Now().In(loc)
	updated_at := updateFields["updated_at"]

	con := db.CreateCon()

	// Construct the SET part of the SQL statement dynamically
	setStatement := "SET "
	values := []interface{}{}
	i := 0

	for fieldName, fieldValue := range updateFields {
		if i > 0 {
			setStatement += ", "
		}
		setStatement += fieldName + " = ?"
		values = append(values, fieldValue)
		i++
	}

	// Construct the final SQL statement
	sqlStatement := "UPDATE user " + setStatement + " WHERE user_id = ?"
	values = append(values, userID)

	stmt, err := con.Prepare(sqlStatement)
	if err != nil {
		return res, err
	}

	result, err := stmt.Exec(values...)
	if err != nil {
		return res, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return res, err
	}

	res.Data = map[string]interface{}{
		"rowsAffected": rowsAffected,
		"updated_at":   updated_at,
	}

	return res, nil
}

func DeleteUser(userID int) (Response, error) {
	var res Response

	con := db.CreateCon()

	sqlStatement := "DELETE FROM user WHERE user_id = ?"

	stmt, err := con.Prepare(sqlStatement)

	if err != nil {
		return res, err
	}

	result, err := stmt.Exec(userID)

	if err != nil {
		return res, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return res, err
	}

	res.Data = map[string]interface{}{
		"rowsAffected":    rowsAffected,
		"deleted_user_id": userID,
	}

	return res, err
}
