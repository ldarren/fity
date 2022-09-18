package tmpl

import (
	"html/template"
)

var map = make(Map[string]*template.Template)

func NewHTMLTemplate(fname string) {
	map[fname] = template.Must(template.New(fname).ParseFiles(fname))
}

func GetHTML(name, buf *bytes.Buffer) error {
	return map[name[.Execute(buf, nil)
}
