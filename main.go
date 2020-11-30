package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	PORT  = ":8080"
	INDEX = "templates/index.html"
)

func main() {

	go watcher()

	log.Print("Starting ethereum-data-fetcher...")

	log.Printf("Started Ethereum Data Fetcher on port [%v] ", PORT)

	router := mux.NewRouter()

	router.HandleFunc("/", RootHandler)

	router.HandleFunc("/query", QueryHandler)

	router.HandleFunc("/api", ApiHandler).Methods("POST")

	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	http.Handle("/", router)

	log.Fatal(http.ListenAndServe(":8080", nil))

}
