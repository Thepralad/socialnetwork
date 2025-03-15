package render

import (
	"html/template"
	"net/http"
)

func Template(res http.ResponseWriter, file string, data interface{}){
	path := "internal/templates/" + file + ".html"
	tmpl := template.Must(template.ParseFiles(path))
	tmpl.Execute(res, data)
}