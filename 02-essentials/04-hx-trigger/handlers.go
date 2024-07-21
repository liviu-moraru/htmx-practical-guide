package main

import (
	"encoding/json"
	"html/template"
	"io"
	"net/http"
	"os"
)

func home(w http.ResponseWriter, r *http.Request) {
	Render(w, r, "./ui/html/pages/home.tmpl")
}

type Data struct {
	Knowledge []string `json:"htmx_knowledge"`
}

func info(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("./data/htmx-info.json")
	if err != nil {
		ServerError(w, r, err)
		return
	}
	defer file.Close()
	byteData, err := io.ReadAll(file)
	if err != nil {
		ServerError(w, r, err)
		return
	}

	var data Data
	err = json.Unmarshal(byteData, &data)
	if err != nil {
		ServerError(w, r, err)
		return
	}
	tmpl, err := template.ParseFiles(
		"./ui/html/fragments/info.tmpl")
	if err != nil {
		ServerError(w, r, err)
		return
	}

	err = tmpl.Execute(w, data.Knowledge)
	if err != nil {
		ServerError(w, r, err)
	}
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
