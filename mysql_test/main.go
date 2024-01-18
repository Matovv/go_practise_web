package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

const dbPath = "admin:sandstone@tcp(elnurdbmysql.cxygyk82uwas.eu-north-1.rds.amazonaws.com)/test01?charset=utf8"

var db *sql.DB
var err error

const port int = 80

func main() {
	fmt.Println("Connecting to database ...")
	db, err = sql.Open("mysql", dbPath)
	check(err)
	defer db.Close()

	err = db.Ping()
	check(err)

	fmt.Println("CONNECTED!")
	fmt.Println("\nListening on port", port)

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/adduser", addUserHandler)
	http.Handle("/favicon.ico", http.NotFoundHandler())

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func indexHandler(res http.ResponseWriter, req *http.Request) {
	_, err :=io.WriteString(res, "Successful Connection to Database!")
	check(err)
	rows, err := db.Query(`SELECT username FROM users`)
	var s, username string
	s = "\nRetrieved Data:\n"
	for rows.Next() {
		err = rows.Scan(&username)
		check(err)
		s += username + "\n"
	}
	check(err)

	_, err =io.WriteString(res, s)
	check(err)
}

func addUserHandler(res http.ResponseWriter, req *http.Request) {
	stmt, err := db.Prepare(`INSERT INTO users (username) VALUE ("Fire Hose");`)
	check(err)
	result, err := stmt.Exec()
	check(err)
	n, err := result.RowsAffected()
	check(err)
	fmt.Fprintln(res,"INSERTED RECORD", n)

}

// check takes error and handles is
func check(err error) {
	if err != nil {
		log.Println(err)
	}
}
