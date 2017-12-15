package web

import (
	"fmt"
	"net/http"
	"encoding/json"
	"github.com/library/storage"
	"log"
	"io/ioutil"
	"io"
)

func Index(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello, this is the library resource")
}

func BooksIndex (w http.ResponseWriter, r *http.Request) {
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
