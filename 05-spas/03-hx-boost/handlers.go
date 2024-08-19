package main

import (
	"errors"
	"html/template"
	"net/http"
	"slices"
)

const (
	homePageTemplate    = "./ui/html/pages/home.tmpl"
	productPageTemplate = "./ui/html/pages/product.tmpl"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	tmpl, err := getTemplate(homePageTemplate)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	if err = tmpl.Execute(w, struct{ Products []Product }{products}); err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) product(w http.ResponseWriter, r *http.Request) {
	productID := r.PathValue("productID")
	if productID == "" {
		app.clientError(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}
	tmpl, err := getTemplate(productPageTemplate)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	index := slices.IndexFunc(products, func(p Product) bool {
		return p.ID == productID
	})

	if index < 0 {
		app.clientError(w, r, http.StatusBadRequest, errors.New("no such product"))
	}

	if err = tmpl.Execute(w, products[index]); err != nil {
		app.serverError(w, r, err)
	}

}

func (app *application) cart(w http.ResponseWriter, r *http.Request) {
	//...
	http.Redirect(w, r, "/", http.StatusFound)
}

func getTemplate(templateFiles ...string) (*template.Template, error) {
	return template.ParseFiles(templateFiles...)
}
