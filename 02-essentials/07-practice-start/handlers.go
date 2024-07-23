package main

import (
	"html/template"
	"net/http"
)

const (
	homePageTemplate     = "./ui/html/pages/home.tmpl"
	goalsPartialTemplate = "./ui/html/partials/goals.tmpl"
)

func home(w http.ResponseWriter, r *http.Request) {
	tmpl, err := getTemplate(homePageTemplate, goalsPartialTemplate)
	if err != nil {
		ServerError(w, r, err)
		return
	}

	err = tmpl.Execute(w, courseGoals)
	if err != nil {
		ServerError(w, r, err)
		return
	}

}

func getTemplate(templateFiles ...string) (*template.Template, error) {
	return template.ParseFiles(templateFiles...)
}
