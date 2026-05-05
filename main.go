package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Book struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Autor string `json:"autor"`
}

var (
	books  = make(map[int]Book)
	nextID = 1
)

func init() {
	books[nextID] = Book{nextID, "My title 1", "Agathya"}
	nextID++
	books[nextID] = Book{nextID, "Dvornik", "Dvorezkii"}
	nextID++
}

func handleBooks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getBooks(w)
	default:
		message := fmt.Sprintf("The method %s is not implemented yet", r.Method)
		http.Error(w, message, http.StatusInternalServerError)
	}
}

func getBooks(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(books)
	if err != nil {
		fmt.Fprintf(w, "Something went wrong: %s", err.Error())
		http.Error(w, "my custom error", http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

}

func main() {
	http.HandleFunc("/books", handleBooks)
	log.Fatal(http.ListenAndServe(":1234", nil))
}
