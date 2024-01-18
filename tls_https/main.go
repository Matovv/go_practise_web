package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

const port int = 10443

func main() {
	fmt.Println("Listening on port", port)

	http.HandleFunc("/", indexHandler)
	http.Handle("/favicon.ico", http.NotFoundHandler())

	log.Fatal(http.ListenAndServeTLS(fmt.Sprintf(":%d", port), "cert.pem", "key.pem", nil))
}

func indexHandler(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, "TSL and HTTPS Test")
}


// generating unsigned certificate
// go run $(go env GOROOT)/src/crypto/tls/generate_cert.go --host=somedomainname.com (for example: localhost)
