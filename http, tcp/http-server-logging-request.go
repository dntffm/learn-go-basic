package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

const (
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
		rw.Write([]byte("Hi " + name))
	},
)

func main() {
	router := mux.NewRouter()
	logfile, err := os.OpenFile("server.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		log.Fatal("error starting http server : ", err)
		return
	}
	router.Handle("/", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(getRequestHandler))).Methods("GET")
	router.Handle("/post", handlers.LoggingHandler(logfile, http.HandlerFunc(postRequestHandler))).Methods("POST")
	router.Handle("/hello/{name}", handlers.LoggingHandler(logfile, http.HandlerFunc(pathVariableHandler))).Methods("GET", "PUT")

	http.ListenAndServe(CONN_HOST+":"+CONN_PORT, router)
}
