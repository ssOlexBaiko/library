package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ssOlexBaiko/library/storage"
	"github.com/twinj/uuid"
)

type handler struct {
	storage Storage
}

type Storage interface {
	GetBooks() (storage.Books, error)
	CreateBook(book storage.Book) error
	GetBook(id string) (storage.Book, error)
	RemoveBook(id string) error
	ChangeBook(changedBook storage.Book) (storage.Book, error)
	PriceFilter(filter storage.BookFilter) (storage.Books, error)
}

func NewHandler(storage Storage) *handler {
	return &handler{
		storage: storage,
	}
}

// IndexHandler handles requests with GET method
func (h *handler) IndexHandler(w http.ResponseWriter, _ *http.Request) {
	log.Println("Index - call")

	_, err := fmt.Fprint(w, "Hello, this is the library resource")
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		log.Println(err)
		return
	}
}

// BooksIndexHandler handles requests with GET method
func (h *handler) BooksIndexHandler(w http.ResponseWriter, _ *http.Request) {
	log.Println("BooksIndex - call")
	books, err := h.storage.GetBooks()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = json.NewEncoder(w).Encode(books)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// BookCreateHandler handles requests with POST method
func (h *handler) BookCreateHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("BookCreate - call")

	var book storage.Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		log.Println(err)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.storage.CreateBook(book)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
}

// GetBookHandler handles requests with GET method
func (h *handler) GetBookHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("GetBook - call")

	vars := mux.Vars(r)
	id, ok := vars["id"]
	_, err := uuid.Parse(id)
	if !ok || err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	book, err := h.storage.GetBook(id)
	if err != nil {
		if err == storage.ErrNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(book)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// RemoveBookHandler handles requests with DELETE method
func (h *handler) RemoveBookHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("RemoveBook - call")

	vars := mux.Vars(r)
	id, ok := vars["id"]
	_, err := uuid.Parse(id)
	if !ok || err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.storage.RemoveBook(id)
	if err != nil {
		if err == storage.ErrNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ChangeBookHandler handles requests with PUT method
func (h *handler) ChangeBookHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("ChangeBook - call")

	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	book, err := h.storage.GetBook(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		log.Println(err)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	book, err = h.storage.ChangeBook(book)
	if err != nil {
		if err == storage.ErrNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	err = json.NewEncoder(w).Encode(book)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// BookFilterHandler handles requests with POST method
func (h *handler) BookFilterHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("BookFilter - call")

	var filter storage.BookFilter
	err := json.NewDecoder(r.Body).Decode(&filter)
	if err != nil {
		log.Println(err)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	books, err := h.storage.PriceFilter(filter)
	if err != nil {
		log.Println(err)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(books)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
