package models

import (
	"fmt"
	"reflect"
	"time"
	"tpt_backend/db"
)

type Purchase struct {
	PurchaseID   int       `json:"purchase_id"`
	SupplierID   int       `json:"supplier_id"`
	PurchaseDate time.Time `json:"purchase_date"`
	TotalItem    int       `json:"total_item"`
	TotalPrice   int       `json:"total_price"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type PurchaseDetail struct {
	PurchaseDetailID int       `json:"purchase_detail_id"`
	PurchaseID       int       `json:"purchase_id"`
	ProductID        int       `json:"product_id"`
	PurchasePrice    int       `json:"purchase_price"`
	Quantity         int       `json:"quantity"`
	Subtotal         int       `json:"subtotal"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type CreatePurchaseRequest struct {
	SupplierID      int              `json:"supplier_id"`
	PurchaseDate    time.Time        `json:"purchase_date"`
	TotalItem       int              `json:"total_item"`
	TotalPrice      int              `json:"total_price"`
	PurchasesDetail []PurchaseDetail `json:"purchases_detail"`
}

func GetAllPurchases(typeName string, page, pageSize int) (Response, error) {
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
	err := con.QueryRow("SELECT COUNT(*) FROM purchase").Scan(&totalItems)
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
	sqlStatement := fmt.Sprintf("SELECT * FROM purchase LIMIT %d OFFSET %d", pageSize, offset)
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

func GetPurchasesDetail(
	typeName string,
	page,
	pageSize int,
	purchaseID int,
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
	err := con.QueryRow("SELECT COUNT(*) FROM purchase_detail WHERE purchase_id = ?", purchaseID).Scan(&totalItems)
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
	sqlStatement := fmt.Sprintf("SELECT * FROM purchase_detail WHERE purchase_id = ? LIMIT %d OFFSET %d", pageSize, offset)
	rows, err := con.Query(sqlStatement, purchaseID)
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
		"purchase_id": purchaseID,
		typeName:      arrobj.Interface(),
		"meta":        meta,
	}

	return res, nil
}

func GetPurchasebyID(purchaseID int) (Response, error) {
	var purchase Purchase
	var res Response

	con := db.CreateCon()

	sqlStatement := "SELECT * FROM purchase  WHERE purchase_id  = ?"

	row := con.QueryRow(sqlStatement, purchaseID)

	err := row.Scan(
		&purchase.PurchaseID,
		&purchase.SupplierID,
		&purchase.PurchaseDate,
		&purchase.TotalItem,
		&purchase.TotalPrice,
		&purchase.CreatedAt,
		&purchase.UpdatedAt,
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
	purchase.PurchaseDate = purchase.PurchaseDate.In(loc)
	purchase.CreatedAt = purchase.CreatedAt.In(loc)
	purchase.UpdatedAt = purchase.UpdatedAt.In(loc)

	res.Data = map[string]interface{}{
		"purchase": purchase,
	}

	return res, nil
}

func CreatePurchase(
	supplierID int,
	purchaseDate time.Time,
	totalItem int,
	totalPrice int,
	purchasesDetail []PurchaseDetail,
) (Response, error) {
	var res Response

	con := db.CreateCon()

	sqlStatement := "INSERT INTO purchase (supplier_id, purchase_date, total_item, total_price, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)"

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
		supplierID,
		purchaseDate,
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
	for _, detail := range purchasesDetail {
		// Assuming purchase_id is obtained from the created purchase
		sqlDetailStatement := "INSERT INTO purchase_detail (purchase_id, product_id, purchase_price, quantity, subtotal, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)"
		detailStmt, err := con.Prepare(sqlDetailStatement)
		if err != nil {
			return res, err
		}

		detailResult, err := detailStmt.Exec(
			getIdLast, // using the purchase_id obtained earlier
			detail.ProductID,
			detail.PurchasePrice,
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
		"purchase_id":   getIdLast,
		"supplier_id":   supplierID,
		"purchase_date": purchaseDate,
		"total_item":    totalItem,
		"total_price":   totalPrice,
		"created_at":    created_at,
		"updated_at":    updated_at,
	}

	return res, nil
}

// func UpdatePurchase(
// 	purchaseID int,
// 	supplierID int,
// 	purchaseDate time.Time,
// 	totalItem int,
// 	totalPrice int,
// 	purchasesDetail []PurchaseDetail,
// ) (Response, error) {
// 	var res Response

// 	con := db.CreateCon()

// 	// Load the UTC+8 time zone
// 	loc, err := time.LoadLocation("Asia/Shanghai")
// 	if err != nil {
// 		return res, err
// 	}

// 	// Construct the SET part of the SQL statement for updating the purchase
// 	setPurchaseStatement := "SET supplier_id = ?, purchase_date = ?, total_item = ?, total_price = ?, updated_at = ?"
// 	values := []interface{}{supplierID, purchaseDate.In(loc), totalItem, totalPrice, time.Now().In(loc)}

// 	// Construct the final SQL statement for updating the purchase
// 	sqlPurchaseStatement := "UPDATE purchase " + setPurchaseStatement + " WHERE purchase_id = ?"
// 	values = append(values, purchaseID)

// 	// Execute the SQL statement to update the purchase
// 	stmtPurchase, err := con.Prepare(sqlPurchaseStatement)
// 	if err != nil {
// 		return res, err
// 	}

// 	resultPurchase, err := stmtPurchase.Exec(values...)
// 	if err != nil {
// 		return res, err
// 	}

// 	// Retrieve the number of rows affected in the purchase update
// 	rowsAffectedPurchase, err := resultPurchase.RowsAffected()
// 	if err != nil {
// 		return res, err
// 	}

// 	// Iterate over purchase details and update or insert them
// 	for _, detail := range purchasesDetail {
// 		sqlDetailStatement := `
// 			INSERT INTO purchase_detail (purchase_id, product_id, purchase_price, quantity, subtotal, created_at, updated_at)
// 			VALUES (?, ?, ?, ?, ?, ?, ?)
// 			ON DUPLICATE KEY UPDATE
// 				purchase_price = VALUES(purchase_price),
// 				quantity = VALUES(quantity),
// 				subtotal = VALUES(subtotal),
// 				updated_at = VALUES(updated_at)
// 		`
// 		valuesDetail := []interface{}{purchaseID, detail.ProductID, detail.PurchasePrice, detail.Quantity, detail.Subtotal, time.Now().In(loc), time.Now().In(loc)}

// 		// Execute the SQL statement to update or insert the purchase detail
// 		stmtDetail, err := con.Prepare(sqlDetailStatement)
// 		if err != nil {
// 			return res, err
// 		}

// 		resultDetail, err := stmtDetail.Exec(valuesDetail...)
// 		if err != nil {
// 			return res, err
// 		}

// 		// Retrieve the number of rows affected in the purchase detail update or insert
// 		_ = resultDetail
// 	}

// 	res.Data = map[string]interface{}{
// 		"rowsAffectedPurchase": rowsAffectedPurchase,
// 		"updated_at":           time.Now().In(loc),
// 	}

// 	return res, nil
// }

func DeletePurchase(purchaseID int) (Response, error) {
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
	sqlDetailStatement := "DELETE FROM purchase_detail WHERE purchase_id = ?"
	detailStmt, err := tx.Prepare(sqlDetailStatement)
	if err != nil {
		return res, err
	}

	_, err = detailStmt.Exec(purchaseID)
	if err != nil {
		return res, err
	}

	// Delete from purchase
	sqlStatement := "DELETE FROM purchase WHERE purchase_id = ?"
	stmt, err := tx.Prepare(sqlStatement)
	if err != nil {
		return res, err
	}

	result, err := stmt.Exec(purchaseID)
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
		"rowsAffected":        rowsAffected,
		"deleted_purchase_id": purchaseID,
	}

	return res, nil
}
