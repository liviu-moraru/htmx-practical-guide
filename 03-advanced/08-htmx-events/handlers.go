package main

import (
	"errors"
	"html/template"
	"net/http"
	"slices"
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

func deleteGoal(w http.ResponseWriter, r *http.Request) {
	goalID := r.PathValue("goalID")
	if goalID == "" {
		ClientError(w, r, http.StatusBadRequest, errors.New("goalID is required"))
		return
	}

	courseGoals = slices.DeleteFunc(courseGoals, func(g Goal) bool {
		return g.ID == goalID
	})
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

	newGoal := Goal{ID: getRandomString(10), Text: goal}
	courseGoals = append(courseGoals, newGoal)

	tmpl, err := getTemplate(goalsPartialTemplate)
	if err != nil {
		ServerError(w, r, err)
		return
	}

	err = tmpl.ExecuteTemplate(w, "goal", newGoal)
	if err != nil {
		ServerError(w, r, err)
		return
	}
}

func getTemplate(templateFiles ...string) (*template.Template, error) {
	return template.ParseFiles(templateFiles...)
}
