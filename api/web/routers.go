package web

import (
	"log"
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
		{"Index", "GET", "/",
			myFirstMiddlewareFunc(http.HandlerFunc(handler.IndexHandler))},
		{"BooksIndex", "GET", "/books",
			myFirstMiddlewareFunc(http.HandlerFunc(handler.BooksIndexHandler))},
		{"BookCreate", "POST", "/books",
			myFirstMiddlewareFunc(http.HandlerFunc(handler.BookCreateHandler))},
		{"GetBook", "GET", "/books/{id}",
			myFirstMiddlewareFunc(http.HandlerFunc(handler.GetBookHandler))},
		{"RemoveBook", "Delete", "/books/{id}",
			myFirstMiddlewareFunc(http.HandlerFunc(handler.RemoveBookHandler))},
		{"ChangeBook", "PUT", "/books/{id}",
			myFirstMiddlewareFunc(http.HandlerFunc(handler.ChangeBookHandler))},
		{"BookFilter", "POST", "/books/filter",
			myFirstMiddlewareFunc(http.HandlerFunc(handler.BookFilterHandler))},
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

func myFirstMiddlewareFunc(h http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Test middleware")
		h.ServeHTTP(w, r)
	})
}
