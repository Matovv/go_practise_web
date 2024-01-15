package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

const port int = 8080

var tpl *template.Template

func init() {
	tpl = template.Must(tpl.ParseGlob("templates/*.gohtml"))
}

func main() {
	fmt.Println("Listening on port", port)

	http.Handle("/", http.FileServer(http.Dir("./assets"))) // serves list of files and folders inside specified folder path
	http.HandleFunc("/wow", wowHandler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func wowHandler(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, "./assets/wow.png") // serves specified file to response
}
