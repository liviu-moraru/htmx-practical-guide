package main

import (
	"errors"
	"html/template"
	"math/rand/v2"
	"net/http"
	"slices"
)

const (
	homePageTemplate                   = "./ui/html/pages/home.tmpl"
	locationPartialTemplate            = "./ui/html/partials/location.tmpl"
	postPlaceFragmentTemplate          = "./ui/html/fragments/postPlace.tmpl"
	deletePlaceFragmentTemplate        = "./ui/html/fragments/deletePlace.tmpl"
	suggestedLocationsFragmentTemplate = "./ui/html/fragments/suggestedLocations.tmpl"
)

type LocationModel struct {
	IsAvailable  bool
	LocationData *Location
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	interestingLocationsModel, notInterestingLocationsModel, suggestedLocationsModel := app.buildLocationsModels()

	locations := struct {
		NotInterestingLocations []LocationModel
		InterestingLocations    []LocationModel
		SuggestedLocations      []LocationModel
	}{
		NotInterestingLocations: notInterestingLocationsModel,
		InterestingLocations:    interestingLocationsModel,
		SuggestedLocations:      suggestedLocationsModel,
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

	_, notInterestingLocationsModel, suggestedLocationsModel := app.buildLocationsModels()

	model := struct {
		Location               LocationModel
		SuggestedLocations     []LocationModel
		NotInterestedLocations []LocationModel
	}{
		Location:               LocationModel{IsAvailable: false, LocationData: location},
		SuggestedLocations:     suggestedLocationsModel,
		NotInterestedLocations: notInterestingLocationsModel,
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

	_, notInterestingLocationsModel, suggestedLocationsModel := app.buildLocationsModels()

	model := struct {
		SuggestedLocations     []LocationModel
		NotInterestedLocations []LocationModel
	}{
		SuggestedLocations:     suggestedLocationsModel,
		NotInterestedLocations: notInterestingLocationsModel,
	}

	tmpl, err := getTemplate(deletePlaceFragmentTemplate, locationPartialTemplate)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = tmpl.Execute(w, model)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) suggestedLocations(w http.ResponseWriter, r *http.Request) {
	suggestedLocations := getSuggestedLocations(availableLocations, interestingLocations)
	suggestedLocationsModel := buildLocationModelSlice(suggestedLocations, true)
	model := struct {
		SuggestedLocations []LocationModel
	}{
		SuggestedLocations: suggestedLocationsModel,
	}

	tmpl, err := getTemplate(suggestedLocationsFragmentTemplate, locationPartialTemplate)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = tmpl.Execute(w, model)
	if err != nil {
		app.serverError(w, r, err)
	}

}

func (app *application) buildLocationsModels() (interestingLocationsModel []LocationModel, notInterestingLocationsModel []LocationModel, suggestedLocationsModel []LocationModel) {
	interestingLocationsModel = buildLocationModelSlice(interestingLocations, false)
	notInterestingLocations := filterLocations(availableLocations, interestingLocations)
	notInterestingLocationsModel = buildLocationModelSlice(notInterestingLocations, true)
	suggestedLocations := getSuggestedLocations(availableLocations, interestingLocations)
	suggestedLocationsModel = buildLocationModelSlice(suggestedLocations, true)
	return interestingLocationsModel, notInterestingLocationsModel, suggestedLocationsModel
}

func getSuggestedLocations(avLoc []Location, intLoc []*Location) []*Location {
	notInterestingLocations := filterLocations(avLoc, intLoc)

	if len(notInterestingLocations) < 2 {
		return notInterestingLocations
	}

	locations := make([]*Location, 0, 2)

	index := rand.IntN(len(notInterestingLocations))
	locations = append(locations, notInterestingLocations[index])

	notInterestingLocations = slices.Delete(notInterestingLocations, index, index+1)

	index2 := rand.IntN(len(notInterestingLocations))
	locations = append(locations, notInterestingLocations[index2])

	return locations
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
