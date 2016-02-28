package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
)

var pathToResources string = "..\\..\\..\\..\\.."

func sayHelloName(w http.ResponseWriter, r *http.Request) {
	fmt.Println("SAYHELLONAME")
	r.ParseForm() //Parse url parameters passed, then parse the response packet for the POST body (request body)
	// attention: If you do not call ParseForm method, the following data can not be obtained form
	fmt.Println(r.Form) // print information on server side.
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Hello Martin!") // write data to response
}

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("LOGIN: method:", r.Method) //get request method
	if r.Method == "GET" {
		t, err := template.ParseFiles(pathToResources + "\\resources\\com\\masterhilli\\WebServerTryout\\Forms\\login.gtpl")
		if (err == nil) {
			t.Execute(w, nil)
		} else {
			fmt.Println("Error on reading resource: " + err.Error());
		}
	} else {
		r.ParseForm()
		// logic part of log in
		fmt.Println("username:", template.HTMLEscapeString(r.Form.Get("username"))) // print at server side
		fmt.Println("password:", template.HTMLEscapeString(r.Form.Get("password")))
		template.HTMLEscape(w, []byte(r.Form.Get("username"))) // responded to clients
	}
}
func main() {
	http.HandleFunc("/say", sayHelloName) // setting router rule
	http.HandleFunc("/login", login)
	err := http.ListenAndServe(":9090", nil) // setting listening port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}