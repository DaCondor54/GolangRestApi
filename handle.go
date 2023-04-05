package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func IndexHandle(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Request Index")

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	err := json.NewEncoder(writer).Encode(MangaArray)
	if err != nil {
		log.Fatalln("Couldn't encode")
	}
}

func CreateHandle(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Request Create")

	writer.Header().Set("Content-Type", "application/json")

	var manga Manga
	err := json.NewDecoder(request.Body).Decode(&manga)
	if err != nil {
		log.Println("Couldn't Decode Post Request")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	MangaArray.AddManga(manga)

	err = json.NewEncoder(writer).Encode(manga)
	if err != nil {
		log.Println("Couldn't Encode ")
	}
}

func UpdateHanlde(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Request Update")

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
}

func ReadHandle(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Read Handle")

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
}
