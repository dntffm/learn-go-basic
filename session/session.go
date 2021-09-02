package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "8080"
)

var store *sessions.CookieStore

func init() {
	store = sessions.NewCookieStore([]byte("secret-key"))
}

func home(rw http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")

	var authenticated interface{} = session.Values["authenticated"]
	if authenticated != nil {
		isAuthenticated := session.Values["authenticated"].(bool)

		if !isAuthenticated {
			http.Error(rw, "You are unauthorized to view the page", http.StatusForbidden)
			return
		}

		fmt.Fprintln(rw, "Homepage")
	} else {
		http.Error(rw, "You are unauthorized to view the page", http.StatusForbidden)
		return
	}
}

func login(rw http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")

	session.Values["authenticated"] = true

	session.Save(r, rw)
	fmt.Fprintln(rw, "You have successfully logged in.")
}

func logout(rw http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")

	session.Values["authenticated"] = false

	session.Save(r, rw)
	fmt.Fprintln(rw, "You have successfully logged out.")
}

func main() {
	http.HandleFunc("/home", home)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)

	err := http.ListenAndServe(CONN_HOST+":"+CONN_PORT, nil)

	if err != nil {
		log.Fatal("error starting http server : ", err)
		return
	}

	fmt.Println("Server running...")
}
