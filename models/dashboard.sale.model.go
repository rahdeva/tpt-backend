package models

import "tpt_backend/db"

type SalesTransactions struct {
	Length int      `json:"length"`
	Title  string   `json:"title"`
	Label  []string `json:"label"`
	Value  []int    `json:"value"`
}

type TotalSalesResponse struct {
	Length int      `json:"length"`
	Title  string   `json:"title"`
	Label  []string `json:"label"`
	Value  []int    `json:"value"`
}

type TotalItemsSold struct {
	Length int      `json:"length"`
	Title  string   `json:"title"`
	Label  []string `json:"label"`
	Value  []int    `json:"value"`
}

type TotalProfit struct {
	Length int      `json:"length"`
	Title  string   `json:"title"`
	Label  []string `json:"label"`
	Value  []int    `json:"value"`
}

type TopBestSellingProducts struct {
	Length int      `json:"length"`
	Title  string   `json:"title"`
	Label  []string `json:"label"`
	Value  []int    `json:"value"`
}

func GetSalesTransactions() (SalesTransactions, error) {
	var salesTransactions SalesTransactions
	salesTransactions.Length = 12
	salesTransactions.Title = "Sales Transactions"
	salesTransactions.Label = make([]string, 12)
	salesTransactions.Value = make([]int, 12)

	con := db.CreateConDW()

	sqlStatement := `
		SELECT 
			CONCAT(SUBSTR(dt.month_name, 1, 3), ' ', dt.year) AS month,
			COALESCE(COUNT(fs.sale_id), 0) AS total_transactions
		FROM 
			dim_time dt
		LEFT JOIN 
			fact_sale fs ON fs.time_id = dt.time_id
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
		return salesTransactions, err
	}
	defer rows.Close()

	index := 0
	for rows.Next() {
		var month string
		var totalTransactions int
		if err := rows.Scan(&month, &totalTransactions); err != nil {
			return salesTransactions, err
		}
		salesTransactions.Label[index] = month
		salesTransactions.Value[index] = totalTransactions
		index++
	}

	if err := rows.Err(); err != nil {
		return salesTransactions, err
	}

	return salesTransactions, nil
}

func GetTotalSales() (TotalSalesResponse, error) {
	var totalSalesResponse TotalSalesResponse
	totalSalesResponse.Length = 12
	totalSalesResponse.Title = "Total Sales"
	totalSalesResponse.Label = make([]string, 12)
	totalSalesResponse.Value = make([]int, 12)

	con := db.CreateConDW()

	sqlStatement := `
        SELECT 
            CONCAT(SUBSTR(dt.month_name, 1, 3), ' ', dt.year) AS month,
            COALESCE(SUM(fs.total_price), 0) AS sales_amount
        FROM 
            Dim_Time dt
        LEFT JOIN 
            Fact_Sale fs ON dt.time_id = fs.time_id
        WHERE
            (dt.year = 2023 AND dt.month BETWEEN 4 AND 12)
            OR (dt.year = 2024 AND dt.month BETWEEN 1 AND 3)
        GROUP BY 
            dt.month, dt.year
        ORDER BY 
            dt.year, dt.month
    `

	rows, err := con.Query(sqlStatement)
	if err != nil {
		return totalSalesResponse, err
	}
	defer rows.Close()

	index := 0
	for rows.Next() {
		var month string
		var salesAmount int
		if err := rows.Scan(&month, &salesAmount); err != nil {
			return totalSalesResponse, err
		}
		totalSalesResponse.Label[index] = month
		totalSalesResponse.Value[index] = salesAmount
		index++
	}

	if err := rows.Err(); err != nil {
		return totalSalesResponse, err
	}

	return totalSalesResponse, nil
}

func GetTotalItemsSold() (TotalItemsSold, error) {
	var totalItemsSoldResponse TotalItemsSold
	totalItemsSoldResponse.Length = 12
	totalItemsSoldResponse.Title = "Total Items Sold"
	totalItemsSoldResponse.Label = make([]string, 12)
	totalItemsSoldResponse.Value = make([]int, 12)

	con := db.CreateConDW()

	sqlStatement := `
		SELECT 
			CONCAT(SUBSTR(dt.month_name, 1, 3), ' ', dt.year) AS month,
			COALESCE(SUM(fs.total_item), 0) AS total_items_sold
		FROM 
			Dim_Time dt
		LEFT JOIN 
			Fact_Sale fs ON dt.time_id = fs.time_id
		WHERE
			(dt.year = 2023 AND dt.month BETWEEN 4 AND 12)
			OR (dt.year = 2024 AND dt.month BETWEEN 1 AND 3)
		GROUP BY 
			dt.month, dt.year
		ORDER BY 
			dt.year, dt.month;
    `

	rows, err := con.Query(sqlStatement)
	if err != nil {
		return totalItemsSoldResponse, err
	}
	defer rows.Close()

	index := 0
	for rows.Next() {
		var month string
		var totalItemsSold int
		if err := rows.Scan(&month, &totalItemsSold); err != nil {
			return totalItemsSoldResponse, err
		}
		totalItemsSoldResponse.Label[index] = month
		totalItemsSoldResponse.Value[index] = totalItemsSold
		index++
	}

	if err := rows.Err(); err != nil {
		return totalItemsSoldResponse, err
	}

	return totalItemsSoldResponse, nil
}

func GetTotalProfit() (TotalProfit, error) {
	var totalProfitResponse TotalProfit
	totalProfitResponse.Length = 12
	totalProfitResponse.Title = "Total Profit"
	totalProfitResponse.Label = make([]string, 12)
	totalProfitResponse.Value = make([]int, 12)

	con := db.CreateConDW()

	sqlStatement := `
		SELECT 
			CONCAT(SUBSTR(dt.month_name, 1, 3), ' ', dt.year) AS month,
			COALESCE(SUM(fs.profit), 0) AS total_profit
		FROM 
			Dim_Time dt
		LEFT JOIN 
			Fact_Sale fs ON dt.time_id = fs.time_id
		WHERE
			(dt.year = 2023 AND dt.month BETWEEN 4 AND 12)
			OR (dt.year = 2024 AND dt.month BETWEEN 1 AND 3)
		GROUP BY 
			dt.month, dt.year
		ORDER BY 
			dt.year, dt.month;
    `

	rows, err := con.Query(sqlStatement)
	if err != nil {
		return totalProfitResponse, err
	}
	defer rows.Close()

	index := 0
	for rows.Next() {
		var month string
		var totalProfit int
		if err := rows.Scan(&month, &totalProfit); err != nil {
			return totalProfitResponse, err
		}
		totalProfitResponse.Label[index] = month
		totalProfitResponse.Value[index] = totalProfit
		index++
	}

	if err := rows.Err(); err != nil {
		return totalProfitResponse, err
	}

	return totalProfitResponse, nil
}

func GetTopBestSellingProducts() (TopBestSellingProducts, error) {
	var topBestSellingProducts TopBestSellingProducts
	topBestSellingProducts.Length = 10
	topBestSellingProducts.Title = "Top Best Selling Products"
	topBestSellingProducts.Label = make([]string, 10)
	topBestSellingProducts.Value = make([]int, 10)

	con := db.CreateConDW()

	sqlStatement := `
		SELECT 
			dp.product_variant_name AS product_name,
			COALESCE(SUM(sd.quantity), 0) AS top_product
		FROM 
			dim_sale_detail sd
		JOIN 
			dim_product_variant dp ON sd.product_variant_id = dp.product_variant_id
		GROUP BY 
			dp.product_variant_name
		ORDER BY 
			top_product DESC
		LIMIT 10;
    `

	rows, err := con.Query(sqlStatement)
	if err != nil {
		return topBestSellingProducts, err
	}
	defer rows.Close()

	index := 0
	for rows.Next() {
		var productName string
		var topSaleProduct int
		if err := rows.Scan(&productName, &topSaleProduct); err != nil {
			return topBestSellingProducts, err
		}
		topBestSellingProducts.Label[index] = productName
		topBestSellingProducts.Value[index] = topSaleProduct
		index++
	}

	if err := rows.Err(); err != nil {
		return topBestSellingProducts, err
	}

	return topBestSellingProducts, nil
}
