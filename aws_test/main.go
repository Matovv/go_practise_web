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
	log.Fatalln(http.ListenAndServe(":80", nil))
}

func indexHandler(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, "Oh Yeah! I'm running on AWS.")
}
