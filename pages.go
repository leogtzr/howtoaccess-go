package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func homePage(w http.ResponseWriter, r *http.Request, accesses *[]Access) {
	parsedTemplates, _ := template.ParseFiles("templates/index.html")
	err := parsedTemplates.Execute(w, *accesses)
	if err != nil {
		log.Print("Error occurred while executing the template or writing its output: ", err)
		return
	}
}

func addPage(w http.ResponseWriter, r *http.Request, accesses *[]Access) {
	parsedTemplates, _ := template.ParseFiles("templates/add.html")
	err := parsedTemplates.Execute(w, *accesses)
	if err != nil {
		log.Print("Error occurred while executing the template or writing its output: ", err)
		return
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
	err = parsedTemplates.Execute(w, access)
	if err != nil {
		log.Print("Error occurred while executing the template or writing its output: ", err)
		return
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
	err = parsedTemplates.Execute(w, access)
	if err != nil {
		log.Print("Error occurred while executing the template or writing its output: ", err)
		return
	}
}
