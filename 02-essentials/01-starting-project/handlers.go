package main

import (
	"html/template"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	Render(w, r, "./ui/html/pages/home.tmpl")
}

func Render(w http.ResponseWriter, r *http.Request, page string) {
	ts, err := template.ParseFiles(page)
	if err != nil {
		ServerError(w, r, err)
		return
	}
	err = ts.Execute(w, nil)
	if err != nil {
		ServerError(w, r, err)
	}
}
