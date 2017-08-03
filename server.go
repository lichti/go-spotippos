package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	//Listening port
	listenAddr := fmt.Sprintf(":%v", 8000)
	log.Printf("Listening on %v...\n", listenAddr)

	//Routes
	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler).Methods("GET")
	r.HandleFunc("/provinces/{action}", provinceHandler).Methods("POST")
	r.HandleFunc("/properties", newPropertyPOSTHandler).Methods("POST")
	r.HandleFunc("/properties", searchPropertiesGETHandler).Methods("GET")
	r.HandleFunc("/properties/{action}", propertiesPOSTHandler).Methods("POST")
	r.HandleFunc("/properties/{id:[0-9]+}", propertiesGETHandler).Methods("GET")
	r.NotFoundHandler = http.HandlerFunc(notFoundHandler)
	//Print access log on Stdout
	loggedRouter := handlers.CombinedLoggingHandler(os.Stdout, r)

	// Listening for requests
	log.Fatal(http.ListenAndServe(listenAddr, loggedRouter))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	StatusUnauthorized(w, "Access Denied")
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	StatusNotFound(w, "Not Found")
}
