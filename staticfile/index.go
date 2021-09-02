package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/schema"
)

const (
	CONN_HOST              = "localhost"
	CONN_PORT              = "8000"
	USERNAME_ERROR_MESSAGE = "Please enter a valid Username"
	PASSWORD_ERROR_MESSAGE = "Please enter a valid Password"
	GENERIC_ERROR_MESSAGE  = "Validation Error"
)

type Person struct {
	Id   string
	Name string
}

type User struct {
	Username string `valid : "alpha, required"`
	Password string `valid : "alpha, required"`
}

func validateFormLogin(rw http.ResponseWriter, r *http.Request, user *User) (bool, string) {
	valid, validationError := govalidator.ValidateStruct(user)

	if !valid {
		usernameError := govalidator.ErrorByField(validationError, "Username")
		passwordError := govalidator.ErrorByField(validationError, "Password")

		if usernameError != "" {
			log.Printf("username validation error : ", usernameError)
			return valid, USERNAME_ERROR_MESSAGE
		}

		if passwordError != "" {
			log.Printf("password validation error : ", passwordError)
			return valid, PASSWORD_ERROR_MESSAGE
		}

	}
	return valid, GENERIC_ERROR_MESSAGE
}

func readForm(r *http.Request) *User {
	r.ParseForm()
	user := new(User)
	decoder := schema.NewDecoder()
	decodeErr := decoder.Decode(user, r.PostForm)

	if decodeErr != nil {
		log.Printf("error mapping parsed form data to struct : ", decodeErr)
	}

	return user
}

func login(rw http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		parsedTemplate, _ := template.ParseFiles("templates/form.html")
		parsedTemplate.Execute(rw, nil)
	} else {
		user := readForm(r)
		valid, validationErrorMessage := validateFormLogin(rw, r, user)

		if !valid {
			fmt.Fprintf(rw, validationErrorMessage)
			return
		}
		fmt.Fprintf(rw, "Hello, "+user.Username)
	}

}

func upload(rw http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		parsedTemplate, _ := template.ParseFiles("templates/upload.html")
		err := parsedTemplate.Execute(rw, nil)

		if err != nil {
			log.Printf("Error occurred while executing the template or writing its output : ", err)
			return
		}
	} else {
		file, header, err := r.FormFile("file")

		if err != nil {
			log.Printf("error getting a file for the provided form key : ", err)
			return
		}

		defer file.Close()

		out, pathError := os.Create("tmp/uploadedFile")
		if pathError != nil {
			log.Printf("error creating a file for writing : ", pathError)
			return
		}

		defer out.Close()

		_, copyFileError := io.Copy(out, file)

		if copyFileError != nil {
			log.Printf("error occurred while file copy : ", copyFileError)
		}

		fmt.Fprintf(rw, "File uploaded successfully : "+header.Filename)
	}
}

func renderTemplate(rw http.ResponseWriter, r *http.Request) {
	person := Person{"11", "Philip"}

	parsedTemplate, _ := template.ParseFiles("templates/template.html")

	err := parsedTemplate.Execute(rw, person)

	if err != nil {
		log.Printf("Error occurred while executing the template or writing its output : ", err)
		return
	}

}

func main() {
	fileServer := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))
	http.HandleFunc("/", renderTemplate)
	http.HandleFunc("/login", login)
	http.HandleFunc("/upload", upload)
	err := http.ListenAndServe(CONN_HOST+":"+CONN_PORT, nil)

	if err != nil {
		log.Fatal("error starting http server : ", err)
		return
	}
}
