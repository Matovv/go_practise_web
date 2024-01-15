package main

import (
	"html/template"
	"io"
	"log"
	"os"
	"strings"
)
var tpl *template.Template

// User type is for testing templates with custom type data
type User struct {
	ID string
	First string
	Last string
	Age int
}


// Initialize templates
func init() {
	tpl = template.Must(template.ParseGlob("go_templates/templates/*.gohtml"))
	log.Println("Init completed!")
}
func main() {
	//ExecuteTemplate("tpl_with_data", "JOSTAAAR")
	//ExecuteTemplate("tpl_complex", []int{1,5,25,100,500})
	/*
	m := map[string]string {
		"23242": "Alfred",
		"00032": "Jhon",
		"09981": "Tiger",
	}
	ExecuteTemplate("tpl_complex", m)
	*/
	u1 := User{"00039","Jhon","Wicked",32}
	u2 := User{"00031","Lily","Deepthroat",19}
	u3 := User{"00107","Phantom","Lancer",26}
	//ExecuteTemplate("tpl_structs", u1)
	ExecuteTemplate("tpl_complexStructs", []User{u1,u2,u3})
}

// ExecuteTemplate executes given template name and passes given data, while handling all error that may occur.
func ExecuteTemplate(tplName string, data any) {
	err := tpl.ExecuteTemplate(os.Stdout, tplName + ".gohtml", data)
	if err != nil {
		log.Fatalln(err)
	}
}

// CreateTemplate creates template based on specifications inside of function
func CreateTemplate() {
	name := "Jhon Wicked"

	tplObj := `
	<!DOCTYPE html>
	<html lang="en">
	<head>
	<meta charset="UTF-8">
	<title>Hallow World!</title>
	</head>
	<body>
	<h1>Welcome, ` + name + `!</h1>
	</body>
	</html> 
	`

	nf, err := os.Create("go_templates/templates/index.gohtml")
	if err != nil {
		log.Fatal("error creating file", err)
	}
	defer nf.Close()

	io.Copy(nf, strings.NewReader(tplObj))
}
