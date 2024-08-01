package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

const (
	jsonFilePath = "./data/available-locations.json"
)

type Image struct {
	Src string `json:"src"`
	Alt string `json:"alt"`
}

type Location struct {
	ID    string  `json:"id"`
	Title string  `json:"title"`
	Image Image   `json:"image"`
	Lat   float64 `json:"lat"`
	Long  float64 `json:"long"`
}

type application struct{}

var availableLocations []Location
var interestingLocations []*Location

var app *application

func main() {
	app = &application{}
	var err error
	availableLocations, err = getDataFromFile(jsonFilePath)
	if err != nil {
		log.Fatal(err)
	}
	mux := app.routes()
	log.Print("starting server on :3000")
	err = http.ListenAndServe(":3000", mux)
	log.Fatal(err)
}

func getDataFromFile(filePath string) ([]Location, error) {
	byteData, err := readFileData(filePath)
	if err != nil {
		return nil, err
	}
	return parseData(byteData)
}

func readFileData(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return io.ReadAll(file)
}

func parseData(byteData []byte) ([]Location, error) {
	var data []Location
	err := json.Unmarshal(byteData, &data)
	return data, err
}
