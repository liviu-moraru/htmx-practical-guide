package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

const (
	infoJSONFilePath = "./data/htmx-info.json"
)

type Data struct {
	Knowledge []string `json:"htmx_knowledge"`
}

var data *Data

func main() {
	var err error
	data, err = getDataFromFile(infoJSONFilePath)
	if err != nil {
		log.Fatal(err)
	}

	mux := routes()
	log.Print("starting server on :3000")
	err = http.ListenAndServe(":3000", mux)
	log.Fatal(err)
}

func getDataFromFile(filePath string) (*Data, error) {
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

func parseData(byteData []byte) (*Data, error) {
	var data Data
	err := json.Unmarshal(byteData, &data)
	return &data, err
}
