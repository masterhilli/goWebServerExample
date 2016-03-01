package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func sayHelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	/*fmt.Println(r.Form)
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])*/
	var name string = ""
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
		if k == "name" {
			name = strings.Join(v, "")
		}
	}
	if len(name) == 0 {
		fmt.Fprintf(w, "Hello Martin!") // send data to client side
	} else {
		fmt.Fprintf(w, fmt.Sprintf("Hello %s!", name))
	}
}

type ServerHandlerMutex struct{}

func (this ServerHandlerMutex) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if (r.URL.Path == "/say") || (r.URL.Path == "/sayMyName/") {
		sayHelloName(w, r)
		return
	}
	http.NotFound(w, r)
}

func main() {
	serverMux := &ServerHandlerMutex{}
	http.HandleFunc("/sayHelloName", sayHelloName) // set router
	err := http.ListenAndServe(":9090", serverMux) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
