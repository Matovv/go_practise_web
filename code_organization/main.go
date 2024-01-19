// This is our service handler

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Matovv/go_practise_web/code_organization/constants"
	"github.com/Matovv/go_practise_web/code_organization/handlers"
	"github.com/julienschmidt/httprouter"
)

func main() {
	fmt.Println("\nListening on port", constants.Port)

	router := httprouter.New()
	router.GET("/user/:id", handlers.GetUser)
	router.POST("/user", handlers.CreateUser)
	//router.DELETE("/user/:id", createUserHandler)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", constants.Port), router))

}
