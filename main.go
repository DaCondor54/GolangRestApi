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
	router.HandleFunc("/", IndexHandle).Methods("GET")
	router.HandleFunc("/create", CreateHandle).Methods("POST")
	router.HandleFunc("/update", IndexHandle).Methods("UPDATE")
	router.HandleFunc("/read", IndexHandle).Methods("GET")

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatalln("Server Couldn't Start")
	}
}
