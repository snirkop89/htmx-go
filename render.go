package main

import (
	"html/template"
	"net/http"
)

func render(w http.ResponseWriter, r *http.Request, tmpl string, data any, partial ...string) {
	tpl := template.Must(template.ParseFiles("templates/layout.html", "templates/"+tmpl))

	if len(partial) > 0 {
		if err := tpl.ExecuteTemplate(w, partial[0], data); err != nil {
			panic(err)
		}
	} else {
		if err := tpl.ExecuteTemplate(w, "base", data); err != nil {
			panic(err)
		}
	}
}
