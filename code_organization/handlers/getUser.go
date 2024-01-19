package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Matovv/go_practise_web/code_organization/config"
	"github.com/Matovv/go_practise_web/code_organization/models"
	"github.com/julienschmidt/httprouter"
)

// GetUser handles GET user request, and returns user by the id param
// if url queries ('name' in this case) specified then does additional functions
func GetUser(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	user := models.User{}
	rows, err := config.DB.Query(`SELECT col1, col2 from users`)
	if err != nil {
		panic(err)
	}
	var s, col1, col2 string
	s = "Retrieved Records:\n"

	for rows.Next() {
		err = rows.Scan(&col1, &col2)
		if err != nil {
			log.Println(err)
		}
		s += col1 + col2 + "\n"
	}
	user.Name=col1
	fmt.Fprintln(res, s)

}
