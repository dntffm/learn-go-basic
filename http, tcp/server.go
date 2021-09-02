package main

import (
	"crypto/subtle"
	"fmt"
	"log"
	"net/http"
)

const (
	CONN_HOST  = "localhost"
	CONN_PORT  = "1234"
	ADMIN_USER = "admin"
	ADMIN_PW   = "admin"
)

func helloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world")
}

func auth(handler http.HandlerFunc, realm string) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()

		if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(ADMIN_USER)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(ADMIN_PW)) != 1 {
			rw.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
			rw.WriteHeader(401)
			rw.Write([]byte("You are not authorized. \n"))
			return
		}
		handler(rw, r)
	}
}

func main() {
	http.HandleFunc("/login", auth(helloWorld, "Please enter your credentials"))
	err := http.ListenAndServe(CONN_HOST+":"+CONN_PORT, nil)

	if err != nil {
		log.Fatal("error starting http server : ", err)
		return

	}
}
