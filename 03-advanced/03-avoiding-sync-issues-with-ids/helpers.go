package main

import (
	"log"
	"math/rand/v2"
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

func ClientError(w http.ResponseWriter, r *http.Request, status int, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
		trace  = string(debug.Stack())
	)
	log.Println(method, uri, trace, err.Error())
	http.Error(w, http.StatusText(status), status)
}

var alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

func getRandomString(length int) string {
	result := make([]byte, length)
	for i := range result {
		result[i] = alphabet[rand.IntN(len(alphabet))]
	}

	return string(result)

}
