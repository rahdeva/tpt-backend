package models

import (
	"fmt"
	"reflect"
	"strings"
	"time"
	"tpt_backend/db"
)

type ProductVariant struct {
	ProductVariantID   int       `json:"product_variant_id"`
	ProductID          int       `json:"product_id"`
	CategoryID         int       `json:"category_id"`
	CategoryName       string    `json:"category_name"`
	VariantName        string    `json:"variant_name"`
	ProductVariantCode string    `json:"product_variant_code"`
	ProductVariantName string    `json:"product_variant_name"`
	ProductQuantity    int       `json:"product_quantity"`
	PurchasePrice      float64   `json:"purchase_price"`
	SalePrice          float64   `json:"sale_price"`
	Sold               int       `json:"sold"`
	Image              string    `json:"image"`
	Brand              string    `json:"brand"`
	VariantStock       int       `json:"variant_stock"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// type Product struct {
//     ProductID    int       `json:"product_id"`
//     CategoryID   int       `json:"category_id"`
//     ProductName  string    `json:"product_name"`
//     Unit         string    `json:"unit"`
//     Brand        string    `json:"brand"`
//     Stock        int       `json:"stock"`
//     CreatedAt    time.Time `json:"created_at"`
//     UpdatedAt    time.Time `json:"updated_at"`
// }

type Product struct {
	ProductID     int       `json:"product_id"`
	CategoryID    int       `json:"category_id"`
	CategoryName  string    `json:"category_name"`
	CategoryColor string    `json:"category_color"`
	ProductName   string    `json:"product_name"`
	PurchasePrice int       `json:"purchase_price"`
	ProductCode   string    `json:"product_code"`
	EceranID      int       `json:"eceran_id"`
	Brand         string    `json:"brand"`
	SalePrice     int       `json:"sale_price"`
	Stock         int       `json:"stock"`
	Sold          int       `json:"sold"`
	Image         string    `json:"image"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func GetAllProductVariants(
	typeName string,
	page,
	pageSize int,
	keyword string,
	categoryID int,
) (Response, error) {
	objType := reflect.TypeOf(ResponseData(typeName))
	if objType == nil {
		return Response{}, fmt.Errorf("invalid type: %s", typeName)
	}

	var res Response
	var arrobj reflect.Value
	var meta Meta

	con := db.CreateCon()

	// Add a WHERE clause to filter products based on the keyword and category_id
	whereClause := ""
	conditions := make([]string, 0)

	if keyword != "" {
		conditions = append(conditions, fmt.Sprintf("(pv.product_variant_name LIKE '%%%s%%' OR p.brand LIKE '%%%s%%')", keyword, keyword))
	}

	if categoryID != 0 {
		conditions = append(conditions, fmt.Sprintf("p.category_id = %d", categoryID))
	}

	if len(conditions) > 0 {
		whereClause = " WHERE " + strings.Join(conditions, " AND ")
	}

	// Calculate the offset based on the page number and page size
	offset := (page - 1) * pageSize
	sqlStatement := fmt.Sprintf(`
		SELECT
			pv.product_variant_id,
			pv.product_id,
			p.category_id,
			c.category_name,
			pv.variant_name,
			pv.product_variant_code,
			pv.product_variant_name,
			pv.product_quantity,
			pv.purchase_price,
			pv.sale_price,
			pv.sold,
			pv.image,
			p.brand,
			FLOOR(p.stock / pv.product_quantity) AS variant_stock,
			pv.created_at,
			pv.updated_at
		FROM
			product_variant pv
		JOIN
			product p ON pv.product_id = p.product_id
		JOIN
			category c ON p.category_id = c.category_id
		%s
		LIMIT %d OFFSET %d;
	`, whereClause, pageSize, offset)
	fmt.Println("SQL Statement:", sqlStatement)
	rows, err := con.Query(sqlStatement)
	if err != nil {
		return res, err
	}
	defer rows.Close()

	// Count total items in the database
	var totalItems int
	err = con.QueryRow(fmt.Sprintf(`
		SELECT 
			COUNT(*) 
		FROM 
			product_variant pv
		JOIN
			product p ON pv.product_id = p.product_id
		JOIN 
			category c ON p.category_id = c.category_id 
		%s
	`, whereClause)).Scan(&totalItems)
	if err != nil {
		return res, err
	}

	// If no items are found, return an empty response data object
	if totalItems == 0 {
		meta.Limit = pageSize
		meta.Page = page
		meta.TotalPages = 0
		meta.TotalItems = totalItems

		res.Data = map[string]interface{}{
			typeName: make([]interface{}, 0), // Empty slice
			"meta":   meta,
		}

		return res, nil
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

func GetAllProducts(
	typeName string,
	page,
	pageSize int,
	keyword string,
	categoryID int,
) (Response, error) {
	objType := reflect.TypeOf(ResponseData(typeName))
	if objType == nil {
		return Response{}, fmt.Errorf("invalid type: %s", typeName)
	}

	var res Response
	var arrobj reflect.Value
	var meta Meta

	con := db.CreateCon()

	// Add a WHERE clause to filter products based on the keyword and category_id
	whereClause := ""
	conditions := make([]string, 0)

	if keyword != "" {
		conditions = append(conditions, fmt.Sprintf("(p.product_name LIKE '%%%s%%' OR p.brand LIKE '%%%s%%')", keyword, keyword))
	}

	if categoryID != 0 {
		conditions = append(conditions, fmt.Sprintf("p.category_id = %d", categoryID))
	}

	if len(conditions) > 0 {
		whereClause = " WHERE " + strings.Join(conditions, " AND ")
	}

	// Calculate the offset based on the page number and page size
	offset := (page - 1) * pageSize
	sqlStatement := fmt.Sprintf(`
		SELECT
			p.product_id,
			p.category_id,
			c.category_name,
			c.category_color,
			p.product_name,
			p.purchase_price,
			p.product_code,
			p.eceran_id,
			p.brand,
			p.sale_price,
			p.stock,
			p.sold,
			p.image,
			p.created_at,
			p.updated_at
		FROM
			product p
		JOIN
			category c ON p.category_id = c.category_id
		%s
		LIMIT %d OFFSET %d;
	`, whereClause, pageSize, offset)
	rows, err := con.Query(sqlStatement)
	if err != nil {
		return res, err
	}
	defer rows.Close()

	// Count total items in the database
	var totalItems int
	err = con.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM product p JOIN category c ON p.category_id = c.category_id %s", whereClause)).Scan(&totalItems)
	if err != nil {
		return res, err
	}

	// If no items are found, return an empty response data object
	if totalItems == 0 {
		meta.Limit = pageSize
		meta.Page = page
		meta.TotalPages = 0
		meta.TotalItems = totalItems

		res.Data = map[string]interface{}{
			typeName: make([]interface{}, 0), // Empty slice
			"meta":   meta,
		}

		return res, nil
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

func GetProductDetail(productID int) (Response, error) {
	var product Product
	var res Response

	con := db.CreateCon()

	sqlStatement := `
		SELECT
			p.product_id,
			p.category_id,
			c.category_name,
			c.category_color,
			p.product_name,
			p.purchase_price,
			p.product_code,
			p.eceran_id,
			p.brand,
			p.sale_price,
			p.stock,
			p.sold,
			p.image,
			p.created_at,
			p.updated_at
		FROM
			product p
		JOIN
			category c ON p.category_id = c.category_id
		WHERE
			p.product_id = ?;
	`

	row := con.QueryRow(sqlStatement, productID)

	err := row.Scan(
		&product.ProductID,
		&product.CategoryID,
		&product.CategoryName,
		&product.CategoryColor,
		&product.ProductName,
		&product.PurchasePrice,
		&product.ProductCode,
		&product.EceranID,
		&product.Brand,
		&product.SalePrice,
		&product.Stock,
		&product.Sold,
		&product.Image,
		&product.CreatedAt,
		&product.UpdatedAt,
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
	product.CreatedAt = product.CreatedAt.In(loc)
	product.UpdatedAt = product.UpdatedAt.In(loc)

	res.Data = map[string]interface{}{
		"product": product,
	}

	return res, nil
}

func CreateProduct(
	product_code string,
	product_name string,
	category_id int,
	eceran_id int,
	brand string,
	purchase_price int,
	sale_price int,
	stock int,
	sold int,
	image string,
) (Response, error) {
	var res Response

	con := db.CreateCon()

	sqlStatement := "INSERT INTO product (product_code, product_name, category_id, eceran_id, brand, purchase_price, sale_price, stock, sold, image, created_at, updated_at) VALUES ( ? , ? , ? , ? , ? , ? , ? , ? , ? , ? , ? , ? )"

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
		product_code,
		product_name,
		category_id,
		eceran_id,
		brand,
		purchase_price,
		sale_price,
		stock,
		sold,
		image,
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

func UpdateProduct(product_id int, updateFields map[string]interface{}) (Response, error) {
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
	sqlStatement := "UPDATE product " + setStatement + " WHERE product_id = ?"
	values = append(values, product_id)

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

func DeleteProduct(product_id int) (Response, error) {
	var res Response

	con := db.CreateCon()

	sqlStatement := "DELETE FROM product WHERE product_id = ?"

	stmt, err := con.Prepare(sqlStatement)

	if err != nil {
		return res, err
	}

	result, err := stmt.Exec(product_id)

	if err != nil {
		return res, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return res, err
	}

	res.Data = map[string]interface{}{
		"rowsAffected":       rowsAffected,
		"deleted_product_id": product_id,
	}

	return res, err
}
