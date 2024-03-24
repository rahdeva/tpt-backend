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
	Image              string    `json:"image"`
	Brand              string    `json:"brand"`
	VariantStock       int       `json:"variant_stock"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type Product struct {
	ProductID    int       `json:"product_id"`
	CategoryID   int       `json:"category_id"`
	CategoryName string    `json:"category_name"`
	ProductName  string    `json:"product_name"`
	Unit         string    `json:"unit"`
	Brand        string    `json:"brand"`
	Stock        int       `json:"stock"`
	TotalVariant int       `json:"total_variant"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ProductDetail struct {
	ProductID      int                    `json:"product_id"`
	CategoryID     int                    `json:"category_id"`
	CategoryName   string                 `json:"category_name"`
	ProductName    string                 `json:"product_name"`
	Unit           string                 `json:"unit"`
	Brand          string                 `json:"brand"`
	Stock          int                    `json:"stock"`
	TotalVariant   int                    `json:"total_variant"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
	ProductVariant []ProductDetailVariant `json:"product_variant"`
}

type ProductDetailVariant struct {
	ProductVariantID   int     `json:"product_variant_id"`
	VariantName        string  `json:"variant_name"`
	ProductVariantCode string  `json:"product_variant_code"`
	ProductVariantName string  `json:"product_variant_name"`
	ProductQuantity    int     `json:"product_quantity"`
	PurchasePrice      float64 `json:"purchase_price"`
	SalePrice          float64 `json:"sale_price"`
	Image              string  `json:"image"`
}

type ProductVariantAdd struct {
	VariantName        string `json:"variant_name"`
	ProductVariantCode string `json:"product_variant_code"`
	ProductVariantName string `json:"product_variant_name"`
	ProductQuantity    int    `json:"product_quantity"`
	PurchasePrice      int    `json:"purchase_price"`
	SalePrice          int    `json:"sale_price"`
	Image              string `json:"image"`
}

type CreateProductRequest struct {
	ProductName string              `json:"product_name"`
	CategoryID  int                 `json:"category_id"`
	Unit        string              `json:"unit"`
	Stock       int                 `json:"stock"`
	Brand       string              `json:"brand"`
	Variants    []ProductVariantAdd `json:"variants"`
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
	sort string,
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

	var orderBy string
	if sort != "" {
		switch sort {
		case "asc":
			orderBy = "p.stock ASC"
		case "desc":
			orderBy = "p.stock DESC"
		default:
			orderBy = "p.product_id ASC"
		}
	}
	if sort == "" {
		orderBy = "p.product_id ASC"
	}

	// Calculate the offset based on the page number and page size
	offset := (page - 1) * pageSize
	sqlStatement := fmt.Sprintf(`
		SELECT
			p.product_id,
			p.category_id,
			c.category_name,
			p.product_name,
			p.unit,
			p.brand,
			p.stock,
			COUNT(pv.product_variant_id) AS total_variant,
			p.created_at,
			p.updated_at
		FROM
			product p
		JOIN
			category c ON p.category_id = c.category_id
		LEFT JOIN
			product_variant pv ON p.product_id = pv.product_id
		%s
		GROUP BY
			p.product_id, c.category_name, p.product_name, p.unit, p.brand, p.created_at, p.updated_at
		ORDER BY
            %s
        LIMIT %d OFFSET %d;
    `, whereClause, orderBy, pageSize, offset)
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
	var productDetail ProductDetail
	var res Response

	con := db.CreateCon()

	// Query untuk mendapatkan detail produk berdasarkan ID
	productStmt := `
		SELECT
			p.product_id,
			p.category_id,
			c.category_name,
			p.product_name,
			p.unit,
			p.brand,
			p.stock,
			COUNT(pv.product_variant_id) AS total_variant,
			p.created_at,
			p.updated_at
		FROM
			product p
		JOIN
			category c ON p.category_id = c.category_id
		LEFT JOIN
			product_variant pv ON p.product_id = pv.product_id
		WHERE
			p.product_id = ?
		GROUP BY
			p.product_id, c.category_name, p.product_name, p.unit, p.brand, p.created_at, p.updated_at;
	`

	// Eksekusi query produk
	row := con.QueryRow(productStmt, productID)

	// Memindai hasil query ke variabel productDetail
	err := row.Scan(
		&productDetail.ProductID,
		&productDetail.CategoryID,
		&productDetail.CategoryName,
		&productDetail.ProductName,
		&productDetail.Unit,
		&productDetail.Brand,
		&productDetail.Stock,
		&productDetail.TotalVariant,
		&productDetail.CreatedAt,
		&productDetail.UpdatedAt,
	)

	// Handle error jika terjadi
	if err != nil {
		return res, err
	}

	// Load the UTC+8 time zone
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return res, err
	}

	// Konversi waktu ke UTC+8 (Asia/Shanghai)
	productDetail.CreatedAt = productDetail.CreatedAt.In(loc)
	productDetail.UpdatedAt = productDetail.UpdatedAt.In(loc)

	// Query untuk mendapatkan detail variant produk
	variantStmt := `
		SELECT
			product_variant_id,
			variant_name,
			product_variant_code,
			product_variant_name,
			product_quantity,
			purchase_price,
			sale_price,
			image
		FROM
			product_variant
		WHERE
			product_id = ?
		ORDER BY
			product_quantity;
	`

	// Eksekusi query variant produk
	rows, err := con.Query(variantStmt, productID)
	if err != nil {
		return res, err
	}
	defer rows.Close()

	// Variabel untuk menyimpan detail variant produk
	var variants []ProductDetailVariant

	// Memindai hasil query variant produk dan menyimpannya dalam variabel variants
	for rows.Next() {
		var variant ProductDetailVariant
		err := rows.Scan(
			&variant.ProductVariantID,
			&variant.VariantName,
			&variant.ProductVariantCode,
			&variant.ProductVariantName,
			&variant.ProductQuantity,
			&variant.PurchasePrice,
			&variant.SalePrice,
			&variant.Image,
		)
		if err != nil {
			return res, err
		}
		variants = append(variants, variant)
	}

	// Mengisi informasi variant produk ke dalam ProductDetail
	productDetail.ProductVariant = variants

	// Mengatur respons dengan data dan tanpa error
	res.Data = map[string]interface{}{
		"product": productDetail,
	}
	res.Error = ""

	return res, nil
}

func CreateProduct(
	productName string,
	categoryID int,
	unit string,
	stock int,
	brand string,
	variants []ProductVariantAdd,
) (Response, error) {
	var res Response

	con := db.CreateCon()

	sqlStatement := "INSERT INTO product (product_name, category_id, unit, stock, brand, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)"

	stmt, err := con.Prepare(sqlStatement)

	if err != nil {
		return res, err
	}

	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return res, err
	}

	createdAt := time.Now()
	updatedAt := time.Now()

	result, err := stmt.Exec(
		productName,
		categoryID,
		unit,
		stock,
		brand,
		createdAt,
		updatedAt,
	)
	if err != nil {
		return res, err
	}

	getIdLast, err := result.LastInsertId()

	if err != nil {
		return res, err
	}

	for _, variant := range variants {
		// Insert variant data into product_variant table
		sqlStatement := `
			INSERT INTO product_variant 
				(product_id, variant_name, product_variant_code, product_variant_name, product_quantity, purchase_price, sale_price, image, created_at, updated_at) 
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`

		stmt, err := con.Prepare(sqlStatement)
		if err != nil {
			return res, err
		}

		result, err := stmt.Exec(
			getIdLast,
			variant.VariantName,
			variant.ProductVariantCode,
			variant.ProductVariantName,
			variant.ProductQuantity,
			variant.PurchasePrice,
			variant.SalePrice,
			variant.Image,
			createdAt,
			updatedAt,
		)
		if err != nil {
			return res, err
		}

		_ = result
	}

	res.Data = map[string]interface{}{
		"product_id": getIdLast,
		"created_at": createdAt.In(loc),
	}

	return res, nil
}

// func UpdateProduct(product_id int, updateFields map[string]interface{}) (Response, error) {
// 	var res Response

// 	// Load the UTC+8 time zone
// 	loc, err := time.LoadLocation("Asia/Shanghai")
// 	if err != nil {
// 		return res, err
// 	}

// 	// Add or update the 'updated_at' field in the updateFields map
// 	updateFields["updated_at"] = time.Now().In(loc)
// 	updated_at := updateFields["updated_at"]

// 	con := db.CreateCon()

// 	// Construct the SET part of the SQL statement dynamically
// 	setStatement := "SET "
// 	values := []interface{}{}
// 	i := 0

// 	for fieldName, fieldValue := range updateFields {
// 		if i > 0 {
// 			setStatement += ", "
// 		}
// 		setStatement += fieldName + " = ?"
// 		values = append(values, fieldValue)
// 		i++
// 	}

// 	// Construct the final SQL statement
// 	sqlStatement := "UPDATE product " + setStatement + " WHERE product_id = ?"
// 	values = append(values, product_id)

// 	stmt, err := con.Prepare(sqlStatement)
// 	if err != nil {
// 		return res, err
// 	}

// 	result, err := stmt.Exec(values...)
// 	if err != nil {
// 		return res, err
// 	}

// 	rowsAffected, err := result.RowsAffected()
// 	if err != nil {
// 		return res, err
// 	}

// 	res.Data = map[string]interface{}{
// 		"rowsAffected": rowsAffected,
// 		"updated_at":   updated_at,
// 	}

// 	return res, nil
// }

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
