package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

const port int = 8080

func main() {
	fmt.Println("Listening on port", port)

	http.HandleFunc("/", incrementHandler)
	http.HandleFunc("/set", setHandler)
	http.HandleFunc("/read", readHandler)
	http.HandleFunc("/expire", expireHandler)
	http.Handle("/favicon.ico", http.NotFoundHandler())

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

// set cookie
func setHandler(res http.ResponseWriter, req *http.Request) {
	http.SetCookie(res, &http.Cookie{
		Name:  "my-cookie",
		Value: "some value",
	})
	fmt.Fprintln(res, "COOKIE IS WRITTEN - CHECK YOUR BROWSER")
	fmt.Fprintln(res, "in chrome go to: dev tools / application / cookies")
}

// read cookie
func readHandler(res http.ResponseWriter, req *http.Request) {
	c, err := req.Cookie("my-cookie")
	if err != nil {
		http.Error(res, err.Error(), http.StatusNotFound)
		return
	}

	fmt.Fprintln(res, "YOUR COOKIE:", c)
}

// expires and deletes cookie, if cookie doesnot exist then redirects to set
func expireHandler(res http.ResponseWriter, req *http.Request) {
	c, err := req.Cookie("my-cookie")
	if err != nil {
		http.Redirect(res,req,"/set",http.StatusSeeOther)
		return
	}

	c.MaxAge = -1
	http.SetCookie(res, c)
	http.Redirect(res,req,"/",http.StatusSeeOther)
}

// increments cookie's value and show it to user
func incrementHandler(res http.ResponseWriter, req *http.Request) {
	c, err := req.Cookie("dynamic-cookie")
	if err == http.ErrNoCookie {
		c = &http.Cookie{
			Name:  "dynamic-cookie",
			Value: "0",
		}
	}

	count, _ := strconv.Atoi(c.Value)
	if err != nil {
		log.Fatalln(err)
	}
	count++
	c.Value = strconv.Itoa(count)

	http.SetCookie(res, c)
	io.WriteString(res, c.Value)

}
