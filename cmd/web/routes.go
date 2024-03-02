package main

import (
	"net/http"

	"github.com/kctjohnson/mid-blog/internal/templates/pages"
)

func (app Application) Index(w http.ResponseWriter, r *http.Request) {
	posts, err := app.PostRepo.All()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	component := pages.Index(posts)
	component.Render(r.Context(), w)
}
