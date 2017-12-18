package web

import (
	"fmt"
	"net/http"
	"encoding/json"
	"github.com/ssOlexBaiko/library/storage"
	"log"
	"github.com/gorilla/mux"
)

func indexHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello, this is the library resource")
}

func booksIndexHandler(w http.ResponseWriter, _ *http.Request) {
	books, err := storage.GetBooks()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	err = json.NewEncoder(w).Encode(books)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func bookCreateHandler(w http.ResponseWriter, r *http.Request) {
	var book storage.Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		log.Println(err)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = storage.CreateBook(book)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
}

func getBookHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	book, err := storage.GetBook(id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	err = json.NewEncoder(w).Encode(book)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func removeBookHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	err := storage.RemoveBook(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func changeBookHandler(w http.ResponseWriter, r *http.Request) {
	var book storage.Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		log.Println(err)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	vars := mux.Vars(r)
	id := vars["id"]
	err = storage.ChangeBook(id, book)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func bookFilterHandler(w http.ResponseWriter, r *http.Request) {
	var filter storage.Filter
	err := json.NewDecoder(r.Body).Decode(&filter)
	if err != nil {
		log.Println(err)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	books, err := storage.PriceFilter(filter)
	if err != nil {
		log.Println(err)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = json.NewEncoder(w).Encode(books)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
