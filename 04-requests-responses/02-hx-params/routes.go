package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("POST /validate", app.validate)
	mux.HandleFunc("POST /login", app.login)
	mux.HandleFunc("GET /authenticated", app.authenticated)
	return mux
}
