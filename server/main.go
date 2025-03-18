package main

import (
	"fmt"
	"net/http"

	"github.com/khalidibnwalid/Luma/core"
	"github.com/khalidibnwalid/Luma/server/middlewares"
)

func main() {
	app := core.NewApp()
	app.Use(middlewares.Logging)
	app.HandleFunc("/hello", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "hello\n")
	})
	v1 := http.NewServeMux()
	v1.Handle("/v1/", http.StripPrefix("/v1", app.Mux))

	server := http.Server{
		Addr:    ":8080",
		Handler: v1,
	}

	server.ListenAndServe()
	// app.ha
}

// func routes() *http.ServeMux {
// 	mux := http.NewServeMux()
// 	mux.HandleFunc("/hello", func(w http.ResponseWriter, req *http.Request) {
// 		fmt.Fprintf(w, "hello\n")
// 	})

// 	return mux
// }
