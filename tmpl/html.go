package tmpl

import (
	"errors"
	"html/template"
	"io"
	"log"
)

var tmpls *template.Template

func Load(path string) error {
	tmpls = template.Must(template.ParseGlob(path + "*.html"))
	if tmpls == nil {
		log.Fatal("Template not found at: " + path)
		return errors.New("Template not found at" + path)	
	}
	return nil
}

func RenderHTML(w io.Writer, name string, data any) error {
	return  tmpls.ExecuteTemplate(w, name + ".html", data)
}
