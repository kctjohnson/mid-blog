package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/kctjohnson/mid-blog/internal/templates/pages"
	"github.com/kctjohnson/mid-blog/internal/templates/pages/admin"
	"github.com/kctjohnson/mid-blog/internal/templates/pages/public"
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

	pages.Index(posts).Render(r.Context(), w)
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

	public.Post(*post, comments).Render(r.Context(), w)
}

func (app Application) Admin(w http.ResponseWriter, r *http.Request) {
	admin.Index().Render(r.Context(), w)
}

func (app Application) AdminPosts(w http.ResponseWriter, r *http.Request) {
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

	admin.Posts(posts).Render(r.Context(), w)
}

func (app Application) AdminPost(w http.ResponseWriter, r *http.Request) {
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

	admin.Post(*post).Render(r.Context(), w)
}

func (app Application) AdminBloggers(w http.ResponseWriter, r *http.Request) {
	bloggers, err := app.BloggerRepo.All()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	admin.Bloggers(bloggers).Render(r.Context(), w)
}

func (app Application) AdminBlogger(w http.ResponseWriter, r *http.Request) {
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

	blogger, err := app.BloggerRepo.FindByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	admin.Blogger(*blogger).Render(r.Context(), w)
}

func (app Application) AdminUsers(w http.ResponseWriter, r *http.Request) {
	users, err := app.UserRepo.All()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	admin.Users(users).Render(r.Context(), w)
}

func (app Application) AdminUser(w http.ResponseWriter, r *http.Request) {
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

	user, err := app.UserRepo.FindByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	comments, err := app.UserRepo.Comments(user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for i := range comments {
		comments[i].User = user
	}

	admin.User(*user).Render(r.Context(), w)
}
