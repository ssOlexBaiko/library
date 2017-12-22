package web

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Route describes route object
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes contains route objects
type Routes []Route

// NewRouter initialize app routers
func NewRouter(handler *handler) *mux.Router {
	var routes = Routes{
		{"Index", "GET", "/", handler.IndexHandler},
		{"BooksIndex", "GET", "/books", handler.BooksIndexHandler},
		{"BookCreate", "POST", "/books", handler.BookCreateHandler},
		{"GetBook", "GET", "/books/{id}", handler.GetBookHandler},
		{"RemoveBook", "Delete", "/books/{id}", handler.RemoveBookHandler},
		{"ChangeBook", "PUT", "/books/{id}", handler.ChangeBookHandler},
		{"BookFilter", "POST", "/books/filter", handler.BookFilterHandler},
	}

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	return router
}
