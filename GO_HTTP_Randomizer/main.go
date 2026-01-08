package main

import (
	"httpRandom/app/randomer"
	"net/http"
)

func main() {
	router := http.NewServeMux()
	randomer.NewRandomHandler(router)
	server := http.Server{
		Addr: ":8081",
		Handler: router,
	}
	server.ListenAndServe()
}
