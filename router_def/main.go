package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type Hotdog struct {
	CustomerName string
	Type string
	Cost string
}

type Car struct {
	CustomerName string
	Brand string
	Model string
	Color string
	Cost string
}

func orderHotdogHandler(res http.ResponseWriter, req *http.Request) {
	//fmt.Fprintf(w, "Hotdog '%s' for customer '%s' is served. Transaction amount is '%s'. Thank you!", m.Type, m.CustomerName, m.Cost)
	err := req.ParseForm()
	if err != nil {
		log.Fatalln(err)
	}
	customerName := req.Form.Get("fname")
	var data Hotdog
	if customerName != "" {
		data = Hotdog{req.Form.Get("fname"),"Big Joe", "$2.95"}
	}
	res.Header().Set("Order-Key", "000223")
	err = tpl.ExecuteTemplate(res,"post_test.gohtml", data)
	if err != nil {
		log.Println(err)
	}
}

func orderCarHandler(res http.ResponseWriter, req *http.Request) {
	
	err := req.ParseForm()
	if err != nil {
		log.Fatalln(err)
	}
	customerName := req.Form.Get("fname")

	log.Println("Car is requested  by -",customerName)

	var data Car
	if customerName != "" {
		data = Car{req.Form.Get("fname"),"BMW", "X7", "Black","$30000"}
	}
	res.Header().Set("Order-Key", "000155")
	err = tpl.ExecuteTemplate(res,"post_test2.gohtml", data)
	if err != nil {
		log.Println(err)
	}

	log.Println("Car Served! Customer -",customerName)
}

var tpl *template.Template
func init() {
	tpl = template.Must(tpl.ParseGlob("templates/*.gohtml"))
}

const port int = 8080

func main() {
	fmt.Println("Listening on port",port)

	http.HandleFunc("/order/hotdog", orderHotdogHandler)
	http.HandleFunc("/order/car", orderCarHandler)

	//using http.Handle with handle func, instead of http.HandleFunc():
	//http.Handle("/", http.HandlerFunc(<any func with ServeHTTP arguments>))

	log.Fatal(http.ListenAndServe(getPort(), nil))
	

}

// getPort returns correctly formatted port string for listening
func getPort() string {
	return fmt.Sprintf(":%d",port)
}


