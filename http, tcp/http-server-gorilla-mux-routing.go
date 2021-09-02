package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

const
(
	CONN_HOST = "localhost"
	CONN_PORT = "8080"
)

var getRequestHandler = http.HandlerFunc(
	func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("It's a GET request"))
	},
)

var postRequestHandler = http.HandlerFunc(
	func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("It's a POST request"))
	},
)

var pathVariableHandler = http.HandlerFunc(
	func(rw http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
		rw.Write([]byte("Hi "+name))
	},
)

func main()  {
	router := mux.NewRouter()
	router.Handle("/", getRequestHandler).Methods("GET")
	router.Handle("/post", postRequestHandler).Methods("POST")
	router.Handle("/hello/{name}", pathVariableHandler).Methods("GET", "PUT")

	http.ListenAndServe(CONN_HOST+":"+CONN_PORT, router)
}