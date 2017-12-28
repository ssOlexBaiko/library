package main

import (
	"log"
	"net/http"

	"flag"

	"github.com/jinzhu/gorm"
	"github.com/ssOlexBaiko/library/api/web"
	"github.com/ssOlexBaiko/library/storage"
)

var (
	libPath = flag.String("libPath", "storage/storage.json", "set path the storage file")
	useSql  = flag.Bool("useSql", false, "use sql db instead of json file")
)

func main() {
	var (
		sqlStorage *gorm.DB
		err        error
		store      web.Storage
	)

	defer func() {
		if sqlStorage != nil {
			err := sqlStorage.Close()
			if err != nil {
				log.Fatal(err)
			}
		}
	}()

	flag.Parse()
	switch *useSql {
	case true:
		sqlStorage, err = storage.InitDB()
		if err != nil {
			log.Fatal(err)
		}

		store = storage.NewSQLLibrary(sqlStorage)
	case false:
		store = storage.NewLibrary(*libPath)
	}

	router := web.NewRouter(
		web.NewHandler(store),
	)

	log.Fatal(http.ListenAndServe("0.0.0.0:8000", router))
}
