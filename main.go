package main

import (
	"net/http"
	"../library/api/web"
)

func main() {
	router := web.NewRouter()
	http.ListenAndServe("0.0.0.0:8000", router)
}