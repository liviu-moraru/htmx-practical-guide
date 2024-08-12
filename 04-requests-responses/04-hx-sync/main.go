package main

import (
	"github.com/justinas/nosurf"
	"log"
	"net/http"
)

type application struct{}

var app *application

func main() {
	app = &application{}
	mux := app.routes()
	log.Print("starting server on :3000")
	err := http.ListenAndServe(":3000", noSurf(mux))
	log.Fatal(err)
}

func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})
	return csrfHandler
}
