package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dbubel/intake"
)

func a(w http.ResponseWriter, r *http.Request) {
	intake.Respond(w, r, http.StatusOK, []byte("yo"))
}

func main() {
	app := intake.New()
	app.OptionsHandler(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, Authorization")
		w.WriteHeader(http.StatusOK)
	})
	app.AddEndpoint(http.MethodGet, "/hello", a, w, w2)

	app.Run(&http.Server{
		Addr:           ":8080",
		Handler:        app.Mux,
		ReadTimeout:    time.Second * 60,
		WriteTimeout:   time.Second * 60,
		MaxHeaderBytes: 1 << 20,
	})
}

func w(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(("sup"))
		next(w, r)
	}
}

func w2(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(("sup 2"))
		next(w, r)
	}
}
