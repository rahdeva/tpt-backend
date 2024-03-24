package models

import (
	"context"
	"fmt"
	"reflect"
	"time"
	"tpt_backend/db"
)

type Purchase struct {
	PurchaseID   int       `json:"purchase_id"`
	UserID       int       `json:"user_id"`
	SupplierID   int       `json:"supplier_id"`
	UserName     string    `json:"user_name"`
	SupplierName string    `json:"supplier_name"`
	PurchaseDate time.Time `json:"purchase_date"`
	TotalItem    int       `json:"total_item"`
	TotalPrice   int       `json:"total_price"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type PurchaseDetail struct {
	PurchaseDetailID   int       `json:"purchase_detail_id"`
	PurchaseID         int       `json:"purchase_id"`
	ProductVariantID   int       `json:"product_variant_id"`
	VariantName        string    `json:"variant_name"`
	ProductVariantCode string    `json:"product_variant_code"`
	ProductVariantName string    `json:"product_variant_name"`
	ProductQuantity    int       `json:"product_quantity"`
	PurchasePrice      int       `json:"purchase_price"`
	Quantity           int       `json:"quantity"`
	Subtotal           int       `json:"subtotal"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type CreatePurchaseRequest struct {
	UserID          int              `json:"user_id"`
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
	sqlStatement := fmt.Sprintf(`
		SELECT
			p.purchase_id,
			p.supplier_id,
			p.user_id,
			u.name AS user_name,
			s.supplier_name,
			p.purchase_date,
			p.total_item,
			p.total_price,
			p.created_at,
			p.updated_at
		FROM
			purchase p
		JOIN
			supplier s ON p.supplier_id = s.supplier_id
		JOIN
			user u ON p.user_id = u.user_id
		ORDER BY p.purchase_id DESC
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
		purchaseDateField, _ := objValue.Type().FieldByName("PurchaseDate") // assuming the field name is "PurchaseDate"

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

		if purchaseDateField.Type == reflect.TypeOf(time.Time{}) {
			purchaseDateFieldIndex := purchaseDateField.Index[0]
			purchaseDateValue := objValue.Field(purchaseDateFieldIndex).Interface().(time.Time)
			objValue.Field(purchaseDateFieldIndex).Set(reflect.ValueOf(purchaseDateValue.In(loc)))
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
	sqlStatement := fmt.Sprintf(`
		SELECT
			pd.purchase_detail_id,
			pd.purchase_id,
			pd.product_variant_id,
			pv.variant_name,
			pv.product_variant_code,
			pv.product_variant_name,
			pv.product_quantity,
			pd.purchase_price,
			pd.quantity,
			pd.subtotal,
			pd.created_at,
			pd.updated_at
		FROM
			purchase_detail pd
		JOIN
			product_variant pv ON pd.product_variant_id = pv.product_variant_id
		WHERE
			pd.purchase_id = ?
		LIMIT %d OFFSET %d;
	`, pageSize, offset)
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

	sqlStatement := `
		SELECT
			p.purchase_id,
			p.supplier_id,
			p.user_id,
			u.name AS user_name,
			s.supplier_name,
			p.purchase_date,
			p.total_item,
			p.total_price,
			p.created_at,
			p.updated_at
		FROM
			purchase p
		JOIN
			supplier s ON p.supplier_id = s.supplier_id
		JOIN
			user u ON p.user_id = u.user_id
		WHERE
			p.purchase_id = ?;
	`

	row := con.QueryRow(sqlStatement, purchaseID)

	// Scan
	err := row.Scan(
		&purchase.PurchaseID,
		&purchase.UserID,
		&purchase.SupplierID,
		&purchase.UserName,
		&purchase.SupplierName,
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
	userId int,
	supplierID int,
	purchaseDate time.Time,
	totalItem int,
	totalPrice int,
	purchasesDetail []PurchaseDetail,
) (Response, error) {
	var res Response

	con := db.CreateCon()

	sqlStatement := "INSERT INTO purchase (user_id, supplier_id, purchase_date, total_item, total_price, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)"

	stmt, err := con.Prepare(sqlStatement)

	if err != nil {
		return res, err
	}

	created_at := time.Now()
	updated_at := time.Now()

	result, err := stmt.Exec(
		userId,
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
		sqlDetailStatement := "INSERT INTO purchase_detail (purchase_id, product_variant_id, purchase_price, quantity, subtotal, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)"
		detailStmt, err := con.Prepare(sqlDetailStatement)
		if err != nil {
			return res, err
		}

		detailResult, err := detailStmt.Exec(
			getIdLast,
			detail.ProductVariantID,
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

	// Launch goroutine to insert data into dim_purchase_detail and fact_purchase asynchronously
	go func(ctx context.Context) {
		// Insert into dim_purchase_detail
		err = InsertIntoDimPurchaseDetail(ctx, purchasesDetail, getIdLast)
		if err != nil {
			fmt.Println("Error inserting into dim_purchase_detail:", err)
			return
		}

		// Insert into fact_purchase
		err = InsertIntoFactPurchase(ctx, getIdLast, userId, purchaseDate, totalItem, totalPrice)
		if err != nil {
			fmt.Println("Error inserting into fact_purchase:", err)
			return
		}
	}(context.Background())

	// Update product stock
	for _, detail := range purchasesDetail {
		sqlUpdateStock := `
			UPDATE product p 
			JOIN product_variant pv 
			ON p.product_id = pv.product_id 
			SET p.stock = p.stock + ? *pv.product_quantity 
			WHERE pv.product_variant_id = ?
		`
		updateStmt, err := con.Prepare(sqlUpdateStock)
		if err != nil {
			return res, err
		}

		_, err = updateStmt.Exec(detail.Quantity, detail.ProductVariantID)
		if err != nil {
			return res, err
		}
	}

	res.Data = map[string]interface{}{
		"purchase_id":   getIdLast,
		"user_id":       userId,
		"supplier_id":   supplierID,
		"purchase_date": purchaseDate,
		"total_item":    totalItem,
		"total_price":   totalPrice,
		"created_at":    created_at,
		"updated_at":    updated_at,
	}

	return res, nil
}

func UpdatePurchase(
	purchaseID int,
	supplierID int,
	purchaseDate time.Time,
	totalItem int,
	totalPrice int,
	purchasesDetail []PurchaseDetail,
) (Response, error) {
	var res Response

	con := db.CreateCon()

	// Load the UTC+8 time zone
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return res, err
	}

	// Construct the SET part of the SQL statement for updating the purchase
	setPurchaseStatement := "SET supplier_id = ?, purchase_date = ?, total_item = ?, total_price = ?, updated_at = ?"
	values := []interface{}{supplierID, purchaseDate, totalItem, totalPrice, time.Now()}

	// Construct the final SQL statement for updating the purchase
	sqlPurchaseStatement := "UPDATE purchase " + setPurchaseStatement + " WHERE purchase_id = ?"
	values = append(values, purchaseID)

	// Execute the SQL statement to update the purchase
	stmtPurchase, err := con.Prepare(sqlPurchaseStatement)
	if err != nil {
		return res, err
	}

	resultPurchase, err := stmtPurchase.Exec(values...)
	if err != nil {
		return res, err
	}

	// Retrieve the number of rows affected in the purchase update
	rowsAffectedPurchase, err := resultPurchase.RowsAffected()
	if err != nil {
		return res, err
	}

	// Get existing purchase details
	existingDetails, err := getExistingPurchaseDetails(purchaseID)
	if err != nil {
		return res, err
	}

	// Iterate over existing purchase details and mark those to be deleted
	detailsToDelete := make(map[int]bool)
	for _, existingDetail := range existingDetails {
		detailsToDelete[existingDetail.PurchaseDetailID] = true
	}

	// Iterate over purchase details and update or insert them
	for _, detail := range purchasesDetail {
		if detail.PurchaseDetailID > 0 {
			// Mark existing detail as not to be deleted
			delete(detailsToDelete, detail.PurchaseDetailID)

			// Update existing purchase detail
			sqlDetailStatement := `
				UPDATE purchase_detail
				SET purchase_price = ?, quantity = ?, subtotal = ?, updated_at = ?
				WHERE purchase_detail_id = ?
			`
			valuesDetail := []interface{}{detail.PurchasePrice, detail.Quantity, detail.Subtotal, time.Now(), detail.PurchaseDetailID}
			_, err := con.Exec(sqlDetailStatement, valuesDetail...)
			if err != nil {
				return res, err
			}
		} else {
			// Insert new purchase detail
			sqlDetailStatement := `
				INSERT INTO purchase_detail (purchase_id, product_id, purchase_price, quantity, subtotal, created_at, updated_at)
				VALUES (?, ?, ?, ?, ?, ?, ?)
			`
			valuesDetail := []interface{}{purchaseID, detail.ProductVariantID, detail.PurchasePrice, detail.Quantity, detail.Subtotal, time.Now(), time.Now()}
			_, err := con.Exec(sqlDetailStatement, valuesDetail...)
			if err != nil {
				return res, err
			}
		}
	}

	// Delete details that are not present in the request
	for detailID := range detailsToDelete {
		sqlDeleteDetail := "DELETE FROM purchase_detail WHERE purchase_detail_id = ?"
		_, err := con.Exec(sqlDeleteDetail, detailID)
		if err != nil {
			return res, err
		}
	}

	res.Data = map[string]interface{}{
		"rowsAffectedPurchase": rowsAffectedPurchase,
		"updated_at":           time.Now().In(loc),
	}

	return res, nil
}

func getExistingPurchaseDetails(purchaseID int) ([]PurchaseDetail, error) {
	con := db.CreateCon()

	// Fetch existing purchase details for the given purchase ID
	sqlStatement := "SELECT purchase_detail_id, product_id, purchase_price, quantity, subtotal FROM purchase_detail WHERE purchase_id = ?"
	rows, err := con.Query(sqlStatement, purchaseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate over rows and populate the existing details
	var existingDetails []PurchaseDetail
	for rows.Next() {
		var detail PurchaseDetail
		if err := rows.Scan(&detail.PurchaseDetailID, &detail.ProductVariantID, &detail.PurchasePrice, &detail.Quantity, &detail.Subtotal); err != nil {
			return nil, err
		}
		existingDetails = append(existingDetails, detail)
	}

	return existingDetails, nil
}

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

// Fungsi baru untuk menyimpan data ke dim_purchase_detail
func InsertIntoDimPurchaseDetail(ctx context.Context, purchasesDetail []PurchaseDetail, getIdLast int64) error {
	// Connect to data warehouse
	conDW := db.CreateConDW()

	for _, detail := range purchasesDetail {
		// Check if context is cancelled
		select {
		case <-ctx.Done():
			fmt.Println("Context cancelled. Exiting function InsertIntoDimPurchaseDetail.")
			return ctx.Err() // Return context error
		default:
			// Continue processing
		}

		// Replace the placeholders below with appropriate column names and values
		sqlDimDetailStatement := "INSERT INTO dim_purchase_detail (purchase_id, product_variant_id, purchase_price, quantity, subtotal) VALUES (?, ?, ?, ?, ?)"
		dimDetailStmt, err := conDW.Prepare(sqlDimDetailStatement)
		if err != nil {
			fmt.Println("Error preparing statement for dim_purchase_detail:", err)
			return err
		}
		_, err = dimDetailStmt.ExecContext(ctx,
			getIdLast,
			detail.ProductVariantID,
			detail.PurchasePrice,
			detail.Quantity,
			detail.Subtotal,
		)
		if err != nil {
			fmt.Println("Error inserting into dim_purchase_detail:", err)
			return err
		}

		// Simulating the insertion process
		fmt.Println("Inserted into dim_purchase_detail:", getIdLast, detail.ProductVariantID, detail.PurchasePrice, detail.Quantity, detail.Subtotal)
	}
	return nil
}

// Fungsi baru untuk menyimpan data ke fact_purchase
func InsertIntoFactPurchase(ctx context.Context, purchaseID int64, userID int, purchaseDate time.Time, totalItem int, totalPrice int) error {
	// Connect to data warehouse
	conDW := db.CreateConDW()

	timeID := purchaseDate.Format("20060102")

	_, err := conDW.ExecContext(ctx, "INSERT INTO fact_purchase (purchase_id, user_id, time_id, total_item, total_price) VALUES (?, ?, ?, ?, ?)",
		purchaseID,
		userID,
		timeID,
		totalItem,
		totalPrice,
	)
	if err != nil {
		fmt.Println("Error inserting into fact_purchase:", err)
		return err
	}
	fmt.Println("Inserted into fact_purchase:", purchaseID, userID, timeID, totalItem, totalPrice)

	return nil
}
