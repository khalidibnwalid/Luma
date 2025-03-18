package main

import (
	"fmt"
	"net/http"
)

func main() {
	v1 := http.NewServeMux()
	v1.Handle("/v1/", http.StripPrefix("/v1", routes()))

	server := http.Server{
		Addr:    ":8080",
		Handler: v1,
	}

	server.ListenAndServe()
}

func routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "hello\n")
	})

	return mux
}
