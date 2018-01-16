package main

import (
	"log"
	"net/http"

	"flag"

	"github.com/ssOlexBaiko/library/api/web"
	"github.com/ssOlexBaiko/library/storage"
)

var (
	fileLibPath = flag.String("libPath", "storage/storage.json", "set path the storage file")
	sqlLibPath  = flag.String("sqlLibPath", "storage/data.db", "set path the storage file")
)

func main() {
	var (
		err   error
		store web.Storage
	)

	flag.Parse()
	if len(*sqlLibPath) != 0 {
		store, err = storage.NewSQLLibrary(*sqlLibPath)
	} else {
		store, err = storage.NewLibrary(*fileLibPath)
	}
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := store.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	router := web.NewRouter(
		web.NewHandler(store),
	)

	log.Fatal(http.ListenAndServe("0.0.0.0:8000", router))
}
