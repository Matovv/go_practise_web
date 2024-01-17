package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	session           = "session"
	sessionLength	  = 600  // seconds
	accessLevel_admin = "ADMIN"
	accessLevel_user  = "USER"
)

// ANSI color codes
const (
	reset   = "\033[0m"
	red     = "\033[31m"
	green   = "\033[32m"
	yellow  = "\033[33m"
	blue    = "\033[34m"
	magenta = "\033[35m"
	cyan    = "\033[36m"
	white   = "\033[37m"
)

type User struct {
	UserName    string
	Password    []byte
	First       string
	Last        string
	AccessLevel string
}

type Session struct {
	UserName     string
	LastActivity time.Time
}

var tpl *template.Template
var dbUsers = map[string]User{}
var dbSessions = map[string]Session{}

func createMockUser(email string, pass string, fname string, lname string, accessLevel string) {
	mockPass, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.MinCost)
	if err != nil {
		log.Fatalln(err)
	}
	mockUser := User{email, mockPass, "System", "Admin", accessLevel}
	dbUsers[mockUser.UserName] = mockUser
}

func init() {
	tpl = template.Must(tpl.ParseGlob("templates/*.gohtml"))
	createMockUser("admin@x.com", "qwerty", "Sys", "Admin", accessLevel_admin)
	
}

const port int = 80

func main() {
	fmt.Println("Listening on port", port)

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/signup", signupHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/profile", userProfileHandler)
	http.HandleFunc("/admin", adminHandler)
	http.Handle("/favicon.ico", http.NotFoundHandler())

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func indexHandler(res http.ResponseWriter, req *http.Request) {
	if alreadyLoggedIn(req) {
		user, ok := getUser(res, req)
		if !ok {
			http.Error(res, "User not found in database", http.StatusInternalServerError)
			return
		}
		tpl.ExecuteTemplate(res, "index.gohtml", user)
		return
	}
	tpl.ExecuteTemplate(res, "index.gohtml", nil)
}

// signupHandler handles user signup, if already logged then redirects to main page
func signupHandler(res http.ResponseWriter, req *http.Request) {

	// check whether user is already logged in
	if alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	// if method POST then proceed to user creation
	// otherwise just open sign up page
	if req.Method == http.MethodPost {

		userName := req.FormValue("username")
		password := req.FormValue("password")
		first := req.FormValue("firstname")
		last := req.FormValue("lastname")

		// check fields to not be empty
		if userName == "" || password == "" || first == "" || last == "" {
			http.Error(res, "Empty fields are not allowed!", http.StatusForbidden)
			return
		}

		// check username to not be already taken
		if _, ok := dbUsers[userName]; ok {
			http.Error(res, "Username already taken!", http.StatusForbidden)
			return
		}

		// set new session
		sID, err := uuid.NewV7()
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		cookie := &http.Cookie{
			Name:  "session",
			Value: sID.String(),
		}
		http.SetCookie(res, cookie)
		dbSessions[cookie.Value] = Session{userName,time.Now()}

		// hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// create user
		user := User{userName, hashedPassword, first, last, accessLevel_user}
		dbUsers[userName] = user

		log.Println("User -", user.UserName, "- signed up.")
		log.Println("Full User Data: ", user)

		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	tpl.ExecuteTemplate(res, "signup.gohtml", nil)

}

// loginHandler handles user login process
func loginHandler(res http.ResponseWriter, req *http.Request) {

	// check whether user is already logged in
	if alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	// if method POST then proceed to login process
	// otherwise just open login page
	if req.Method == http.MethodPost {
		userName := req.FormValue("username")
		password := req.FormValue("password")

		// check fields to not be empty
		if userName == "" || password == "" {
			http.Error(res, "Empty fields are not allowed!", http.StatusForbidden)
			return
		}

		log.Println(userName, "tried to login!")
		// check whether such username exist
		user, ok := dbUsers[userName]
		if !ok {
			log.Println(userName, "does not exist!")
			http.Error(res, "Username and/or password do not match!", http.StatusForbidden)
			return
		}
		// check whether passwords match
		err := bcrypt.CompareHashAndPassword(user.Password, []byte(password))
		if err != nil {
			log.Println(userName, "entered wrong password!")
			http.Error(res, "Username and/or password do not match!", http.StatusForbidden)
			return
		}
		// create session
		sID, err := uuid.NewV7()
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		cookie := &http.Cookie{
			Name:  session,
			Value: sID.String(),
		}
		cookie.MaxAge = sessionLength
		http.SetCookie(res, cookie)

		dbSessions[cookie.Value] = Session{userName,time.Now()}

		log.Println(userName, "logged in!")
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	tpl.ExecuteTemplate(res, "login.gohtml", nil)
}

// logoutHandler handles user logout process
func logoutHandler(res http.ResponseWriter, req *http.Request) {
	// redirect to main page if not logged in yet
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	cookie, _ := req.Cookie(session)
	// delete session
	log.Println(dbSessions[cookie.Value].UserName, "logged out!")
	delete(dbSessions, cookie.Value)
	// remove (expire) the cookie
	cookie = &http.Cookie{
		Name:   session,
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(res, cookie)

	cleanSessions()

	http.Redirect(res, req, "/", http.StatusSeeOther)
}

// userProfileHandler shows user profile if he exist and has session cookie.
// Otherwise redirects to main page
func userProfileHandler(res http.ResponseWriter, req *http.Request) {

	// check whether user is already logged in, if not, then redirect to home page
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	// get user
	user, ok := getUser(res, req)
	if !ok {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	log.Println("User -", user.UserName, "- visited his profile.")

	// open profile page and pass user data to it
	tpl.ExecuteTemplate(res, "user_profile.gohtml", user)

}

func adminHandler(res http.ResponseWriter, req *http.Request) {
	//check if user logged in
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	user, ok := getUser(res, req)
	if !ok {
		http.Error(res, "User not found!", http.StatusInternalServerError)
		return
	}

	log.Println("User -", user.UserName, "- visited admin page.")

	if !checkAccess(user, accessLevel_admin) {
		http.Error(res, "ACCESS RESTRICTED!", http.StatusForbidden)
		return
	}

	tpl.ExecuteTemplate(res, "admin.gohtml", user)

}

// alreadyLoggedIn checks whether user is already logged or not through session look up
func alreadyLoggedIn(req *http.Request) bool {
	// get cookie
	cookie, err := req.Cookie(session)
	if err != nil {
		return false
	}
	// if cookie value corresponds to user's session id, return true
	userName := dbSessions[cookie.Value].UserName
	_, ok := dbUsers[userName]
	return ok

}

// checkAccess check access level of request user and if it matches, returns true
func checkAccess(user User, accessLevel string) bool {
		log.Println("User -", user.UserName, "- tries to access restricted area!")
		if user.AccessLevel == accessLevel {
			log.Println("User -", user.UserName, "-", green, "ACCESS GRANTED!", reset)
			return true
		}
		log.Println("User -", user.UserName, "-", red, "ACCESS DENIED!", reset)
		return false
}

// getUser returns the requested user and confirmation from database
func getUser(res http.ResponseWriter, req *http.Request) (User, bool) {
	var user User
	// get cookie
	cookie, err := req.Cookie(session)
	if err != nil {
		sID, err := uuid.NewV7()
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return user, false
		}
		cookie = &http.Cookie{
			Name:  session,
			Value: sID.String(),
		}
	}
	cookie.MaxAge = sessionLength
    // set the max age outside of cookie creation above, otherwise it will only be set 
	// on creating, without ability to update it when it exists. 

	http.SetCookie(res, cookie)

	// if the user exists already, return user and true. Otherwise return empty user and false
	if userSession, ok := dbSessions[cookie.Value]; ok {
		user = dbUsers[userSession.UserName]
	} else {
		return user, false
	}

	return user, true

}

func cleanSessions() {
	log.Println("Cleaning sessions...")
	count := 0
	for k, v:= range dbSessions {
		if(time.Since(v.LastActivity) > (time.Second * sessionLength)) {
			delete(dbSessions, k)
			count++
		}
	}
	log.Println("Total  ",count,"  sessions cleaned!")
}

