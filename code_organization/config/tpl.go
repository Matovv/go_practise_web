// In this config we initialize our templates as soon as service starts

package config

import "text/template"

var TPL *template.Template

func init() {
	TPL = template.Must(TPL.ParseGlob("/templates/*.gohtml"))
}
