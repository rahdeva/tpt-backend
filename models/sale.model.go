package models

import (
	"fmt"
	"reflect"
	"time"
	"tpt_backend/db"
)

type Sale struct {
	SaleID     int       `json:"sale_id"`
	UserID     int       `json:"user_id"`
	UserName   string    `json:"user_name"`
	SaleDate   time.Time `json:"sale_date"`
	TotalItem  int       `json:"total_item"`
	TotalPrice int       `json:"total_price"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type SaleDetail struct {
	SaleDetailID int       `json:"sale_detail_id"`
	SaleID       int       `json:"sale_id"`
	ProductID    int       `json:"product_id"`
	ProductCode  string    `json:"product_code"`
	ProductName  string    `json:"product_name"`
	SalePrice    int       `json:"sale_price"`
	Quantity     int       `json:"quantity"`
	Subtotal     int       `json:"subtotal"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CreateSaleRequest struct {
	UserID      int          `json:"user_id"`
	SaleDate    time.Time    `json:"sale_date"`
	TotalItem   int          `json:"total_item"`
	TotalPrice  int          `json:"total_price"`
	SalesDetail []SaleDetail `json:"sales_detail"`
}

type UpdateSaleRequest struct {
	SaleID      int          `json:"sale_id"`
	UserID      int          `json:"user_id"`
	SaleDate    time.Time    `json:"sale_date"`
	TotalItem   int          `json:"total_item"`
	TotalPrice  int          `json:"total_price"`
	SalesDetail []SaleDetail `json:"sales_detail"`
}

func GetAllSales(typeName string, page, pageSize int) (Response, error) {
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
	err := con.QueryRow("SELECT COUNT(*) FROM sale").Scan(&totalItems)
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
	sqlStatement := fmt.Sprintf(`
		SELECT
			s.sale_id,
			s.user_id,
			u.name AS user_name,
			s.sale_date,
			s.total_item,
			s.total_price,
			s.created_at,
			s.updated_at
		FROM
			sale s
		JOIN
			user u ON s.user_id = u.user_id
		ORDER BY s.sale_id DESC
		LIMIT %d OFFSET %d;
	`, pageSize, offset)
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
		saleDateField, _ := objValue.Type().FieldByName("SaleDate") // assuming the field name is "SaleDate"

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

		if saleDateField.Type == reflect.TypeOf(time.Time{}) {
			saleDateFieldIndex := saleDateField.Index[0]
			saleDateValue := objValue.Field(saleDateFieldIndex).Interface().(time.Time)
			objValue.Field(saleDateFieldIndex).Set(reflect.ValueOf(saleDateValue.In(loc)))
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

func GetSalesDetail(
	typeName string,
	page,
	pageSize int,
	saleID int,
) (Response, error) {
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
	err := con.QueryRow("SELECT COUNT(*) FROM sale_detail WHERE sale_id = ?", saleID).Scan(&totalItems)
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
	sqlStatement := fmt.Sprintf(`
		SELECT
			sd.sale_detail_id,
			sd.sale_id,
			sd.product_id,
			p.product_code,
			p.product_name,
			sd.sale_price,
			sd.quantity,
			sd.subtotal,
			sd.created_at,
			sd.updated_at
		FROM
			sale_detail sd
		JOIN
			product p ON sd.product_id = p.product_id
		WHERE
			sd.sale_id = ?
		LIMIT %d OFFSET %d;
	`, pageSize, offset)
	rows, err := con.Query(sqlStatement, saleID)
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
		"sale_id": saleID,
		typeName:  arrobj.Interface(),
		"meta":    meta,
	}

	return res, nil
}

func GetSaleByID(saleID int) (Response, error) {
	var sale Sale
	var res Response

	con := db.CreateCon()

	sqlStatement := `
		SELECT
			s.sale_id,
			s.user_id,
			u.name AS user_name,
			s.sale_date,
			s.total_item,
			s.total_price,
			s.created_at,
			s.updated_at
		FROM
			sale s
		JOIN
			user u ON s.user_id = u.user_id
		WHERE
			s.sale_id = ?;
	`

	row := con.QueryRow(sqlStatement, saleID)

	err := row.Scan(
		&sale.SaleID,
		&sale.UserID,
		&sale.UserName,
		&sale.SaleDate,
		&sale.TotalItem,
		&sale.TotalPrice,
		&sale.CreatedAt,
		&sale.UpdatedAt,
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
	sale.SaleDate = sale.SaleDate.In(loc)
	sale.CreatedAt = sale.CreatedAt.In(loc)
	sale.UpdatedAt = sale.UpdatedAt.In(loc)

	res.Data = map[string]interface{}{
		"sale": sale,
	}

	return res, nil
}

func CreateSale(
	userId int,
	saleDate time.Time,
	totalItem int,
	totalPrice int,
	salesDetail []SaleDetail,
) (Response, error) {
	var res Response

	con := db.CreateCon()

	sqlStatement := "INSERT INTO sale (user_id, sale_date, total_item, total_price, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)"

	stmt, err := con.Prepare(sqlStatement)

	if err != nil {
		return res, err
	}

	created_at := time.Now()
	updated_at := time.Now()

	result, err := stmt.Exec(
		userId,
		saleDate,
		totalItem,
		totalPrice,
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

	// Insert purchase details
	for _, detail := range salesDetail {
		// Assuming purchase_id is obtained from the created purchase
		sqlDetailStatement := "INSERT INTO sale_detail (sale_id, product_id, sale_price, quantity, subtotal, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)"
		detailStmt, err := con.Prepare(sqlDetailStatement)
		if err != nil {
			return res, err
		}

		detailResult, err := detailStmt.Exec(
			getIdLast, // using the sale_id obtained earlier
			detail.ProductID,
			detail.SalePrice,
			detail.Quantity,
			detail.Subtotal,
			created_at,
			updated_at,
		)

		if err != nil {
			return res, err
		}

		// Use the detail result or handle as needed
		_ = detailResult
	}

	res.Data = map[string]interface{}{
		"sale_id":     getIdLast,
		"user_id":     userId,
		"sale_date":   saleDate,
		"total_item":  totalItem,
		"total_price": totalPrice,
		"created_at":  created_at,
		"updated_at":  updated_at,
	}

	return res, nil
}

func UpdateSale(
	saleID int,
	userID int,
	saleDate time.Time,
	totalItem int,
	totalPrice int,
	salesDetail []SaleDetail,
) (Response, error) {
	var res Response

	con := db.CreateCon()

	// Load the UTC+8 time zone
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return res, err
	}

	// Construct the SET part of the SQL statement for updating the sale
	setSaleStatement := "SET user_id = ?, sale_date = ?, total_item = ?, total_price = ?, updated_at = ?"
	values := []interface{}{userID, saleDate, totalItem, totalPrice, time.Now()}

	// Construct the final SQL statement for updating the sale
	sqlSaleStatement := "UPDATE sale " + setSaleStatement + " WHERE sale_id = ?"
	values = append(values, saleID)

	// Execute the SQL statement to update the sale
	stmtSale, err := con.Prepare(sqlSaleStatement)
	if err != nil {
		return res, err
	}

	resultSale, err := stmtSale.Exec(values...)
	if err != nil {
		return res, err
	}

	// Retrieve the number of rows affected in the sale update
	rowsAffectedSale, err := resultSale.RowsAffected()
	if err != nil {
		return res, err
	}

	// Get existing sale details
	existingDetails, err := getExistingSaleDetails(saleID)
	if err != nil {
		return res, err
	}

	// Iterate over existing sale details and mark those to be deleted
	detailsToDelete := make(map[int]bool)
	for _, existingDetail := range existingDetails {
		detailsToDelete[existingDetail.SaleDetailID] = true
	}

	// Iterate over sale details and update or insert them
	for _, detail := range salesDetail {
		if detail.SaleDetailID > 0 {
			// Mark existing detail as not to be deleted
			delete(detailsToDelete, detail.SaleDetailID)

			// Update existing sale detail
			sqlDetailStatement := `
				UPDATE sale_detail
				SET sale_price = ?, quantity = ?, subtotal = ?, updated_at = ?
				WHERE sale_detail_id = ?
			`
			valuesDetail := []interface{}{detail.SalePrice, detail.Quantity, detail.Subtotal, time.Now(), detail.SaleDetailID}
			_, err := con.Exec(sqlDetailStatement, valuesDetail...)
			if err != nil {
				return res, err
			}
		} else {
			// Insert new sale detail
			sqlDetailStatement := `
				INSERT INTO sale_detail (sale_id, product_id, sale_price, quantity, subtotal, created_at, updated_at)
				VALUES (?, ?, ?, ?, ?, ?, ?)
			`
			valuesDetail := []interface{}{saleID, detail.ProductID, detail.SalePrice, detail.Quantity, detail.Subtotal, time.Now(), time.Now()}
			_, err := con.Exec(sqlDetailStatement, valuesDetail...)
			if err != nil {
				return res, err
			}
		}
	}

	// Delete details that are not present in the request
	for detailID := range detailsToDelete {
		sqlDeleteDetail := "DELETE FROM sale_detail WHERE sale_detail_id = ?"
		_, err := con.Exec(sqlDeleteDetail, detailID)
		if err != nil {
			return res, err
		}
	}

	res.Data = map[string]interface{}{
		"rowsAffectedSale": rowsAffectedSale,
		"updated_at":       time.Now().In(loc),
	}

	return res, nil
}

func getExistingSaleDetails(saleID int) ([]SaleDetail, error) {
	con := db.CreateCon()

	// Fetch existing sale details for the given sale ID
	sqlStatement := "SELECT sale_detail_id, product_id, sale_price, quantity, subtotal FROM sale_detail WHERE sale_id = ?"
	rows, err := con.Query(sqlStatement, saleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate over rows and populate the existing details
	var existingDetails []SaleDetail
	for rows.Next() {
		var detail SaleDetail
		if err := rows.Scan(&detail.SaleDetailID, &detail.ProductID, &detail.SalePrice, &detail.Quantity, &detail.Subtotal); err != nil {
			return nil, err
		}
		existingDetails = append(existingDetails, detail)
	}

	return existingDetails, nil
}

func DeleteSale(saleID int) (Response, error) {
	var res Response

	con := db.CreateCon()

	// Begin a transaction
	tx, err := con.Begin()
	if err != nil {
		return res, err
	}

	// Defer a function to handle rollback in case of an error
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// Delete from purchase_detail first
	sqlDetailStatement := "DELETE FROM sale_detail WHERE sale_id = ?"
	detailStmt, err := tx.Prepare(sqlDetailStatement)
	if err != nil {
		return res, err
	}

	_, err = detailStmt.Exec(saleID)
	if err != nil {
		return res, err
	}

	// Delete from purchase
	sqlStatement := "DELETE FROM sale WHERE sale_id = ?"
	stmt, err := tx.Prepare(sqlStatement)
	if err != nil {
		return res, err
	}

	result, err := stmt.Exec(saleID)
	if err != nil {
		return res, err
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return res, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return res, err
	}

	res.Data = map[string]interface{}{
		"rowsAffected":    rowsAffected,
		"deleted_sale_id": saleID,
	}

	return res, nil
}
