package main

import (
	"crypto/subtle"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

var (
	user     = "root"
	password = "lein23"
	connHost = "localhost"
	connPort = "8081"

	inputFile *string

	routes = Routes{

		// Route{
		// 	"editserver",
		// 	"POST",
		// 	"/editserver",
		// 	auth(homePage, enterYourUserNamePassword, accesses),
		// },

		// Route{
		// 	"getPersons",
		// 	"GET",
		// 	"/persons",
		// 	auth(getPersons, enterYourUserNamePassword),
		// },
	}
)

func auth(handler HandlerFunc2, realm string, accesses *[]Access) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(user)) != 1 ||
			subtle.ConstantTimeCompare([]byte(pass), []byte(password)) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("You are Unauthorized to access the application.\n"))
			return
		}

		handler(w, r, accesses)
	}
}

func addRoutes(router *mux.Router) *mux.Router {
	for _, route := range routes {
		router.Methods(route.Method).Path(route.Pattern).Name(route.Name).
			Handler(route.HandlerFunc)

	}
	return router
}

func sendHome(w *http.ResponseWriter, accesses *[]Access) {
	parsedTemplates, _ := template.ParseFiles("templates/home.html")
	err := parsedTemplates.Execute(*w, *accesses)
	fmt.Fprintf(os.Stderr, err.Error())
}

func editServer(w http.ResponseWriter, r *http.Request, accesses *[]Access) {
	r.ParseForm()
	access := new(Access)
	if decodeErr := schema.NewDecoder().Decode(access, r.PostForm); decodeErr != nil {
		http.Error(w, decodeErr.Error(), http.StatusInternalServerError)
		return
	}

	accessToEditIndex := getIndexByID(access.ID, accesses)
	if accessToEditIndex < 0 {
		http.Error(w, "access to edit not found", http.StatusInternalServerError)
		return
	}

	(*accesses)[accessToEditIndex] = *access
}

func addServer(w http.ResponseWriter, r *http.Request, accesses *[]Access) {
	r.ParseForm()
	access := new(Access)
	if decodeErr := schema.NewDecoder().Decode(access, r.PostForm); decodeErr != nil {
		http.Error(w, decodeErr.Error(), http.StatusInternalServerError)
		return
	}

	index := getNextIndex(accesses)
	access.ID = index
	*accesses = append(*accesses, *access)

	if err := save(accesses); err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
	}
}

func deleteServer(w http.ResponseWriter, r *http.Request, accesses *[]Access) {
	r.ParseForm()
	access := new(
		struct {
			ID int `json:"id"`
		})
	if decodeErr := schema.NewDecoder().Decode(access, r.PostForm); decodeErr != nil {
		http.Error(w, decodeErr.Error(), http.StatusInternalServerError)
		return
	}
	removeElementByID(access.ID, accesses)
	if err := save(accesses); err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
	}
}

func save(accesses *[]Access) error {
	var sb strings.Builder
	sb.WriteString("Destination (Server Name),User (Destination),Access from,Notes\n")
	for _, a := range *accesses {
		sb.WriteString(a.ToCSV())
		sb.WriteString("\n")
	}
	return ioutil.WriteFile(*inputFile, []byte(sb.String()), 0644)
}

func main() {

	inputFile = flag.String("input", "", "csv file")
	flag.Parse()

	if len(*inputFile) == 0 {
		log.Fatal("input csv file missing")
	}

	file, err := os.Open(*inputFile)
	accesses, err := extractAccessesFromFile(file)
	if err != nil {
		panic(err)
	}
	file.Close()

	router := NewRouter(&accesses)

	if err = http.ListenAndServe(connHost+":"+connPort, router); err != nil {
		log.Fatal("error starting http server: ", err)
	}
}
