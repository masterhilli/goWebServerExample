package FormUser

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	"crypto/md5"
	"io"
	"strconv"
	"io/ioutil"
)

const userNameKey string = "username"
const passwordKey string = "password"

var pathToResources string = "..\\..\\..\\..\\.."
var logger *log.Logger = log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)

func sayHelloName(w http.ResponseWriter, r *http.Request) {
	logger.Println("SAYHELLONAME: method:" +  r.Method)
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
	logger.Println("LOGIN: method:" +  r.Method) //get request method
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

var tokens map[string]bool = make(map[string]bool)
func tokenizerLogin(w http.ResponseWriter, r *http.Request) {
	logger.Println("LOGIN2 method:", r.Method) // get request method
	r.ParseForm()
	logger.Println(r.RequestURI)
	if r.Method == "GET" {
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))
		t, _ := template.ParseFiles(pathToResources + "\\resources\\com\\masterhilli\\WebServerTryout\\Forms\\login2.gtpl")
		t.Execute(w, token)
		logger.Println("again we are in the getter, but why?")
	} else {
		// log in request
		r.ParseForm()
		token := r.Form.Get("token")
		if token == "" {
			retVal := "No token received, invalid call to login method!"
			fmt.Fprintln(w, retVal)
			logger.Println(retVal)
			return
		}else if (tokens[token] ) {
			logger.Println("Token: " + token +" duplicate request. Ignore the request, but print the same")
		}
		tokens[token] = true
		logger.Println("username length:", len(r.Form["username"][0]))
		logger.Println("username:", template.HTMLEscapeString(r.Form.Get("username"))) // print in server side
		logger.Println("password:", template.HTMLEscapeString(r.Form.Get("password")))
		template.HTMLEscape(w, []byte(r.Form.Get("username"))) // respond to client
	}
}

func tryEscapeSequences(w http.ResponseWriter, z *http.Request) {
	t, _ := template.New("foo").Parse(
		`{{define "T"}}Hello, {{.}}!{{end}}<br/>
		{{define "H"}} Escaped HTML: {{.}}!{{end}}<br/>
		{{define "N"}} Not escaped HTML: {{.}}!{{end}}`)
	_ = t.ExecuteTemplate(w, "T", template.HTML("<script>alert('you have been pwned')</script>"))
	t.ExecuteTemplate(w, "H", template.HTMLEscapeString("<script>alert('you have been pwned')</script>"))
	t.ExecuteTemplate(w, "N", "<script>alert('you have been pwned')</script>") // okay I do not get why this is escaped???

}

// upload logic
func upload(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method)
	if r.Method == "GET" {
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))
		t, _ := template.ParseFiles(pathToResources + "\\resources\\com\\masterhilli\\WebServerTryout\\Forms\\upload.gtpl")
		t.Execute(w, token)
	} else {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		fmt.Fprintf(w, "%v", handler.Header)
		f, err := os.OpenFile("./.tmp/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
	}
}

func bootStrapGet(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method)
	if r.Method == "GET" {
		r.ParseForm()
		bytes, err := ioutil.ReadFile("."+ r.RequestURI)
		if err != nil {
			logger.Println(err.Error())
		}
		w.Header().Set("Content-Type", "text/css")
		w.Write(bytes)
		logger.Println(http.DetectContentType(bytes))

	}
}

func validateField(key, value string, validator Validator) bool {
	return validator.Validate(key, value)
}

func RunWebServer(resourceRootFolderPath string) {
	pathToResources = resourceRootFolderPath
	http.HandleFunc("/escape", tryEscapeSequences)
	http.HandleFunc("/upload", upload)
	http.HandleFunc("/say", sayHelloName) // setting router rule
	http.HandleFunc("/login", login)
	http.HandleFunc("/login2", tokenizerLogin)
	http.HandleFunc("/bootstrap.min.css", bootStrapGet)
	err := http.ListenAndServe(":9090", nil) // setting listening port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
