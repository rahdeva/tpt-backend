package models

import (
	"time"
	"tpt_backend/db"
)

type Home struct {
	TotalProduk           int `json:"total_produk"`
	TotalKategori         int `json:"total_kategori"`
	TotalSupplier         int `json:"total_supplier"`
	TotalUser             int `json:"total_user"`
	TotalPenjualanHariIni int `json:"total_penjualan_hari_ini"`
	TotalPembelianHariIni int `json:"total_pembelian_hari_ini"`
}

func GetHomeData() (Response, error) {
	var home Home
	var res Response

	con := db.CreateCon()

	sqlStatement := `
	SELECT 
	(SELECT COUNT(*) FROM product_variant) AS total_produk,
	(SELECT COUNT(*) FROM category) AS total_kategori,
	(SELECT COUNT(*) FROM supplier) AS total_supplier,
	(SELECT COUNT(*) FROM user) AS total_user,
	(SELECT COALESCE(SUM(total_price), 0) FROM sale WHERE sale_date BETWEEN ? AND ?) AS total_penjualan_hari_ini,
	(SELECT COALESCE(SUM(total_price), 0) FROM purchase WHERE purchase_date BETWEEN ? AND ?) AS total_pembelian_hari_ini
	`

	startOfDay := time.Now().UTC().Truncate(24 * time.Hour)
	endOfDay := startOfDay.Add(24 * time.Hour).Add(-time.Second)

	row := con.QueryRow(sqlStatement, startOfDay, endOfDay, startOfDay, endOfDay)

	err := row.Scan(
		&home.TotalProduk,
		&home.TotalKategori,
		&home.TotalSupplier,
		&home.TotalUser,
		&home.TotalPenjualanHariIni,
		&home.TotalPembelianHariIni,
	)

	if err != nil {
		return res, err
	}

	res.Data = map[string]interface{}{
		"home": home,
	}

	return res, nil
}
