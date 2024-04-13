package models

import (
	"math"
	"strconv"
	"strings"
	"time"
	"tpt_backend/db"
)

type ForecastDetail struct {
	ForecastDetailID        int       `json:"forecast_detail_id"`
	Title                   string    `json:"title"`
	Method                  string    `json:"method"`
	TimeID                  int       `json:"time_id"`
	ForecastDate            string    `json:"forecast_date"`
	LastWeekDate            string    `json:"last_week_date"`
	LastWeekTransaction     int       `json:"last_week_transaction"`
	MAE                     float64   `json:"mae"`
	MSE                     float64   `json:"mse"`
	RSME                    float64   `json:"rsme"`
	Length                  int       `json:"length"`
	Week                    []int     `json:"week"`
	WeekDate                []string  `json:"week_date"`
	ActualTotalTransaction  []int     `json:"actual_total_transaction"`
	PredictTotalTransaction []float64 `json:"predict_total_transaction"`
	Label                   []string  `json:"label"`
}

func GetSaleForecast() (ForecastDetail, error) {
	var forecastDetail ForecastDetail
	var weekStr, weekDateStr, actualTotalTransactionStr, predictTotalTransactionStr string

	con := db.CreateConDW()

	sqlStatement := `
        SELECT
            dfd.forecast_detail_id,
            dfd.method,
			fsf.time_id,
            dfd.last_week_date,
            dfd.last_week_transaction,
            dfd.mae,
            dfd.mse,
            dfd.rsme,
            COUNT(fsf.week) AS length,
            GROUP_CONCAT(fsf.week) AS week,
            GROUP_CONCAT(fsf.week_date) AS week_date,
            GROUP_CONCAT(fsf.actual_total_transaction) AS actual_total_transaction,
            GROUP_CONCAT(fsf.predict_total_transaction) AS predict_total_transaction
        FROM
            dim_forecast_detail dfd
        JOIN
            fact_sale_forecast fsf ON dfd.forecast_detail_id = fsf.forecast_detail_id
        WHERE
            dfd.forecast_detail_id = (SELECT MAX(forecast_detail_id) FROM dim_forecast_detail)
        GROUP BY
            dfd.forecast_detail_id
    `

	err := con.QueryRow(sqlStatement).Scan(
		&forecastDetail.ForecastDetailID,
		&forecastDetail.Method,
		&forecastDetail.TimeID,
		&forecastDetail.LastWeekDate,
		&forecastDetail.LastWeekTransaction,
		&forecastDetail.MAE,
		&forecastDetail.MSE,
		&forecastDetail.RSME,
		&forecastDetail.Length,
		&weekStr,
		&weekDateStr,
		&actualTotalTransactionStr,
		&predictTotalTransactionStr,
	)
	if err != nil {
		return forecastDetail, err
	}

	forecastDetail.Title = "Forecasting Transaksi Penjualan 3 Bulan Mendatang"
	forecastDetail.ForecastDate = convertTimeIDToDate(forecastDetail.TimeID)

	forecastDetail.Week, err = parseIntArray(weekStr)
	if err != nil {
		return forecastDetail, err
	}

	forecastDetail.WeekDate, err = parseStringArray(weekDateStr)
	if err != nil {
		return forecastDetail, err
	}

	forecastDetail.ActualTotalTransaction, err = parseIntArray(actualTotalTransactionStr)
	if err != nil {
		return forecastDetail, err
	}

	forecastDetail.PredictTotalTransaction, err = parseDoubleArray(predictTotalTransactionStr)
	if err != nil {
		return forecastDetail, err
	}

	forecastDetail.Label = generateLabels(forecastDetail.WeekDate)

	forecastDetail.MAE = roundToTwoDecimalPlaces(forecastDetail.MAE)
	forecastDetail.MSE = roundToTwoDecimalPlaces(forecastDetail.MSE)
	forecastDetail.RSME = roundToTwoDecimalPlaces(forecastDetail.RSME)

	return forecastDetail, nil
}

func parseIntArray(str string) ([]int, error) {
	var array []int
	if str != "" {
		// Parse string to int array
		values := strings.Split(str, ",")
		for _, value := range values {
			floatValue, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return nil, err
			}
			intValue := int(floatValue) // Konversi nilai float menjadi integer
			array = append(array, intValue)
		}
	}
	return array, nil
}

func parseStringArray(str string) ([]string, error) {
	var array []string
	if str != "" {
		// Parse string to string array
		array = strings.Split(str, ",")
	}
	return array, nil
}

func parseDoubleArray(str string) ([]float64, error) {
	var array []float64
	if str != "" {
		// Parse string to float64 array
		values := strings.Split(str, ",")
		for _, value := range values {
			floatValue, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return nil, err
			}
			// Ambil 2 angka belakang koma
			roundedValue := math.Round(floatValue*100) / 100
			array = append(array, roundedValue)
		}
	}
	return array, nil
}

func roundToTwoDecimalPlaces(value float64) float64 {
	return math.Round(value*100) / 100
}

func convertTimeIDToDate(timeID int) string {
	layout := "20060102"
	timeStr := strconv.Itoa(timeID)
	timeValue, _ := time.Parse(layout, timeStr)
	return timeValue.Format("2006-01-02")
}

func generateLabels(weekDate []string) []string {
	var labels []string

	startIndex := len(weekDate) - 1
	for i := startIndex; i >= 1; i -= 12 {
		labels = append(labels, weekDate[i])
	}

	// Memanggil fungsi reverseSlice untuk membalikkan slice labels
	return reverseSlice(labels)
}

// Fungsi untuk membalikkan slice
func reverseSlice(s []string) []string {
	for i := 0; i < len(s)/2; i++ {
		j := len(s) - i - 1
		s[i], s[j] = s[j], s[i]
	}
	return s
}
