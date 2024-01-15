package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const port int = 8080

func main() {
	fmt.Println("Listening on port", port)

	http.HandleFunc("/", formHandler)
	http.HandleFunc("/upload", uploadHandler)
	http.Handle("/favicon.ico", http.NotFoundHandler())

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func formHandler(res http.ResponseWriter, req *http.Request) {
	// when you call, for example, "localhost:8080/?q=dog"
	// it will read that q element's value and assign it to v
	// in case q element does not exist, then it will assing empty string to v
	v := req.FormValue("q")
	//io.WriteString(res, "Do my search: "+v)  // example result will be: Do my search: dog

	// write header that will specify response as html element, instead of plain text
	// otherwise form will not work, and will be just a string of text on browser
	res.Header().Set("Content-Type", "text/html; charset=utf-8")

	io.WriteString(res, `
	<form method="post">
		<input type="text" name="q">
		<input type="submit">
	</form>
	<br>`+v)
}

// uploadHandler handles the upload file request
func uploadHandler(res http.ResponseWriter, req *http.Request) {
	var s string
	fmt.Println(req.Method)
	if req.Method == http.MethodPost {
		// open
		f, h, err := req.FormFile("q")
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		defer f.Close()

		// for your information
		fmt.Println("\nfilename:", h.Filename, "\nerr:", err)

		// read
		bs, err := io.ReadAll(f)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		s = string(bs)

		// store on the server
		dst, err := os.Create(filepath.Join("./assets/", h.Filename))
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = dst.Write(bs)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}

		log.Printf("File <%s> have been uploaded to server!",h.Filename)
	}

	res.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(res, `
	<form method="POST" enctype="multipart/form-data">
		<input type="file" name="q">
		<input type="submit">
	</form>
	<br>`+s)
}
