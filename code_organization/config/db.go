// In this config we initialize our database as soon as service starts

package config

import (
	"database/sql"
	"fmt"
)

var DB *sql.DB

func init() {
	var err error
	const dbPath = "admin:sandstone@tcp(elnurdbmysql.cxygyk82uwas.eu-north-1.rds.amazonaws.com)/test01?charset=utf8"
	DB, err = sql.Open("mysql", dbPath)
	if err != nil {
		panic(err)
	}

	if err = DB.Ping(); err != nil {
		panic(err)
	}

	fmt.Println("Connected to Database!")
}
