package main

import (
	"errors"
	"html/template"
	"net/http"
	"slices"
)

const (
	homePageTemplate          = "./ui/html/pages/home.tmpl"
	locationPartialTemplate   = "./ui/html/partials/location.tmpl"
	postPlaceFragmentTemplate = "./ui/html/fragments/postPlace.tmpl"
)

type LocationModel struct {
	IsAvailable  bool
	LocationData *Location
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	interestingLocationsModel, notInterestingLocationsModel := app.buildLocationsModels()

	locations := struct {
		NotInterestingLocations []LocationModel
		InterestingLocations    []LocationModel
	}{
		NotInterestingLocations: notInterestingLocationsModel,
		InterestingLocations:    interestingLocationsModel,
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

	_, notInterestingLocationsModel := app.buildLocationsModels()

	model := struct {
		Location  LocationModel
		Locations []LocationModel
	}{
		Location:  LocationModel{IsAvailable: false, LocationData: location},
		Locations: notInterestingLocationsModel,
	}

	tmpl, err := getTemplate(postPlaceFragmentTemplate, locationPartialTemplate)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = tmpl.Execute(w, model)
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

func (app *application) buildLocationsModels() (interestingLocationsModel []LocationModel, notInterestingLocationsModel []LocationModel) {
	interestingLocationsModel = buildLocationModelSlice(interestingLocations, false)
	notInterestingLocations := filterLocations(availableLocations, interestingLocations)
	notInterestingLocationsModel = buildLocationModelSlice(notInterestingLocations, true)
	return interestingLocationsModel, notInterestingLocationsModel
}

func buildLocationModelSlice(locations []*Location, isAvailable bool) []LocationModel {
	locationModels := make([]LocationModel, len(locations))
	for i := range locations {
		locationModels[i] = LocationModel{
			IsAvailable:  isAvailable,
			LocationData: locations[i],
		}
	}
	return locationModels
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
