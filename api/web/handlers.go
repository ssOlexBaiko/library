package web

import (
	"fmt"
	"net/http"
	"encoding/json"
	"github.com/library/storage"
	"log"
	"io/ioutil"
	"io"
	"github.com/twinj/uuid"
)

var db = storage.InitDB()

func RepoCreateBook (b storage.Book) storage.Book {
	b.ID = fmt.Sprint(uuid.NewV4())
	db = append(db, b)
	return b
}

func Index(w http.ResponseWriter, r *http.Request) {
	//w.WriteHeader(http.StatusOK)
	//fmt.Fprintf(w, "Hello, this is the library resource")
	if err := json.NewEncoder(w).Encode(storage.GetBooks()); err != nil {
		log.Fatal(err)
	}
	w.WriteHeader(http.StatusOK)
}

func BooksIndex (w http.ResponseWriter, r *http.Request) {
	if err := json.NewEncoder(w).Encode(db); err != nil {
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
	b := RepoCreateBook(book)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(b); err != nil {
		panic(err)
	}
}
