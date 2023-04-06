package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var MangaArray Mangas

func main() {
	fmt.Println("Hello World")

	router := mux.NewRouter()
	router.HandleFunc("/", limitNumClients(IndexHandle, 5)).Methods("GET")
	router.HandleFunc("/create/manga", limitNumClients(CreateHandle, 5)).Methods("POST")
	router.HandleFunc("/create/mangas", limitNumClients(CreateManyHandle, 5)).Methods("POST")
	router.HandleFunc("/update/{id}", limitNumClients(UpdateHanlde, 5)).Methods("PUT")
	router.HandleFunc("/read/{id}", limitNumClients(ReadHandle, 5)).Methods("GET")

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatalln("Server Couldn't Start")
	}
}
