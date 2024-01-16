package main

import (
	"fmt"
	"log"
	"net/http"
	"html/template"

	"github.com/google/uuid"
)

type user struct {
	UserName string
	First    string
	Last     string
}

var tpl *template.Template
var dbUsers = map[string]user{}
var dbSessions = map[string]string{}

func init() {
	tpl = template.Must(tpl.ParseGlob("templates/*.gohtml"))
}

const port int = 8080

func main() {
	fmt.Println("Listening on port", port)

	http.HandleFunc("/", sessionHandler)
	http.HandleFunc("/profile", userProfileHandler)
	http.Handle("/favicon.ico", http.NotFoundHandler())

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

// sessionHandler checks if session cookie is available, in this case
// uses value in cookie to get user from user database by looking up the session database to find it's name.
// Creates a new session cookie if one does not exist.
// If used as POST then creates new user and puts it to user database and session database.
func sessionHandler(res http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie("session")
	if err != nil {
		sessionID, err := uuid.NewV7()
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		cookie = &http.Cookie{
			Name:  "session",
			Value: sessionID.String(),
			// Secure: true,  // when this is true, you can only send it through https, which makes it more secure
			HttpOnly: true, // when this is true, you cant access this cookie with javascript, which makes it secure
		}
		http.SetCookie(res, cookie)
	}

	var u user
	if userName, ok := dbSessions[cookie.Value]; ok {
		u = dbUsers[userName];
	}

	if req.Method == http.MethodPost {
		userName := req.FormValue("username")
		first := req.FormValue("firstname")
		last := req.FormValue("lastname")
		u = user{userName,first,last}
		dbSessions[cookie.Value] = userName
		dbUsers[userName] = u
	}

	tpl.ExecuteTemplate(res, "session_form.gohtml", u)

}


// userProfileHandler shows user profile if he exist and has session cookie.
// Otherwise redirects to main page
func userProfileHandler(res http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie("session")
	if err != nil {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	userName, ok := dbSessions[cookie.Value]
	if !ok {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	user, ok := dbUsers[userName]
	if !ok {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	tpl.ExecuteTemplate(res, "user_profile.gohtml", user)

}
