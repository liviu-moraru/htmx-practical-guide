package main

import (
	"log"
	"net/http"
)

type application struct{}

var app *application

func main() {
	app = &application{}
	mux := app.routes()
	log.Print("starting server on :3000")
	err := http.ListenAndServe(":3000", mux)
	log.Fatal(err)
}
