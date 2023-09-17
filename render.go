package main

import (
	"html/template"
	"net/http"
)

var templateFuncs template.FuncMap = template.FuncMap{
	"inc": func(i int) int {
		return i + 1
	},
	"dec": func(i int) int {
		return i - 1
	},
	"multiply": func(x, y float64) float64 {
		return x * y
	},
}

func render(w http.ResponseWriter, r *http.Request, tmpl string, data any, partial ...string) {
	tpl := template.Must(template.New("g").Funcs(templateFuncs).ParseFiles("templates/layout.html", "templates/"+tmpl))

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
