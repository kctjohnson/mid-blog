package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/kctjohnson/mid-blog/internal/templates/pages"
)

func (app Application) Index(w http.ResponseWriter, r *http.Request) {
	posts, err := app.PostRepo.All()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for i := range posts {
		blogger, err := app.BloggerRepo.FindByID(posts[i].BloggerID)
		if err != nil {
			return
		}
		posts[i].Blogger = blogger
	}

	component := pages.Index(posts)
	component.Render(r.Context(), w)
}

func (app Application) Post(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	post, err := app.PostRepo.FindByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	blogger, err := app.BloggerRepo.FindByID(post.BloggerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	post.Blogger = blogger

	comments, err := app.PostRepo.Comments(post.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Fill in the comments with the user who made them
	for i := range comments {
		user, err := app.UserRepo.FindByID(comments[i].UserID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		comments[i].User = user
	}

	component := pages.Post(*post, comments)
	component.Render(r.Context(), w)
}
