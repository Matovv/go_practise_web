// handlers package provide handler methods for our service
package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Matovv/go_practise_web/code_organization/models"
	"github.com/julienschmidt/httprouter"
)

// CreateUser handles POST create user request, takes the json data from body,
// creates new user with that data, marshals it to json and sends back in response
func CreateUser(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
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
