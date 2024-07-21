package main

import (
	"log"
	"net/http"
	"runtime/debug"
)

func ServerError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
		trace  = string(debug.Stack())
	)
	log.Println(method, uri, trace, err.Error())
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
