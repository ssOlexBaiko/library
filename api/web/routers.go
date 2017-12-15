package web

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {

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

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	Route{
		"BooksIndex",
		"GET",
		"/books",
		BooksIndex,
	},
	Route{
		"BookCreate",
		"POST",
		"/books",
		BookCreate,
	},
	Route{
		"GetBook",
		"GET",
		"/books/{id}",
		GetBook,
	},
	Route{
		"RemoveBook",
		"Delete",
		"/books/{id}",
		RemoveBook,
	},
}
