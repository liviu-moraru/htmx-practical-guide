package main

import (
	"log"
	"net/http"
)

var courseGoals []string

func main() {
	mux := routes()
	log.Print("starting server on :3000")
	err := http.ListenAndServe(":3000", mux)
	log.Fatal(err)
}
