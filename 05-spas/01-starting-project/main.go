package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

const (
	jsonFilePath = "./data/products.json"
)

type Product struct {
	ID          string  `json:"id"`
	Image       string  `json:"image"`
	Title       string  `json:"title"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
}

var products []Product

type application struct{}

var app *application

func main() {
	var err error
	products, err = getDataFromFile(jsonFilePath)
	if err != nil {
		log.Fatal(err)
	}
	app = &application{}
	mux := app.routes()
	log.Print("starting server on :3000")
	err = http.ListenAndServe(":3000", mux)
	log.Fatal(err)
}

func getDataFromFile(filePath string) ([]Product, error) {
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

func parseData(byteData []byte) ([]Product, error) {
	var data []Product
	err := json.Unmarshal(byteData, &data)
	return data, err
}
