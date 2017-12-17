package main

import (
	"net/http"
	"github.com/ssOlexBaiko/library/api/web"
)

func main() {
	router := web.NewRouter()
	http.ListenAndServe("0.0.0.0:8000", router)
}