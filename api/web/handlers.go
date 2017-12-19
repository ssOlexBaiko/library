package web

import (
	"fmt"
	"net/http"
	"encoding/json"
	"github.com/ssOlexBaiko/library/storage"
	"log"
	"github.com/gorilla/mux"
)
// IndexHandler handles requests with GET method
func IndexHandler(w http.ResponseWriter, _ *http.Request) {
	log.Println("Index - call")
	w.WriteHeader(http.StatusOK)
	_, err := fmt.Fprintf(w, "Hello, this is the library resource")
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// BooksIndexHandler handles requests with GET method
func BooksIndexHandler(w http.ResponseWriter, _ *http.Request) {
	log.Println("BooksIndex - call")
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

// BookCreateHandler handles requests with POST method
func BookCreateHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("BookCreate - call")
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

// GetBookHandler handles requests with GET method
func GetBookHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("GetBook - call")
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

// RemoveBookHandler handles requests with DELETE method
func RemoveBookHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("RemoveBook - call")
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

// ChangeBookHandler handles requests with PUT method
func ChangeBookHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("ChangeBook - call")
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

// BookFilterHandler handles requests with POST method
func BookFilterHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("BookFilter - call")
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
