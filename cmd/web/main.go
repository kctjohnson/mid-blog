package main

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kctjohnson/mid-blog/internal/config"
	"github.com/kctjohnson/mid-blog/internal/db"
	"github.com/kctjohnson/mid-blog/internal/db/repos"
	"github.com/kctjohnson/mid-blog/internal/templates"
	ut "github.com/kctjohnson/mid-blog/internal/templates/utils"
)

type Application struct {
	BloggerRepo *repos.BloggerRepository
	CommentRepo *repos.CommentRepository
	PostRepo    *repos.PostRepository
	UserRepo    *repos.UserRepository
}

func main() {
	cfg := config.New(".")

	db, err := db.New(cfg.DSN, cfg.AutoMigrate)
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

	r.Route("/admin", func(r chi.Router) {
		r.Get("/", app.Admin)
		r.Get("/posts", app.AdminPosts)
		r.Get("/posts/{id}", templ.Handler(ut.SkeletonPostCard()).ServeHTTP)
		r.Get("/bloggers", app.AdminBloggers)
		r.Get("/bloggers/{id}", templ.Handler(ut.SkeletonPostCard()).ServeHTTP)
		r.Get("/users", app.AdminUsers)
		r.Get("/users/{id}", templ.Handler(ut.SkeletonPostCard()).ServeHTTP)
	})

	fmt.Println("Listening on :4231")
	if err := http.ListenAndServe(":4231", r); err != nil {
		return err
	}

	return nil
}
