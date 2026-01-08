package main

import (
	"flag"
	"httpRandom/app/randomer"
	"net/http"
)

func main() {
	number := flag.Int("number", 6, "enter max value for random number")
	flag.Parse()
	router := http.NewServeMux()
	randomer.NewRandomHandler(router, *number)
	server := http.Server{
		Addr:    ":8081",
		Handler: router,
	}
	server.ListenAndServe()
}
