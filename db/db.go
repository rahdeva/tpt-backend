package db

import (
	"database/sql"
	"fmt"
	"tpt_backend/config"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var err error

func Init() {
	conf := config.GetConfig()

	// username:password@protocol(address)/dbname?param=value
	connectionString := conf.DB_USERNAME + ":" + conf.DB_PASSWORD + "@tcp(" + conf.DB_HOST + ":" + conf.DB_PORT + ")/" + conf.DB_NAME

	db, err = sql.Open("mysql", connectionString+"?parseTime=true")

	if err != nil {
		panic("Connection Error")
	}

	err := db.Ping()

	if err != nil {
		fmt.Println(connectionString)
		panic("DSN Error")
	}
}

func CreateCon() *sql.DB {
	return db
}
