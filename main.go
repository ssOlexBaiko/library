package main

import (
	"log"
	"net/http"

	"flag"
	"os"
	"path/filepath"

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
		path       string
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

	path, err = filepath.Abs(*libPath)
	if err != nil {
		log.Fatal(err)
	}
	flag.Parse()
	if *useSql {
		sqlStorage, err = storage.InitDB()
		if err != nil {
			log.Fatal(err)
		}

		store = storage.NewSQLLibrary(sqlStorage)
	} else {
		file, err := os.OpenFile(path, os.O_RDWR, 0660)
		defer file.Close()
		if err != nil {
			log.Fatal(err)
		}
		store = storage.NewLibrary(file)
	}

	router := web.NewRouter(
		web.NewHandler(store),
	)

	log.Fatal(http.ListenAndServe("0.0.0.0:8000", router))
}
