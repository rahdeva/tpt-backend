package models

import "tpt_backend/db"

type TotalSalesResponse struct {
	Length int      `json:"length"`
	Title  string   `json:"title"`
	Label  []string `json:"label"`
	Value  []int    `json:"value"`
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
