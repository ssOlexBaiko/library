package main

import (
	"log"
	"net/http"
	"github.com/ssOlexBaiko/library/api/web"
)

func main() {
	router := web.NewRouter()
	log.Fatal(http.ListenAndServe("0.0.0.0:8000", router))
}