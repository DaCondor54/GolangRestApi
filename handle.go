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
	log.Println("Request Index")
	writer.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(writer).Encode(MangaArray); err != nil {
		log.Println("Couldn't encode")
		writer.WriteHeader(http.StatusInternalServerError)
	}
}

func CreateHandle(writer http.ResponseWriter, request *http.Request) {
	log.Println("Request Create")
	writer.Header().Set("Content-Type", "application/json")

	var manga Manga
	if err := json.NewDecoder(request.Body).Decode(&manga); err != nil {
		log.Println("Couldn't Decode Post Request")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	MangaArray.AddManga(manga)
	if err := json.NewEncoder(writer).Encode(manga); err != nil {
		log.Println("Couldn't Encode ")
		writer.WriteHeader(http.StatusInternalServerError)
	}
}

func CreateManyHandle(writer http.ResponseWriter, request *http.Request) {
	log.Println("Request Create Many")
	writer.Header().Set("Content-Type", "application/json")

	var mangas Mangas
	if err := json.NewDecoder(request.Body).Decode(&mangas); err != nil {
		log.Println("Couldn't Decode Mangas Collection")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	for _, v := range mangas.MangaArray {
		MangaArray.AddManga(v)
	}

	if err := json.NewEncoder(writer).Encode(MangaArray); err != nil {
		log.Println("Couldn't Encode Manga Collection")
		writer.WriteHeader(http.StatusInternalServerError)
		return
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
}

func limitNumClients(f http.HandlerFunc, maxClients int) http.HandlerFunc {
	sema := make(chan struct{})

	return func(w http.ResponseWriter, req *http.Request) {
		sema <- struct{}{}
		defer func() { <-sema }()
		f(w, req)
	}
}
