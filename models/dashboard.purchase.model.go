package models

import (
	"tpt_backend/db"
)

type PurchaseTransactions struct {
	Length int      `json:"length"`
	Title  string   `json:"title"`
	Label  []string `json:"label"`
	Value  []int    `json:"value"`
}

type TotalPurchase struct {
	Length int      `json:"length"`
	Title  string   `json:"title"`
	Label  []string `json:"label"`
	Value  []int    `json:"value"`
}

type TotalItemsPurchased struct {
	Length int      `json:"length"`
	Title  string   `json:"title"`
	Label  []string `json:"label"`
	Value  []int    `json:"value"`
}

type TopSuppliersByTransaction struct {
	Length int      `json:"length"`
	Title  string   `json:"title"`
	Label  []string `json:"label"`
	Value  []int    `json:"value"`
}

type TopSuppliersByTotalItem struct {
	Length int      `json:"length"`
	Title  string   `json:"title"`
	Label  []string `json:"label"`
	Value  []int    `json:"value"`
}

func GetPurchaseTransactions() (PurchaseTransactions, error) {
	var purchaseTransactions PurchaseTransactions
	purchaseTransactions.Length = 12
	purchaseTransactions.Title = "Purchase Transactions"
	purchaseTransactions.Label = make([]string, 12)
	purchaseTransactions.Value = make([]int, 12)

	con := db.CreateConDW()

	sqlStatement := `
		SELECT 
			CONCAT(SUBSTR(dt.month_name, 1, 3), ' ', dt.year) AS month,
			COALESCE(COUNT(fp.purchase_id), 0) AS total_transactions
		FROM 
			dim_time dt
		LEFT JOIN 
			fact_purchase fp ON fp.time_id = dt.time_id
		WHERE
			(dt.year = 2023 AND dt.month BETWEEN 4 AND 12)
			OR (dt.year = 2024 AND dt.month BETWEEN 1 AND 3)
		GROUP BY 
			dt.year,
			dt.month_name
		ORDER BY 
			dt.year,
			dt.month;
    `

	rows, err := con.Query(sqlStatement)
	if err != nil {
		return purchaseTransactions, err
	}
	defer rows.Close()

	index := 0
	for rows.Next() {
		var month string
		var totalTransactions int
		if err := rows.Scan(&month, &totalTransactions); err != nil {
			return purchaseTransactions, err
		}
		purchaseTransactions.Label[index] = month
		purchaseTransactions.Value[index] = totalTransactions
		index++
	}

	if err := rows.Err(); err != nil {
		return purchaseTransactions, err
	}

	return purchaseTransactions, nil
}

func GetTotalPurchase() (TotalPurchase, error) {
	var totalPurchase TotalPurchase
	totalPurchase.Length = 12
	totalPurchase.Title = "Total Purchase"
	totalPurchase.Label = make([]string, 12)
	totalPurchase.Value = make([]int, 12)

	con := db.CreateConDW()

	sqlStatement := `
		SELECT 
			CONCAT(SUBSTR(dt.month_name, 1, 3), ' ', dt.year) AS month,
			COALESCE(SUM(fp.total_price), 0) AS purchase_amount
		FROM 
			dim_time dt
		LEFT JOIN 
			fact_purchase fp ON fp.time_id = dt.time_id
		WHERE
			(dt.year = 2023 AND dt.month BETWEEN 4 AND 12)
			OR (dt.year = 2024 AND dt.month BETWEEN 1 AND 3)
		GROUP BY 
			dt.year,
			dt.month_name
		ORDER BY 
			dt.year,
			dt.month;
    `

	rows, err := con.Query(sqlStatement)
	if err != nil {
		return totalPurchase, err
	}
	defer rows.Close()

	index := 0
	for rows.Next() {
		var month string
		var purchaseAmount int
		if err := rows.Scan(&month, &purchaseAmount); err != nil {
			return totalPurchase, err
		}
		totalPurchase.Label[index] = month
		totalPurchase.Value[index] = purchaseAmount
		index++
	}

	if err := rows.Err(); err != nil {
		return totalPurchase, err
	}

	return totalPurchase, nil
}

func GetTotalItemsPurchased() (TotalItemsPurchased, error) {
	var totalItemsPurchasedResponse TotalItemsPurchased
	totalItemsPurchasedResponse.Length = 12
	totalItemsPurchasedResponse.Title = "Total Items Purchased"
	totalItemsPurchasedResponse.Label = make([]string, 12)
	totalItemsPurchasedResponse.Value = make([]int, 12)

	con := db.CreateConDW()

	sqlStatement := `
		SELECT 
			CONCAT(SUBSTR(dt.month_name, 1, 3), ' ', dt.year) AS month,
			COALESCE(SUM(fp.total_item), 0) AS total_items_purchased
		FROM 
			Dim_Time dt
		LEFT JOIN 
			Fact_Purchase fp ON dt.time_id = fp.time_id
		WHERE
			(dt.year = 2023 AND dt.month BETWEEN 4 AND 12)
			OR (dt.year = 2024 AND dt.month BETWEEN 1 AND 3)
		GROUP BY 
			dt.year,
			dt.month_name
		ORDER BY 
			dt.year,
			dt.month;
    `

	rows, err := con.Query(sqlStatement)
	if err != nil {
		return totalItemsPurchasedResponse, err
	}
	defer rows.Close()

	index := 0
	for rows.Next() {
		var month string
		var totalItemsPurchased int
		if err := rows.Scan(&month, &totalItemsPurchased); err != nil {
			return totalItemsPurchasedResponse, err
		}
		totalItemsPurchasedResponse.Label[index] = month
		totalItemsPurchasedResponse.Value[index] = totalItemsPurchased
		index++
	}

	if err := rows.Err(); err != nil {
		return totalItemsPurchasedResponse, err
	}

	return totalItemsPurchasedResponse, nil
}

func GetTopSuppliersByTransaction() (TopSuppliersByTransaction, error) {
	var topSuppliersByTransaction TopSuppliersByTransaction
	topSuppliersByTransaction.Length = 5
	topSuppliersByTransaction.Title = "Top Suppliers by Transaction Frequency"
	topSuppliersByTransaction.Label = make([]string, 5)
	topSuppliersByTransaction.Value = make([]int, 5)

	con := db.CreateConDW()

	sqlStatement := `
		SELECT 
			ds.supplier_name,
			COUNT(fp.purchase_id) AS transaction_count
		FROM 
			dim_supplier ds
		JOIN 
			fact_purchase fp ON ds.supplier_id = fp.supplier_id
		GROUP BY 
			ds.supplier_name
		ORDER BY 
			transaction_count DESC
		LIMIT 5;
    `

	rows, err := con.Query(sqlStatement)
	if err != nil {
		return topSuppliersByTransaction, err
	}
	defer rows.Close()

	index := 0
	for rows.Next() {
		var supplierName string
		var transactionCount int
		if err := rows.Scan(&supplierName, &transactionCount); err != nil {
			return topSuppliersByTransaction, err
		}
		topSuppliersByTransaction.Label[index] = supplierName
		topSuppliersByTransaction.Value[index] = transactionCount
		index++
	}

	if err := rows.Err(); err != nil {
		return topSuppliersByTransaction, err
	}

	return topSuppliersByTransaction, nil
}

func GetTopSuppliersByTotalItem() (TopSuppliersByTotalItem, error) {
	var topSuppliersByTotalItem TopSuppliersByTotalItem
	topSuppliersByTotalItem.Length = 5
	topSuppliersByTotalItem.Title = "Top Suppliers by Total Item Purchased"
	topSuppliersByTotalItem.Label = make([]string, 5)
	topSuppliersByTotalItem.Value = make([]int, 5)

	con := db.CreateConDW()

	sqlStatement := `
		SELECT 
			ds.supplier_name,
			COALESCE(SUM(fp.total_item), 0) AS total_items_purchased
		FROM 
			dim_supplier ds
		JOIN 
			fact_purchase fp ON ds.supplier_id = fp.supplier_id
		GROUP BY 
			ds.supplier_name
		ORDER BY 
			total_items_purchased DESC
		LIMIT 5;
    `

	rows, err := con.Query(sqlStatement)
	if err != nil {
		return topSuppliersByTotalItem, err
	}
	defer rows.Close()

	index := 0
	for rows.Next() {
		var supplierName string
		var totalItemsPurchased int
		if err := rows.Scan(&supplierName, &totalItemsPurchased); err != nil {
			return topSuppliersByTotalItem, err
		}
		topSuppliersByTotalItem.Label[index] = supplierName
		topSuppliersByTotalItem.Value[index] = totalItemsPurchased
		index++
	}

	if err := rows.Err(); err != nil {
		return topSuppliersByTotalItem, err
	}

	return topSuppliersByTotalItem, nil
}
