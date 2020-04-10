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
	"github.com/gorilla/schema"
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
	ID                int    `json:"id"`
	ServerDestination string `json:"serverDestination"`
	UserDestination   string `json:"userDestination"`
	From              string `json:"from"`
	Notes             string `json:"notes"`
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

func getIndexByID(id int, accesses *[]Access) int {
	idx := -1
	for i, a := range *accesses {
		if a.ID == id {
			idx = i
			break
		}
	}
	return idx
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

func editServer(w http.ResponseWriter, r *http.Request, accesses *[]Access) {

	r.ParseForm()
	access := new(Access)
	decoder := schema.NewDecoder()
	decodeErr := decoder.Decode(access, r.PostForm)
	if decodeErr != nil {
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

	router.
		Methods("POST").Path("/editserver").
		Name("editserver").
		Handler(auth(editServer, enterYourUserNamePassword, &accesses))

	router.HandleFunc("/", auth(homePage, enterYourUserNamePassword, &accesses))
	router.HandleFunc("/edit/{id}", auth(editPage, enterYourUserNamePassword, &accesses))
	// router.HandleFunc("/personas", auth(personasPage, enterYourUserNamePassword))
	// router.HandleFunc("/stats", auth(statsPage, enterYourUserNamePassword))
	// router.HandleFunc("/person/{id}", auth(personInfoPage, enterYourUserNamePassword))
	// router.HandleFunc("/interaction/{id}", auth(interactionInfoPage, enterYourUserNamePassword))

	router.PathPrefix("/").Handler(http.StripPrefix("/static", http.FileServer(http.Dir("static/"))))

	fileSave := time.NewTicker(10 * time.Second)
	go func(tick *time.Ticker, accesses *[]Access) {
		for {
			select {
			case <-fileSave.C:
				fmt.Printf("Time: %s\n", time.Now().String())
			}
		}
	}(fileSave, &accesses)

	fmt.Println("Starting server ...")
	err = http.ListenAndServe(connHost+":"+connPort, router)
	if err != nil {
		log.Fatal("error starting http server: ", err)
		return
	}
}
