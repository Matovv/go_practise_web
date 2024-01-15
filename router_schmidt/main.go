package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Hotdog struct {
	CustomerName string
	Type         string
	Cost         string
}

type Car struct {
	CustomerName string
	Brand        string
	Model        string
	Color        string
	Cost         string
}

func orderHotdog(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	//fmt.Fprintf(w, "Hotdog '%s' for customer '%s' is served. Transaction amount is '%s'. Thank you!", m.Type, m.CustomerName, m.Cost)
	err := req.ParseForm()
	if err != nil {
		log.Fatalln(err)
	}
	customerName := params.ByName("fname")
	var data Hotdog
	if customerName != "" {
		data = Hotdog{customerName, "Big Joe", "$2.95"}
	}
	res.Header().Set("Order-Key", "000223")
	err = tpl.ExecuteTemplate(res, "post_test.gohtml", data)
	if err != nil {
		log.Println(err)
	}
}

func orderCar(res http.ResponseWriter, req *http.Request, params httprouter.Params) {

	err := req.ParseForm()
	if err != nil {
		log.Fatalln(err)
	}
	customerName := params.ByName("fname")

	log.Println("Car is requested  by -", customerName)

	var data Car
	if customerName != "" {
		data = Car{customerName, "BMW", "X7", "Black", "$30000"}
	}
	res.Header().Set("Order-Key", "000155")
	err = tpl.ExecuteTemplate(res, "post_test2.gohtml", data)
	if err != nil {
		log.Println(err)
	}

	log.Println("Car Served! Customer -", customerName)
}

var tpl *template.Template

func init() {
	tpl = template.Must(tpl.ParseGlob("templates/*.gohtml"))
}

const port int = 8080

func main() {
	fmt.Println("Listening on port", port)

	router := httprouter.New()
	router.GET("/order/hotdog", orderHotdog)
	router.GET("/order/car", orderCar)

	log.Fatal(http.ListenAndServe(getPort(), router))

}

// getPort returns correctly formatted port string for listening
func getPort() string {
	return fmt.Sprintf(":%d", port)
}
