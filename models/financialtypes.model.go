package models

import (
	"fmt"
	"reflect"
	"time"
	"tpt_backend/db"
)

type FinancialType struct {
	FinancialTypeID   int       `json:"financial_type_id"`
	FinancialTypeName string    `json:"financial_type_name"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func GetAllFinancialTypes(typeName string, page, pageSize int) (Response, error) {
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
	err := con.QueryRow("SELECT COUNT(*) FROM financial_type").Scan(&totalItems)
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
	sqlStatement := fmt.Sprintf("SELECT * FROM financial_type LIMIT %d OFFSET %d", pageSize, offset)
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

func GetFinancialTypeDetail(financialTypeID int) (Response, error) {
	var financialType FinancialType
	var res Response

	con := db.CreateCon()

	sqlStatement := "SELECT * FROM financial_type WHERE financial_type_id = ?"

	row := con.QueryRow(sqlStatement, financialTypeID)

	err := row.Scan(
		&financialType.FinancialTypeID,
		&financialType.FinancialTypeName,
		&financialType.CreatedAt,
		&financialType.UpdatedAt,
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
	financialType.CreatedAt = financialType.CreatedAt.In(loc)
	financialType.UpdatedAt = financialType.UpdatedAt.In(loc)

	res.Data = map[string]interface{}{
		"financialType": financialType,
	}

	return res, nil
}

func CreateFinancialType(financialTypeName string) (Response, error) {
	var res Response

	con := db.CreateCon()

	sqlStatement := "INSERT INTO financial_type (financial_type_name, created_at, updated_at) VALUES (?, ?, ?)"

	stmt, err := con.Prepare(sqlStatement)

	if err != nil {
		return res, err
	}

	// Load the UTC+8 time zone
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return res, err
	}

	created_at := time.Now()
	updated_at := time.Now()

	result, err := stmt.Exec(
		financialTypeName,
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
		"created_at": created_at.In(loc),
	}

	return res, nil
}

func UpdateFinancialType(financialTypeID int, updateFields map[string]interface{}) (Response, error) {
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
	sqlStatement := "UPDATE financial_type " + setStatement + " WHERE financial_type_id  = ?"
	values = append(values, financialTypeID)

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

func DeleteFinancialType(financialTypeID int) (Response, error) {
	var res Response

	con := db.CreateCon()

	sqlStatement := "DELETE FROM financial_type WHERE financial_type_id = ?"

	stmt, err := con.Prepare(sqlStatement)

	if err != nil {
		return res, err
	}

	result, err := stmt.Exec(financialTypeID)

	if err != nil {
		return res, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return res, err
	}

	res.Data = map[string]interface{}{
		"rowsAffected":              rowsAffected,
		"deleted_financial_type_id": financialTypeID,
	}

	return res, err
}
