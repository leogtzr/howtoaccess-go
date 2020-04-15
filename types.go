package main

import "net/http"

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

// HealthState ...
type HealthState struct {
	State         int
	ErrorMessages []string
}
