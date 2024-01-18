package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Matovv/go_practise_web/mvc_mongo/models"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/mongo"
)

// Declare controller
type UserController struct {
	session *mongo.Client
}

// NewUserController takes mongodb session pointer argument and sets it to User Controller pointer
// and then returns the pointer to our UserController
func NewUserController(session *mongo.Client) *UserController {
	return &UserController{session}
}

// GetUser handles GET user request, and returns user by the id param
// if url queries ('name' in this case) specified then does additional functions
func (controller UserController) GetUser(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	user := models.User{
		ID:     params.ByName("id"),
		Name:   "Jhon Wicked",
		Gender: "Male",
		Age:    32,
	}

	query_name := req.URL.Query().Get("name")
	if query_name != "" {
		user.Name = query_name
	}

	// Marshal to json
	json, err := json.Marshal(user)
	if err != nil {
		log.Println(err)
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK) // 200
	fmt.Fprintf(res, "User: %s\n", json)
}

// CreateUser handles POST create user request, takes the json data from body,
// creates new user with that data, marshals it to json and sends back in response
func (controller UserController) CreateUser(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	user := models.User{}
	// decode json from post request body into our struct
	json.NewDecoder(req.Body).Decode(&user)
	user.Name = "Fakintosh"

	userJson, err := json.Marshal(user)
	if err != nil {
		log.Println(err)
	}
	// set header
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated) // 201

	log.Println("user have been created:", user.Name)
	fmt.Fprintf(res, "%s\n", userJson)
}
