package models

import (
	"fmt"
	"reflect"
	"time"
	"tpt_backend/db"
)

type Product struct {
	ProductID     int       `json:"product_id"`
	ProductCode   string    `json:"product_code"`
	ProductName   string    `json:"product_name"`
	CategoryID    int       `json:"category_id"`
	Brand         string    `json:"brand"`
	PurchasePrice int       `json:"purchase_price"`
	SalePrice     int       `json:"sale_price"`
	Stock         int       `json:"stock"`
	Sold          int       `json:"sold"`
	Image         string    `json:"image"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func GetAllProducts(typeName string, page, pageSize int) (Response, error) {
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
	err := con.QueryRow("SELECT COUNT(*) FROM product").Scan(&totalItems)
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
	sqlStatement := fmt.Sprintf("SELECT * FROM product LIMIT %d OFFSET %d", pageSize, offset)
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

// func CreateBarang(kode_barang string, nama_barang string) (Response, error) {
// 	var res Response

// 	con := db.CreateCon()

// 	sqlStatement := "INSERT INTO barang (kode_barang, nama_barang) VALUES(? , ? )"

// 	stmt, err := con.Prepare(sqlStatement)

// 	if err != nil {
// 		return res, err
// 	}

// 	result, err := stmt.Exec(kode_barang, nama_barang)

// 	if err != nil {
// 		return res, err
// 	}

// 	getIdLast, err := result.LastInsertId()

// 	if err != nil {
// 		return res, err
// 	}

// 	res.Status = http.StatusOK
// 	res.Message = "Sukses"
// 	res.Data = map[string]int64{
// 		"getIdLast": getIdLast,
// 	}

// 	return res, nil
// }

// func UpdateBarang(id int, kode_barang string, nama_barang string) (Response, error) {
// 	var res Response

// 	con := db.CreateCon()

// 	sqlStatement := "UPDATE barang SET nama_barang = ? , kode_barang = ? WHERE id = ?"

// 	stmt, err := con.Prepare(sqlStatement)

// 	if err != nil {
// 		return res, err
// 	}

// 	result, err := stmt.Exec(nama_barang, kode_barang, id)

// 	if err != nil {
// 		return res, err
// 	}

// 	rowsAffected, err := result.RowsAffected()

// 	if err != nil {
// 		return res, err
// 	}

// 	res.Status = http.StatusOK
// 	res.Message = "Sukses"
// 	res.Data = map[string]int64{
// 		"rows": rowsAffected,
// 	}

// 	return res, nil
// }

// func DeleteBarang(id int) (Response, error) {
// 	var res Response

// 	con := db.CreateCon()

// 	sqlStatement := "DELETE FROM barang WHERE id = ? "

// 	stmt, err := con.Prepare(sqlStatement)

// 	if err != nil {
// 		return res, err
// 	}

// 	result, err := stmt.Exec(id)

// 	if err != nil {
// 		return res, err
// 	}

// 	rowsAffected, err := result.RowsAffected()

// 	if err != nil {
// 		return res, err
// 	}

// 	res.Status = http.StatusOK
// 	res.Message = "Sukses"
// 	res.Data = map[string]int64{
// 		"rows": rowsAffected,
// 	}

// 	return res, err
// }
