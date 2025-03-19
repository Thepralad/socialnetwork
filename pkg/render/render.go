package render

import (
	"html/template"
	"net/http"
	"path/filepath"
)

var templates *template.Template

func Init() error {
	var err error
	templates, err = template.ParseGlob(filepath.Join("internal", "templates", "*.html"))
	return err
}

func Template(res http.ResponseWriter, file string, data interface{}) error {
	if templates == nil {
		if err := Init(); err != nil {
			return err
		}
	}
	return templates.ExecuteTemplate(res, file+".html", data)
}