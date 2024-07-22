package main

import (
	"html/template"
	"net/http"
)

const (
	homePageTemplate    = "./ui/html/pages/home.tmpl"
	infoPartialTemplate = "./ui/html/partials/info.tmpl"
)

func home(w http.ResponseWriter, r *http.Request) {
	tmpl, err := getTemplate(homePageTemplate, infoPartialTemplate)
	if err != nil {
		ServerError(w, r, err)
		return
	}

	err = tmpl.Execute(w, data.Knowledge)
	if err != nil {
		ServerError(w, r, err)
		return
	}

}

func note(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		ClientError(w, r, http.StatusBadRequest, err)
	}

	note := r.PostFormValue("note")
	if note != "" {
		data.Knowledge = append([]string{note}, data.Knowledge...)
	}

	http.Redirect(w, r, "/", http.StatusFound)

	//tmpl, err := getTemplate(infoPartialTemplate)
	//if err != nil {
	//	ServerError(w, r, err)
	//	return
	//}
	//
	//err = tmpl.ExecuteTemplate(w, "info", data.Knowledge)
	//if err != nil {
	//	ServerError(w, r, err)
	//	return
	//}

}

func getTemplate(templateFiles ...string) (*template.Template, error) {
	return template.ParseFiles(templateFiles...)
}
