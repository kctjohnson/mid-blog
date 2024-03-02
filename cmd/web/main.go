package main

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kctjohnson/mid-blog/internal/templates"
	"github.com/kctjohnson/mid-blog/internal/templates/pages"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Handle("/static/*", http.FileServer(http.FS(templates.Static)))

	r.Get("/", templ.Handler(pages.Index()).ServeHTTP)

	fmt.Println("Listening on :4231")
	if err := http.ListenAndServe(":4231", r); err != nil {
		panic(err)
	}
}
