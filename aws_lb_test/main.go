package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Listening on port 80")
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/ping", pingHandler)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	log.Fatalln(http.ListenAndServe(":80", nil))
}

func indexHandler(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, "This is AWS Load Balancer TEST")
}

func pingHandler(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, "OK")
}
