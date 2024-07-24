package main

import (
	"errors"
	"html/template"
	"net/http"
)

const (
	homePageTemplate     = "./ui/html/pages/home.tmpl"
	goalsPartialTemplate = "./ui/html/partials/goals.tmpl"
	goalFragmentTemplate = "./ui/html/fragments/goal.tmpl"
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

type TemplateData struct {
	Index int
	Item  string
}

func goals(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		ClientError(w, r, http.StatusBadRequest, err)
		return
	}
	goal := r.Form.Get("goal")
	if goal == "" {
		ClientError(w, r, http.StatusBadRequest, errors.New("missing goal parameter"))
		return
	}
	courseGoals = append(courseGoals, goal)

	tmpl, err := getTemplate(goalFragmentTemplate)
	if err != nil {
		ServerError(w, r, err)
		return
	}

	err = tmpl.Execute(w, TemplateData{Index: len(courseGoals) - 1, Item: goal})
	if err != nil {
		ServerError(w, r, err)
		return
	}
}

func getTemplate(templateFiles ...string) (*template.Template, error) {
	return template.ParseFiles(templateFiles...)
}
