package main

import (
	"fmt"
	"net/http"

	"github.com/khalidibnwalid/Luma/core"
)

func TestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("1 Start")
		next.ServeHTTP(w, r)
		fmt.Println("1 done")
	})
}

func Test2Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("2 Start")
		next.ServeHTTP(w, r)
		fmt.Println("2 done")
	})
}

func main() {
	app := core.NewApp()
	app.Use(TestMiddleware)
    app.Use(Test2Middleware)
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
