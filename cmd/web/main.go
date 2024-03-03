package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kctjohnson/mid-blog/internal/config"
	"github.com/kctjohnson/mid-blog/internal/db"
	"github.com/kctjohnson/mid-blog/internal/db/repos"
	"github.com/kctjohnson/mid-blog/internal/templates"
)

type Application struct {
	BloggerRepo *repos.BloggerRepository
	CommentRepo *repos.CommentRepository
	PostRepo    *repos.PostRepository
	UserRepo    *repos.UserRepository
}

func main() {
	cfg := config.New(".")

	db, err := db.New(cfg.DSN, true)
	if err != nil {
		panic(err)
	}

	bloggerRepo := repos.NewBloggerRepository(db)
	commentRepo := repos.NewCommentRepository(db)
	postRepo := repos.NewPostRepository(db)
	userRepo := repos.NewUserRepository(db)

	app := Application{
		BloggerRepo: bloggerRepo,
		CommentRepo: commentRepo,
		PostRepo:    postRepo,
		UserRepo:    userRepo,
	}

	app.StartServer()
}

func (app Application) StartServer() error {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Handle("/static/*", http.FileServer(http.FS(templates.Static)))

	r.Get("/", app.Index)
	r.Get("/{id}", app.Post)

	fmt.Println("Listening on :4231")
	if err := http.ListenAndServe(":4231", r); err != nil {
		return err
	}

	return nil
}
