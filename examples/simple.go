package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dbubel/intake"
)

func stream(w http.ResponseWriter, r *http.Request) {
	js := intake.NewStreamingJSONWriter(&w)
	js.Write(map[string]string{
		"hello": "test",
	})
	js.Write(map[string]string{
		"hello": "world",
	})
}

func testSimple(w http.ResponseWriter, r *http.Request) {
	intake.RespondJSON(w, r, http.StatusOK, map[string]string{
		"status": "world",
	})
	return
}

func middlewareOne(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			var s int
			intake.FromContext(r, "response-code", &s)
		}()
		next(w, r)
	}
}

func middlewareTwo(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			var s int
			intake.FromContext(r, "response-code", &s)
		}()
		next(w, r)
	}
}

func main() {
	app := intake.New()

	app.AddEndpoint("/hello", http.MethodGet, func(w http.ResponseWriter, t *http.Request) {
		fmt.Fprint(w, "hello, world")
	})

	app.Run(&http.Server{
		Addr:           ":8000",
		Handler:        app.Router,
		ReadTimeout:    time.Second * 30,
		WriteTimeout:   time.Second * 30,
		MaxHeaderBytes: 1 << 20,
	})
}
