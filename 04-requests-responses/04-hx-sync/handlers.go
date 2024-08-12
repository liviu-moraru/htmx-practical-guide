package main

import (
	"github.com/justinas/nosurf"
	"html/template"
	"net/http"
	"strings"
)

const (
	homePageTemplate      = "./ui/html/pages/home.tmpl"
	loginFragmentTemplate = "./ui/html/fragments/login.tmpl"
	authenticatedTemplate = "./ui/html/pages/authenticated.tmpl"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	tmpl, err := getTemplate(homePageTemplate)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	token := nosurf.Token(r)

	err = tmpl.Execute(w, struct{ Token string }{token})
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) validate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		app.clientError(w, r, http.StatusBadRequest, err)
		return
	}

	t := "{{.Error}}"
	tmpl, err := template.New("validate").Parse(t)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	var text string

	email := r.PostForm.Get("email")

	if email != "" {
		if !strings.Contains(email, "@") {
			text = "E-Mail address is invalid."
		}
		err := tmpl.Execute(w, struct{ Error string }{Error: text})
		if err != nil {
			app.serverError(w, r, err)
			return
		}
		return
	}

	password := r.PostForm.Get("password")
	if password != "" {
		if len(strings.TrimSpace(password)) < 8 {
			text = "Password must be at least 8 characters long."
		}
		err := tmpl.Execute(w, struct{ Error string }{Error: text})
		if err != nil {
			app.serverError(w, r, err)
		}
	}
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		app.clientError(w, r, http.StatusBadRequest, err)
		return
	}
	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")
	var errors []string
	if email == "" || !strings.Contains(email, "@") {
		errors = append(errors, "Please enter a valid email address.")
	}
	if password == "" || len(strings.TrimSpace(password)) < 8 {
		errors = append(errors, "Password must be at least 8 characters long.")
	}

	tmpl, err := getTemplate(loginFragmentTemplate)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	err = tmpl.Execute(w, struct{ Errors []string }{Errors: errors})
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) authenticated(w http.ResponseWriter, r *http.Request) {
	tmpl, err := getTemplate(authenticatedTemplate)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func getTemplate(templateFiles ...string) (*template.Template, error) {
	return template.ParseFiles(templateFiles...)
}
