package main

import (
	"encoding/gob"
	"fmt"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
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

	BloggerAI      *chatgpt.Client
	SessionManager *scs.SessionManager
}

func main() {
	cfg := config.New(".env")

	client, err := chatgpt.NewClient(cfg.OpenAIKey)
	if err != nil {
		panic(err)
	}

	gob.Register(UserInfo{})
	sessionManager := scs.New()
	sessionManager.Lifetime = 24 * time.Hour

	db, err := db.New(cfg.DSN, cfg.AutoMigrate)
	if err != nil {
		panic(err)
	}

	bloggerRepo := repos.NewBloggerRepository(db)
	commentRepo := repos.NewCommentRepository(db)
	postRepo := repos.NewPostRepository(db)
	userRepo := repos.NewUserRepository(db)

	app := Application{
		BloggerRepo:    bloggerRepo,
		CommentRepo:    commentRepo,
		PostRepo:       postRepo,
		UserRepo:       userRepo,
		BloggerAI:      client,
		SessionManager: sessionManager,
	}

	app.StartServer()
}

func (app Application) StartServer() error {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Handle("/static/*", http.FileServer(http.FS(templates.Static)))

	r.Get("/", app.Index)
	r.Get("/{id}", app.Post)
	r.Get("/membership", app.Membership)

	r.Get("/login", app.LoginPage)
	r.Get("/register", app.RegisterPage)
	r.Get("/unauthorized", app.UnauthorizedPage)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", app.Login)
		r.Post("/register", app.Register)
		r.Get("/logout", app.Logout)
	})

	r.Route("/api", func(r chi.Router) {
		r.Post("/posts/comment", app.protectedRoute(app.Comment))
		r.Post("/posts/{id}/like", app.protectedRoute(app.LikePost))
		r.Post("/posts/{id}/dislike", app.protectedRoute(app.DislikePost))
		r.Post("/comments/{id}/like", app.protectedRoute(app.LikeComment))
		r.Post("/comments/{id}/dislike", app.protectedRoute(app.DislikeComment))
	})

	r.Route("/admin", func(r chi.Router) {
		r.Get("/", app.protectedRoute(app.Admin))
		r.Get("/posts", app.protectedRoute(app.AdminPosts))
		r.Get("/posts/{id}", app.protectedRoute(app.AdminPost))
		r.Get("/bloggers", app.protectedRoute(app.AdminBloggers))
		r.Get("/bloggers/{id}", app.protectedRoute(app.AdminBlogger))
		r.Get("/users", app.protectedRoute(app.AdminUsers))
		r.Get("/users/{id}", app.protectedRoute(app.AdminUser))

		r.Delete("/posts/{id}", app.protectedRoute(app.DeletePost))
		r.Delete("/bloggers/{id}", app.protectedRoute(app.DeleteBlogger))
		r.Delete("/users/{id}", app.protectedRoute(app.DeleteUser))

		r.Get("/initializebloggers", app.protectedRoute(app.InitializeBloggers))

		r.Post("/createpostrand", app.protectedRoute(app.CreatePostRandom))
	})

	fmt.Println("Listening on :4231")
	if err := http.ListenAndServe(":4231", app.SessionManager.LoadAndSave(r)); err != nil {
		return err
	}

	return nil
}

// Makes it so that calls to this route validate the auth token FIRST before passing through the page
func (app *Application) protectedRoute(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userDataInterface := app.SessionManager.Get(r.Context(), "user_data")

		if userDataInterface == nil {
			http.Redirect(w, r, "/unauthorized", http.StatusSeeOther)
			return
		}

		h(w, r)
	}
}
