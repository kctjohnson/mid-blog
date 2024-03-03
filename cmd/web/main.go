package main

import (
	"fmt"
	"net/http"

	"github.com/ayush6624/go-chatgpt"
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
	BloggerAI   *chatgpt.Client
}

func main() {
	cfg := config.New(".")

	client, err := chatgpt.NewClient(cfg.OpenAIKey)
	if err != nil {
		panic(err)
	}

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
		BloggerAI:   client,
	}

	app.StartServer()
}

func (app Application) StartServer() error {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Handle("/static/*", http.FileServer(http.FS(templates.Static)))

	r.Get("/", app.Index)
	r.Get("/{id}", app.Post)

	r.Route("/api", func(r chi.Router) {
		r.Post("/posts/{id}/like", app.LikePost)
		r.Post("/posts/{id}/dislike", app.DislikePost)
	})

	r.Route("/admin", func(r chi.Router) {
		r.Get("/", app.Admin)
		r.Get("/posts", app.AdminPosts)
		r.Get("/posts/{id}", app.AdminPost)
		r.Get("/bloggers", app.AdminBloggers)
		r.Get("/bloggers/{id}", app.AdminBlogger)
		r.Get("/users", app.AdminUsers)
		r.Get("/users/{id}", app.AdminUser)

		r.Delete("/posts/{id}", app.DeletePost)
		r.Delete("/bloggers/{id}", app.DeleteBlogger)
		r.Delete("/users/{id}", app.DeleteUser)

		r.Get("/initializebloggers", app.InitializeBloggers)

		r.Post("/createpostrand", app.CreatePostRandom)
	})

	fmt.Println("Listening on :4231")
	if err := http.ListenAndServe(":4231", r); err != nil {
		return err
	}

	return nil
}
