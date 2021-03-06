package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

// String ...
func (a Access) String() string {
	if len(a.Notes) > 0 {
		return fmt.Sprintf("%d: %s@%s from %s (%s)", a.ID, a.ServerDestination, a.UserDestination, a.From, a.Notes)
	}
	return fmt.Sprintf("%d: %s@%s from %s", a.ID, a.ServerDestination, a.UserDestination, a.From)
}

// ToCSV ...
func (a Access) ToCSV() string {
	return fmt.Sprintf("%s,%s,%s,%s", a.ServerDestination, a.UserDestination, a.From, a.Notes)
}

func extractAccessesFromFile(file io.Reader) ([]Access, error) {
	records := make([]Access, 0)

	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		return []Access{}, err
	}

	for i, record := range lines {
		if i == 0 {
			continue
		}
		if len(record) != requiredNumberOfFields {
			return []Access{}, fmt.Errorf("wrong number of fields in line: %d", i)
		}
		records = append(records, Access{
			ID:                i,
			ServerDestination: record[0],
			UserDestination:   record[1],
			From:              record[2],
			Notes:             record[3],
		})
	}

	return records, nil
}

func equal(a, b []Access) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func getNextIndex(accesses *[]Access) int {
	if len(*accesses) == 0 {
		return -1
	}
	lastAccessElement := (*accesses)[len(*accesses)-1]
	return lastAccessElement.ID + 1
}

func removeElementByID(id int, accesses *[]Access) {
	id = getIndexByID(id, accesses)
	if id == -1 {
		return
	}
	*accesses = append((*accesses)[:id], (*accesses)[id+1:]...)
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

// NewRouter ...
func NewRouter(accesses *[]Access) *mux.Router {
	router := mux.NewRouter().StrictSlash(false)
	router = addRoutes(router)

	router.
		Methods("POST").Path("/editserver").
		Name("editserver").
		Handler(auth(editServer, enterYourUserNamePassword, accesses))

	router.
		Methods("POST").Path("/addserver").
		Name("addserver").
		Handler(auth(addServer, enterYourUserNamePassword, accesses))

	router.
		Methods("POST").Path("/deleteserver").
		Name("deleteserver").
		Handler(auth(deleteServer, enterYourUserNamePassword, accesses))

	router.
		Methods("GET").Path("/healthcheck").
		Name("healthcheck").
		Handler(healthCheckHandler(accesses))

	router.HandleFunc("/", auth(homePage, enterYourUserNamePassword, accesses))
	router.HandleFunc("/edit/{id}", auth(editPage, enterYourUserNamePassword, accesses))
	router.HandleFunc("/delete/{id}", auth(deletePage, enterYourUserNamePassword, accesses))
	router.HandleFunc("/add.html", auth(addPage, enterYourUserNamePassword, accesses))
	router.PathPrefix("/").Handler(http.StripPrefix("/static", http.FileServer(http.Dir("static/"))))

	return router
}
