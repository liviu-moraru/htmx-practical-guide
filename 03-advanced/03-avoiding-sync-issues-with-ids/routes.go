package main

import "net/http"

func routes() *http.ServeMux {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("POST /goals", goals)
	mux.HandleFunc("DELETE /goals/{goalID}", deleteGoal)
	return mux
}
