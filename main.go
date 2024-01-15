package main

import (
	"fmt"
	"log"
	"net/http"
)

const port int = 8080

func main() {
	fmt.Println("Listening on port", port)
	
	http.Handle("/favicon.ico", http.NotFoundHandler())

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
