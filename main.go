package main

import (
	"crypto/subtle"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// Route ...
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Access ...
type Access struct {
	ID                                              int
	ServerDestination, UserDestination, From, Notes string
}

// HandlerFunc2 ...
type HandlerFunc2 func(http.ResponseWriter, *http.Request, *[]Access)

// Routes ...
type Routes []Route

const enterYourUserNamePassword = "Please enter your username and password"

var (
	cbeUser     = "root"
	cbePassword = "lein23"
	connHost    = "localhost"
	connPort    = "8081"

	routes = Routes{

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
		if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(cbeUser)) != 1 ||
			subtle.ConstantTimeCompare([]byte(pass), []byte(cbePassword)) != 1 {
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
		router.Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)

	}
	return router
}

func homePage(w http.ResponseWriter, r *http.Request, accesses *[]Access) {
	parsedTemplates, _ := template.ParseFiles("templates/index.html")
	err := parsedTemplates.Execute(w, *accesses)
	if err != nil {
		log.Print("Error occurred while executing the template or writing its output: ", err)
		return
	}
}

func searchByID(id int, accesses *[]Access) (Access, bool) {
	found := false
	var acc Access
	for _, a := range *accesses {
		if a.ID == id {
			acc = a
			found = true
			break
		}
	}

	return acc, found
}

func sendHome(w *http.ResponseWriter, accesses *[]Access) {
	parsedTemplates, _ := template.ParseFiles("templates/home.html")
	err := parsedTemplates.Execute(*w, *accesses)
	fmt.Fprintf(os.Stderr, err.Error())
}

func editPage(w http.ResponseWriter, r *http.Request, accesses *[]Access) {

	vars := mux.Vars(r)
	idParam := vars["id"]

	id, err := strconv.Atoi(idParam)
	if err != nil {
		sendHome(&w, accesses)
		return
	}
	access, found := searchByID(id, accesses)
	if !found {
		sendHome(&w, accesses)
		return
	}

	parsedTemplates, _ := template.ParseFiles("templates/edit.html")
	err = parsedTemplates.Execute(w, access)
	if err != nil {
		log.Print("Error occurred while executing the template or writing its output: ", err)
		return
	}
}

func main() {

	inputFile := flag.String("input", "", "csv file")
	flag.Parse()

	if len(*inputFile) == 0 {
		log.Fatal("input csv file missing")
	}

	file, err := os.Open(*inputFile)
	accesses, err := extractAccessesFromFile(file)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	router := mux.NewRouter().StrictSlash(false)
	router = addRoutes(router)

	router.HandleFunc("/", auth(homePage, enterYourUserNamePassword, &accesses))
	router.HandleFunc("/edit/{id}", auth(editPage, enterYourUserNamePassword, &accesses))
	// router.HandleFunc("/personas", auth(personasPage, enterYourUserNamePassword))
	// router.HandleFunc("/stats", auth(statsPage, enterYourUserNamePassword))
	// router.HandleFunc("/person/{id}", auth(personInfoPage, enterYourUserNamePassword))
	// router.HandleFunc("/interaction/{id}", auth(interactionInfoPage, enterYourUserNamePassword))

	router.PathPrefix("/").Handler(http.StripPrefix("/static", http.FileServer(http.Dir("static/"))))

	fileSave := time.NewTicker(10 * time.Second)
	go func(tick *time.Ticker) {
		for {
			select {
			case <-fileSave.C:
				fmt.Println("Holis ... ")
			}
		}
	}(fileSave)

	fmt.Println("Starting server ...")
	err = http.ListenAndServe(connHost+":"+connPort, router)
	if err != nil {
		log.Fatal("error starting http server: ", err)
		return
	}
}
