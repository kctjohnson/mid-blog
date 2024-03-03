package main

import (
	"fmt"
	"math/rand"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kctjohnson/mid-blog/internal/config"
	"github.com/kctjohnson/mid-blog/internal/db"
	"github.com/kctjohnson/mid-blog/internal/db/models"
	"github.com/kctjohnson/mid-blog/internal/db/repos"
	"github.com/kctjohnson/mid-blog/internal/templates"
	"github.com/kctjohnson/mid-blog/internal/templates/utils"
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

		r.Get("/testcreate", app.CreateStuff)
	})

	fmt.Println("Listening on :4231")
	if err := http.ListenAndServe(":4231", r); err != nil {
		return err
	}

	return nil
}

func (app Application) CreateStuff(w http.ResponseWriter, r *http.Request) {
	me, err := app.UserRepo.Insert(repos.UserInsertParameters{
		Username: utils.WordGenerator.Word(),
		Password: utils.WordGenerator.Word(),
	})
	if err != nil {
		panic(err)
	}

	kris, err := app.BloggerRepo.Insert(repos.BloggerInsertParameters{
		FirstName: utils.WordGenerator.Word(),
		LastName:  utils.WordGenerator.Word(),
		Email: fmt.Sprintf(
			"%s_%s@mid.com",
			utils.WordGenerator.Word(),
			utils.WordGenerator.Word(),
		),
		Age:    rand.Int()%75 + 15,
		Gender: models.Male,
		Bio:    utils.WordGenerator.Sentence(),
	})
	if err != nil {
		panic(err)
	}

	post, err := app.PostRepo.Insert(repos.PostInsertParameters{
		BloggerID: kris.ID,
		Title:     utils.WordGenerator.Words(5),
		Content:   utils.WordGenerator.Paragraphs(3),
	})
	if err != nil {
		panic(err)
	}

	_, err = app.CommentRepo.Insert(repos.CommentInsertParameters{
		UserID:  me.ID,
		PostID:  post.ID,
		Content: utils.WordGenerator.Sentences(2),
	})
	if err != nil {
		panic(err)
	}

	w.Write([]byte("Created stuff!"))
}
