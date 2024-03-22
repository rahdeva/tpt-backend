package db

import (
	"database/sql"
	"fmt"
	"tpt_backend/config"

	_ "github.com/go-sql-driver/mysql"
)

var dw *sql.DB
var err2 error

func DWInit() {

	DWconf := config.GetDWConfig()

	// username:password@protocol(address)/dbname?param=value
	DWconnectionString := DWconf.DB_USERNAME + ":" + DWconf.DB_PASSWORD + "@tcp(" + DWconf.DB_HOST + ":" + DWconf.DB_PORT + ")/" + DWconf.DB_NAME

	dw, err2 = sql.Open("mysql", DWconnectionString+"?parseTime=true")

	if err2 != nil {
		panic("Connection Error")
	}

	err2 := dw.Ping()

	if err2 != nil {
		fmt.Println(DWconnectionString)
		panic("DSN Error")
	}
}

func CreateConDW() *sql.DB {
	return dw
}
