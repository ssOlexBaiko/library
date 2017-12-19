package web

import (
	"github.com/gorilla/mux"
	"net/http"
)

// Route describes route object
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

//Routes contains route objects
type Routes []Route

// NewRouter initialize app routers
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
		IndexHandler,
	},
	Route{
		"BooksIndex",
		"GET",
		"/books",
		BooksIndexHandler,
	},
	Route{
		"BookCreate",
		"POST",
		"/books",
		BookCreateHandler,
	},
	Route{
		"GetBook",
		"GET",
		"/books/{id}",
		GetBookHandler,
	},
	Route{
		"RemoveBook",
		"Delete",
		"/books/{id}",
		RemoveBookHandler,
	},
	Route{
		"ChangeBook",
		"PUT",
		"/books/{id}",
		ChangeBookHandler,
	},
	Route{
		"BookFilter",
		"POST",
		"/books/filter",
		BookFilterHandler,
	},
}
