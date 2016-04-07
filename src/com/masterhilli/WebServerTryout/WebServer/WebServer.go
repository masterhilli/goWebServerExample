package Webserver

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
	. "./DB"
	. "./Logger"
)

const (
	userNameKey string = "Username"
 	passwordKey string = "Password"
	loginRelativeResourcePath string = "\\resources\\com\\masterhilli\\WebServerTryout\\Forms\\"
)

// TODO: make a singleton out of the server
var (
	pathToResources string = "..\\..\\..\\..\\.."
	tokens map[string]bool = make(map[string]bool)

	dbConnection DBConnector = DBConnector{}
)

func login(w http.ResponseWriter, r *http.Request) {
	LOGGER.Println("LOGIN: method:" +  r.Method) //get request method
	if r.Method == "GET" {
		// TODO: add fields to test validator or unit tests
		t, err := template.ParseFiles(pathToResources + loginRelativeResourcePath +"login.gtpl")
		if err == nil {
			t.Execute(w, nil)
		} else {
			LOGGER.Println("Error on reading resource: " + err.Error())
		}
	} else {
		r.ParseForm()
		// logic part of log in
		username := r.Form.Get(userNameKey)
		password := r.Form.Get(passwordKey)
		pwdFromDB := dbConnection.SelectTable(TABLE_NAME_USERS).ReceiveStringWhere(COL_USER_NAME, username)
		if strings.Compare(pwdFromDB, password) == 0 {
			//TODO: return something Successfull ;)
			fmt.Fprintf(w, "You lucky fuck, you got in!")
		} else {
			fmt.Fprintf(w, "Me scuzzi, you are nota gona geta in!")
		}
		/*
		template.HTMLEscapeString(r.Form.Get(userNameKey))
		password := template.HTMLEscapeString()

		validateField(userNameKey, username, RequiredFieldValidator{})
		validateField(passwordKey, password, NumberFieldValidator{})

		logger.Printf("Usr: \"%s\" / Pwd: \"%s\"", username, password)
		template.HTMLEscape(w, []byte(r.Form.Get(userNameKey))) // responded to clients */
	}
}


func bootStrapGet(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method)
	if r.Method == "GET" {
		r.ParseForm()
		bytes, err := ioutil.ReadFile("."+ r.RequestURI)
		if err != nil {
			LOGGER.Println(err.Error())
		}
		w.Header().Set("Content-Type", "text/css")
		w.Write(bytes)
		LOGGER.Println(http.DetectContentType(bytes))

	}
}

func RunWebServer(resourceRootFolderPath string) {

	dbConnection.Initialize("MartinsWebServer")
	defer dbConnection.Close()
	pathToResources = resourceRootFolderPath
	//http.HandleFunc("/escape", tryEscapeSequences)
	//http.HandleFunc("/upload", upload)
	//http.HandleFunc("/say", sayHelloName) // setting router rule
	http.HandleFunc("/login", login)
	//http.HandleFunc("/login2", tokenizerLogin)
	http.HandleFunc("/bootstrap.min.css", bootStrapGet)
	err := http.ListenAndServe(":9090", nil) // setting listening port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}


//**************************************************************************************
//********************* TRY OUT CODE ***************************************************

func sayHelloName(w http.ResponseWriter, r *http.Request) {
	LOGGER.Println("SAYHELLONAME: method:" +  r.Method)
	r.ParseForm() //Parse url parameters passed, then parse the response packet for the POST body (request body)
	// attention: If you do not call ParseForm method, the following data can not be obtained form
	LOGGER.Println(r.Form) // print information on server side.
	LOGGER.Println("path", r.URL.Path)
	LOGGER.Println("scheme", r.URL.Scheme)
	LOGGER.Println(r.Form["url_long"])
	for k, v := range r.Form {
		LOGGER.Println("key:", k)
		LOGGER.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Hello Martin!") // write data to response
}


func validateField(key, value string, validator Validator) bool {
	return validator.Validate(key, value)
}


func tokenizerLogin(w http.ResponseWriter, r *http.Request) {
	LOGGER.Println("LOGIN2 method:", r.Method) // get request method
	r.ParseForm()
	LOGGER.Println(r.RequestURI)
	if r.Method == "GET" {
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))
		t, _ := template.ParseFiles(pathToResources + "\\resources\\com\\masterhilli\\WebServerTryout\\Forms\\login2.gtpl")
		t.Execute(w, token)
		LOGGER.Println("again we are in the getter, but why?")
	} else {
		// log in request
		r.ParseForm()
		token := r.Form.Get("token")
		if token == "" {
			retVal := "No token received, invalid call to login method!"
			fmt.Fprintln(w, retVal)
			LOGGER.Println(retVal)
			return
		}else if (tokens[token] ) {
			LOGGER.Println("Token: " + token +" duplicate request. Ignore the request, but print the same")
		}
		tokens[token] = true
		LOGGER.Println("username length:", len(r.Form["username"][0]))
		LOGGER.Println("username:", template.HTMLEscapeString(r.Form.Get("username"))) // print in server side
		LOGGER.Println("password:", template.HTMLEscapeString(r.Form.Get("password")))
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