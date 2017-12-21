package main

import (
	"log"
	"net/http"

	"flag"

	"github.com/ssOlexBaiko/library/api/web"
)

func main() {
	flag.Parse()
	router := web.NewRouter()
	log.Fatal(http.ListenAndServe("0.0.0.0:8000", router))
}
