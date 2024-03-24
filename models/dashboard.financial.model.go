package models

import (
	"tpt_backend/db"
)

type CashFlow struct {
	Length int      `json:"length"`
	Title  string   `json:"title"`
	Label  []string `json:"label"`
	Value  []int    `json:"value"` // Change the type to []int
}

type FinancialTransaction struct {
	Length int      `json:"length"`
	Title  string   `json:"title"`
	Labels []string `json:"labels"`
	Values []int    `json:"values"`
}

type CashIn struct {
	Length int      `json:"length"`
	Title  string   `json:"title"`
	Labels []string `json:"labels"`
	Values []int    `json:"values"`
}

type CashOut struct {
	Length int      `json:"length"`
	Title  string   `json:"title"`
	Labels []string `json:"labels"`
	Values []int    `json:"values"`
}

func GetCashFlow() (CashFlow, error) {
	var cashFlow CashFlow
	cashFlow.Length = 2
	cashFlow.Title = "Uang Masuk vs Uang Keluar"
	cashFlow.Label = []string{"Uang Masuk", "Uang Keluar"}
	cashFlow.Value = make([]int, 2) // Change the type to []int

	// Retrieve cash flow information
	con := db.CreateConDW()

	sqlStatement := `
		SELECT
			CEILING(SUM(cash_in)) AS total_cash_in,
			SUM(cash_out) AS total_cash_out
		FROM
			fact_financial;
    `

	err := con.QueryRow(sqlStatement).Scan(&cashFlow.Value[0], &cashFlow.Value[1])
	if err != nil {
		return cashFlow, err
	}

	// Calculate percentages
	total := cashFlow.Value[0] + cashFlow.Value[1]
	cashFlow.Value[0] = (cashFlow.Value[0] * 100) / total
	cashFlow.Value[1] = (cashFlow.Value[1] * 100) / total

	return cashFlow, nil
}

func GetFinancialByType() (FinancialTransaction, error) {
	var financialTransactions FinancialTransaction
	financialTransactions.Length = 0
	financialTransactions.Title = "Jenis Keuangan"

	con := db.CreateConDW()

	sqlStatement := `
        SELECT
            dft.financial_type_name,
            COALESCE(SUM(ff.cash_in), 0) + COALESCE(SUM(ff.cash_out), 0) AS total_cash
        FROM
            dim_financial_type dft
        LEFT JOIN
            fact_financial ff ON dft.financial_type_id = ff.financial_type_id
        GROUP BY
            dft.financial_type_name
        ORDER BY
            dft.financial_type_id;
    `

	rows, err := con.Query(sqlStatement)
	if err != nil {
		return financialTransactions, err
	}
	defer rows.Close()

	labels := make([]string, 0)
	values := make([]int, 0)

	for rows.Next() {
		var item struct {
			Type  string `json:"type"`
			Total int    `json:"total"`
		}
		if err := rows.Scan(&item.Type, &item.Total); err != nil {
			return financialTransactions, err
		}
		labels = append(labels, item.Type)
		values = append(values, item.Total)
		financialTransactions.Length++
	}

	if err := rows.Err(); err != nil {
		return financialTransactions, err
	}

	financialTransactions.Labels = labels
	financialTransactions.Values = values

	return financialTransactions, nil
}

func GetCashIn() (CashIn, error) {
	var cashIn CashIn
	cashIn.Length = 0
	cashIn.Title = "Kas Masuk"

	con := db.CreateConDW()

	sqlStatement := `
        SELECT
            CONCAT(SUBSTR(dt.month_name, 1, 3), ' ', dt.year) AS month,
            COALESCE(SUM(ff.cash_in), 0) AS total_cash_in
        FROM
            dim_time dt
        LEFT JOIN
            fact_financial ff ON dt.time_id = ff.time_id
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
		return cashIn, err
	}
	defer rows.Close()

	labels := make([]string, 0)
	values := make([]int, 0)

	for rows.Next() {
		var month string
		var totalCashIn int
		if err := rows.Scan(&month, &totalCashIn); err != nil {
			return cashIn, err
		}
		labels = append(labels, month)
		values = append(values, totalCashIn)
		cashIn.Length++
	}

	if err := rows.Err(); err != nil {
		return cashIn, err
	}

	cashIn.Labels = labels
	cashIn.Values = values

	return cashIn, nil
}

func GetCashOut() (CashOut, error) {
	var cashOut CashOut
	cashOut.Length = 0
	cashOut.Title = "Kas Keluar"

	con := db.CreateConDW()

	sqlStatement := `
        SELECT
            CONCAT(SUBSTR(dt.month_name, 1, 3), ' ', dt.year) AS month,
            COALESCE(SUM(ff.cash_out), 0) AS total_cash_out
        FROM
            dim_time dt
        LEFT JOIN
            fact_financial ff ON dt.time_id = ff.time_id
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
		return cashOut, err
	}
	defer rows.Close()

	labels := make([]string, 0)
	values := make([]int, 0)

	for rows.Next() {
		var month string
		var totalCashOut int
		if err := rows.Scan(&month, &totalCashOut); err != nil {
			return cashOut, err
		}
		labels = append(labels, month)
		values = append(values, totalCashOut)
		cashOut.Length++
	}

	if err := rows.Err(); err != nil {
		return cashOut, err
	}

	cashOut.Labels = labels
	cashOut.Values = values

	return cashOut, nil
}
