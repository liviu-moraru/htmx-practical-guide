package main

import (
	"errors"
	"html/template"
	"net/http"
	"slices"
)

const (
	homePageTemplate        = "./ui/html/pages/home.tmpl"
	locationPartialTemplate = "./ui/html/partials/location.tmpl"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	notInterestingLocations := filterLocations(availableLocations, interestingLocations)

	locations := struct {
		NotInterestingLocations []*Location
		InterestingLocations    []*Location
	}{
		NotInterestingLocations: notInterestingLocations,
		InterestingLocations:    interestingLocations,
	}

	tmpl, err := getTemplate(homePageTemplate, locationPartialTemplate)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = tmpl.Execute(w, locations)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

}

func (app *application) places(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		app.clientError(w, r, http.StatusBadRequest, err)
		return
	}
	locationID := r.PostForm.Get("locationId")

	location, err := getLocationByID(locationID)
	if err != nil {
		app.clientError(w, r, http.StatusBadRequest, err)
		return
	}
	interestingLocations = append(interestingLocations, location)
	//TODO
	tmpl, err := getTemplate(locationPartialTemplate)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	err = tmpl.ExecuteTemplate(w, "location", location)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) delete(w http.ResponseWriter, r *http.Request) {
	locationID := r.PathValue("id")
	if locationID == "" {
		app.clientError(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}
	interestingLocations = slices.DeleteFunc(interestingLocations, func(loc *Location) bool {
		return loc.ID == locationID
	})
}

func filterLocations(availableLocations []Location, interestingLocations []*Location) []*Location {
	var filteredLocations []*Location
	for _, loc := range availableLocations {
		if !slices.ContainsFunc(interestingLocations, func(iloc *Location) bool {
			return iloc.ID == loc.ID
		}) {
			filteredLocations = append(filteredLocations, &loc)
		}
	}
	return filteredLocations
}

func getTemplate(templateFiles ...string) (*template.Template, error) {
	return template.ParseFiles(templateFiles...)
}

func getLocationByID(locationID string) (*Location, error) {
	index := slices.IndexFunc(availableLocations, func(l Location) bool {
		return l.ID == locationID
	})
	if index == -1 {
		return nil, errors.New("location not found")
	}
	return &availableLocations[index], nil
}
