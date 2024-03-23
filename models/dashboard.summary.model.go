package models

import (
	"tpt_backend/db"
)

// ProductByCategory represents the product quantity by category
type ProductByCategory struct {
	Length int      `json:"length"`
	Title  string   `json:"title"`
	Label  []string `json:"label"`
	Value  []int    `json:"value"`
}

type TopProductStock struct {
	Length int      `json:"length"`
	Title  string   `json:"title"`
	Label  []string `json:"label"`
	Value  []int    `json:"value"`
}

type TopExpensiveProducts struct {
	Length int      `json:"length"`
	Title  string   `json:"title"`
	Label  []string `json:"label"`
	Value  []int    `json:"value"`
}

func GetProductByCategory() (ProductByCategory, error) {
	var productByCategory ProductByCategory
	productByCategory.Length = 4
	productByCategory.Title = "Product by Category"
	productByCategory.Label = make([]string, 4)
	productByCategory.Value = make([]int, 4)

	con := db.CreateConDW()

	sqlStatement := `
		SELECT 
			category_name,
			SUM(product_quantity) AS total_barang
		FROM 
			dim_product_variant
		GROUP BY 
			category_name
		ORDER BY 
			total_barang DESC
    `

	rows, err := con.Query(sqlStatement)
	if err != nil {
		return productByCategory, err
	}
	defer rows.Close()

	index := 0
	for rows.Next() {
		var category_name string
		var total_barang int
		if err := rows.Scan(&category_name, &total_barang); err != nil {
			return productByCategory, err
		}
		productByCategory.Label[index] = category_name
		productByCategory.Value[index] = total_barang
		index++
	}

	if err := rows.Err(); err != nil {
		return productByCategory, err
	}

	return productByCategory, nil
}

func GetTopProductStock() (TopProductStock, error) {
	var topProducts TopProductStock
	topProducts.Length = 5
	topProducts.Title = "Top 5 Products by Ending Stock"
	topProducts.Label = make([]string, 5)
	topProducts.Value = make([]int, 5)

	con := db.CreateConDW()

	sqlStatement := `
		SELECT 
			product_variant_name,
			variant_stock
		FROM 
			dim_product_variant
		ORDER BY 
			variant_stock DESC
		LIMIT 5
    `

	rows, err := con.Query(sqlStatement)
	if err != nil {
		return topProducts, err
	}
	defer rows.Close()

	index := 0
	for rows.Next() {
		var productVariantName string
		var variantStock int
		if err := rows.Scan(&productVariantName, &variantStock); err != nil {
			return topProducts, err
		}
		topProducts.Label[index] = productVariantName
		topProducts.Value[index] = variantStock
		index++
	}

	if err := rows.Err(); err != nil {
		return topProducts, err
	}

	return topProducts, nil
}

func GetTopExpensiveProducts() (TopExpensiveProducts, error) {
	var topExpensiveProducts TopExpensiveProducts
	topExpensiveProducts.Length = 5 // Ubah ke 1 jika hanya ingin satu produk teratas
	topExpensiveProducts.Title = "Top Products by Sale Price"
	topExpensiveProducts.Label = make([]string, 5) // Ubah ke 1 jika hanya ingin satu produk teratas
	topExpensiveProducts.Value = make([]int, 5)    // Ubah ke 1 jika hanya ingin satu produk teratas

	con := db.CreateConDW()

	sqlStatement := `
		SELECT 
			product_variant_name,
			sale_price
		FROM 
			dim_product_variant
		ORDER BY 
			sale_price DESC
		LIMIT 5
    `

	rows, err := con.Query(sqlStatement)
	if err != nil {
		return topExpensiveProducts, err
	}
	defer rows.Close()

	index := 0
	for rows.Next() {
		var productVariantName string
		var salePrice int
		if err := rows.Scan(&productVariantName, &salePrice); err != nil {
			return topExpensiveProducts, err
		}
		topExpensiveProducts.Label[index] = productVariantName
		topExpensiveProducts.Value[index] = salePrice
		index++
	}

	if err := rows.Err(); err != nil {
		return topExpensiveProducts, err
	}

	return topExpensiveProducts, nil
}
