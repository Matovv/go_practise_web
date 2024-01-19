package main

import (
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func main() {
	router := httprouter.New()
	router.GET("/", indexHandler)
	http.ListenAndServe(":80", router)

}

func indexHandler(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	io.WriteString(res, "Hello from a Docker Container!")
}
