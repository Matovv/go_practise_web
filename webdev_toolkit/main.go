package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

const (
	cookie_session    = "session"
	cookie_expireTime = 1800
)

var tpl *template.Template

const port int = 8080

func init() {
	tpl = template.Must(tpl.ParseGlob("../templates/*.gohtml"))
}

func main() {
	fmt.Println("\nListening on port", port)

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/tool1", tool1Handler)
	http.HandleFunc("/ajax", ajaxHandler)
	http.HandleFunc("/context", contextHandler)
	http.HandleFunc("/ping", pingHandler)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func indexHandler(res http.ResponseWriter, req *http.Request) {
	cookie := getCookie(res, req)
	tpl.ExecuteTemplate(res, "index_toolkit.gohtml", cookie)
}

func tool1Handler(res http.ResponseWriter, req *http.Request) {
	s := fmt.Sprintf("\nHash1 - %s -\nHash2 - %s -\n", getCode("test@example.com"), getCode2("test@example.com"))
	tpl.ExecuteTemplate(res, "index_toolkit.gohtml", s)
}

func ajaxHandler(res http.ResponseWriter, req *http.Request) {
	log.Println("HAHA! You got ajax-ed!")
	s := "Here is your data, crybaby!"
	fmt.Fprintln(res,s)
}

func contextHandler(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	//ctx = context.WithValue(ctx, "userID", 777)
	//ctx = context.WithValue(ctx, "userName", "Bond")

	results, err := dbAccess(ctx)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(res, results)
}

func dbAccess(ctx context.Context) (int, error) {
	log.Println(ctx)
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	ch := make(chan int)
	go func(){
		//uid:=ctx.Value("userID").(int)
		uid:=777
		time.Sleep(10*time.Second)

		if ctx.Err() != nil {
			return
		}
		ch <- uid
	}()
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		case i := <-ch:
			return i, nil
		}
}

func pingHandler(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, "OK")
}

func getCode(s string) string {
	h := hmac.New(sha256.New, []byte("ourkey"))
	io.WriteString(h, s)
	return fmt.Sprintf("%x", h.Sum(nil))
}
func getCode2(s string) string {
	h := hmac.New(sha256.New, []byte("keyour"))
	io.WriteString(h, s)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func getCookie(res http.ResponseWriter, req *http.Request) *http.Cookie {
	cookie, err := req.Cookie(cookie_session)
	if err != nil {
		sID, err := uuid.NewV7()
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return cookie
		}
		cookie = &http.Cookie{
			Name:  cookie_session,
			Value: sID.String(),
		}
	}
	cookie.MaxAge = cookie_expireTime
	http.SetCookie(res, cookie)
	return cookie

}
