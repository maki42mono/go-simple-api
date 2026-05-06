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

type App struct {
	Books map[int]Book
	ID    int
}

func (app *App) incId() int {
	defer func() {
		app.ID++
	}()

	return app.ID
}

func NewApp() App {
	app := App{make(map[int]Book), 1}
	id := app.incId()
	app.Books[id] = Book{id, "Govyanka", "Griboedov"}
	id = app.incId()
	app.Books[id] = Book{id, "Nazdorovka", "Kristina"}

	return app
}

func (app *App) handleBooks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		app.getBooks(w)
	case http.MethodPost:
		app.addBook(w, r)
	default:
		message := fmt.Sprintf("The method %s is not implemented yet", r.Method)
		http.Error(w, message, http.StatusInternalServerError)
	}
}

func (app *App) getBooks(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(app.Books)
	if err != nil {
		fmt.Fprintf(w, "Something went wrong: %s", err.Error())
		http.Error(w, "my custom error", http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

}

func (app *App) addBook(w http.ResponseWriter, r *http.Request) {
	var input Book
	err := json.NewDecoder(r.Body).Decode(&input)

	if err != nil {
		fmt.Fprintf(w, "Something went wrong: %s", err.Error())
		http.Error(w, "my custom error", http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	input.ID = app.incId()
	app.Books[input.ID] = input
	url := fmt.Sprintf("/book/%d", input.ID)
	http.Redirect(w, r, url, http.StatusCreated)

}

func main() {
	app := NewApp()

	http.HandleFunc("/books", app.handleBooks)
	log.Fatal(http.ListenAndServe(":1234", nil))
}
