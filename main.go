package main

import (
	"log"
	"net/http"

	"flag"

	"github.com/ssOlexBaiko/library/api/web"
	"github.com/ssOlexBaiko/library/storage"
)

var libPath = flag.String("libPath", "storage/storage.json", "set path the storage file")

func main() {
	flag.Parse()
	library := storage.Lib{*libPath}
	_ = library // library has to be passed into web.NewRouter() as param
	router := web.NewRouter()

	log.Fatal(http.ListenAndServe("0.0.0.0:8000", router))
}
