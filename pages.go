package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

func homePage(w http.ResponseWriter, r *http.Request, accesses *[]Access) {
	parsedTemplates, _ := template.ParseFiles("templates/index.html")
	if err := parsedTemplates.Execute(w, *accesses); err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
	}
}

func addPage(w http.ResponseWriter, r *http.Request, accesses *[]Access) {
	parsedTemplates, _ := template.ParseFiles("templates/add.html")
	if err := parsedTemplates.Execute(w, *accesses); err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
	}
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
	if err = parsedTemplates.Execute(w, access); err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
	}
}

func deletePage(w http.ResponseWriter, r *http.Request, accesses *[]Access) {
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

	parsedTemplates, _ := template.ParseFiles("templates/delete.html")
	if err = parsedTemplates.Execute(w, access); err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		return
	}
}
