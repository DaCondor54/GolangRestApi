package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
)

func IndexHandle(writer http.ResponseWriter, request *http.Request) {
	log.Println("Request Index")
	writer.Header().Set("Content-Type", "application/json")

	urlExample := "postgres://user:password@localhost:5432/dbname"
	conn, err := pgx.Connect(context.Background(), urlExample)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return
	}
	defer conn.Close(context.Background())

	rows, err := conn.Query(context.Background(), "SELECT * FROM manga;")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to Query into manga DB: %v\n", err)
		return
	}
	if rows.Err() != nil {
		fmt.Fprintf(os.Stderr, "Unable to Query Row into mangas DB: %v\n", rows.Err())
		return
	}

	var mangaCollection Mangas
	for rows.Next() {
		var manga Manga
		if err != rows.Scan(&manga.BookId, &manga.Title, &manga.Year, &manga.Volumes) {
			fmt.Fprintf(os.Stderr, "Unable to Scan Rows into mangas DB: %v\n", err)
			return
		}
		mangaCollection.AddManga(&manga)
	}

	if err := json.NewEncoder(writer).Encode(&mangaCollection); err != nil {
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

	urlExample := "postgres://user:password@localhost:5432/dbname"
	conn, err := pgx.Connect(context.Background(), urlExample)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return
	}
	defer conn.Close(context.Background())
	if _, err := conn.Exec(context.Background(), "INSERT INTO  manga(title, year, volumes) VALUES($1, $2, $3);", manga.Title, manga.Year, manga.Volumes); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to insert into mangas DB: %v\n", err)
		return
	}

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

	urlExample := "postgres://user:password@localhost:5432/dbname"
	conn, err := pgx.Connect(context.Background(), urlExample)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return
	}
	defer conn.Close(context.Background())
	for _, manga := range mangas.MangaArray {
		if _, err := conn.Exec(context.Background(), "INSERT INTO  manga(title, year, volumes) VALUES($1, $2, $3);", manga.Title, manga.Year, manga.Volumes); err != nil {
			fmt.Fprintf(os.Stderr, "Unable to insert into mangas DB: %v\n", err)
			return
		}
	}

	if err := json.NewEncoder(writer).Encode(&mangas); err != nil {
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

	var manga Manga
	if err := json.NewDecoder(request.Body).Decode(&manga); err != nil {
		log.Println("Could not decode manga")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	urlExample := "postgres://user:password@localhost:5432/dbname"
	conn, err := pgx.Connect(context.Background(), urlExample)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return
	}
	defer conn.Close(context.Background())

	if _, err := conn.Exec(context.Background(), "UPDATE manga SET title=$1, year=$2, volumes=$3 WHERE book_id=$4", manga.Title, manga.Year, manga.Volumes, idInt); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to conn Exec Update to database: %v\n", err)
		return
	}
	if err := json.NewEncoder(writer).Encode(manga); err != nil {
		log.Println("Could not encode manga")
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func ReadHandle(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Read Id Handle")
	writer.Header().Set("Content-Type", "application/json")
	id := mux.Vars(request)["id"]

	idInt, err := strconv.Atoi(id)
	if err != nil {
		log.Println("Could not parse Int")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	urlExample := "postgres://user:password@localhost:5432/dbname"
	conn, err := pgx.Connect(context.Background(), urlExample)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return
	}
	defer conn.Close(context.Background())
	row := conn.QueryRow(context.Background(), "SELECT book_id, title, year, volumes FROM manga WHERE book_id=$1;", idInt)
	var manga Manga
	if err := row.Scan(&manga.BookId, &manga.Title, &manga.Year, &manga.Volumes); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return
	}

	if err := json.NewEncoder(writer).Encode(manga); err != nil {
		log.Println("Could not Encode the manga")
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func limitNumClients(f http.HandlerFunc, maxClients int) http.HandlerFunc {
	sema := make(chan struct{}, maxClients)

	return func(w http.ResponseWriter, req *http.Request) {
		sema <- struct{}{}
		defer func() { <-sema }()
		f(w, req)
	}
}
