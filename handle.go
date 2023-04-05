package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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

	id := mux.Vars(request)["id"]
	idInt, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("Couldnt parse parameter ID to Number")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	mangaIndex := MangaArray.SameId(idInt)
	if mangaIndex == -1 {
		log.Println("Wrong Index")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	var manga Manga
	if err := json.NewDecoder(request.Body).Decode(&manga); err != nil {
		log.Println("Could not decode manga")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	manga.BookId = idInt
	MangaArray.MangaArray[mangaIndex] = manga

	if err := json.NewEncoder(writer).Encode(manga); err != nil {
		log.Println("Could not encode manga")
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
}

func ReadHandle(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Read Id Handle")
	id := mux.Vars(request)["id"]

	idInt, err := strconv.Atoi(id)
	if err != nil {
		log.Println("Could not parse Int")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	mangaIndex := MangaArray.SameId(idInt)

	if err := json.NewEncoder(writer).Encode(MangaArray.MangaArray[mangaIndex]); err != nil {
		log.Println("Could not Encode the manga")
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusOK)
}
