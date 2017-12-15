package web

import (
	"fmt"
	"net/http"
	"encoding/json"
	"github.com/library/storage"
	"log"
	"io/ioutil"
	"io"
	"github.com/gorilla/mux"
)

func Index(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello, this is the library resource")
}

func BooksIndex(w http.ResponseWriter, _ *http.Request) {
	books, err := storage.GetBooks()
	if err != nil {
		log.Fatal(err)
	}
	err = json.NewEncoder(w).Encode(books)
	if err != nil {
		log.Fatal(err)
	}
	w.WriteHeader(http.StatusOK)
}

func BookCreate(w http.ResponseWriter, r *http.Request) {
	var book storage.Book
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		log.Fatal(err)
	}
	if err := r.Body.Close(); err != nil {
		log.Fatal(err)
	}
	if err := json.Unmarshal(body, &book); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	err = storage.CreateBook(book)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
	}
}

func GetBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	book, err := storage.GetBook(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Fatal(err)

	}
	err = json.NewEncoder(w).Encode(book)
	if err != nil {
		log.Fatal(err)
	}
	w.WriteHeader(http.StatusOK)
}

func RemoveBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	err := storage.RemoveBook(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Fatal(err)
	}
	w.WriteHeader(http.StatusNoContent)
}
