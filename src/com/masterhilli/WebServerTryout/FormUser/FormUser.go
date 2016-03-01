package FormUser

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

const userNameKey string = "username"
const passwordKey string = "password"

var pathToResources string = "..\\..\\..\\..\\.."
var logger *log.Logger = log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)

func sayHelloName(w http.ResponseWriter, r *http.Request) {
	logger.Println("SAYHELLONAME")
	r.ParseForm() //Parse url parameters passed, then parse the response packet for the POST body (request body)
	// attention: If you do not call ParseForm method, the following data can not be obtained form
	logger.Println(r.Form) // print information on server side.
	logger.Println("path", r.URL.Path)
	logger.Println("scheme", r.URL.Scheme)
	logger.Println(r.Form["url_long"])
	for k, v := range r.Form {
		logger.Println("key:", k)
		logger.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Hello Martin!") // write data to response
}

func login(w http.ResponseWriter, r *http.Request) {
	logger.Println("LOGIN: method:", r.Method) //get request method
	if r.Method == "GET" {
		// TODO: add fields to test validator or unit tests
		t, err := template.ParseFiles(pathToResources + "\\resources\\com\\masterhilli\\WebServerTryout\\Forms\\login.gtpl")
		if err == nil {
			t.Execute(w, nil)
		} else {
			logger.Println("Error on reading resource: " + err.Error())
		}
	} else {
		r.ParseForm()
		// logic part of log in
		username := template.HTMLEscapeString(r.Form.Get(userNameKey))
		password := template.HTMLEscapeString(r.Form.Get(passwordKey))

		validateField(userNameKey, username, RequiredFieldValidator{})
		validateField(passwordKey, password, NumberFieldValidator{})

		logger.Printf("Usr: \"%s\" / Pwd: \"%s\"", username, password)
		template.HTMLEscape(w, []byte(r.Form.Get(userNameKey))) // responded to clients
	}
}

func validateField(key, value string, validator Validator) bool {
	return validator.Validate(key, value)
}

func RunWebServer(resourceRootFolderPath string) {
	pathToResources = resourceRootFolderPath
	http.HandleFunc("/say", sayHelloName) // setting router rule
	http.HandleFunc("/login", login)
	err := http.ListenAndServe(":9090", nil) // setting listening port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
