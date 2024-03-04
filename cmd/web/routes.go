package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/kctjohnson/mid-blog/internal/db/models"
	"github.com/kctjohnson/mid-blog/internal/db/repos"
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

	userData := app.SessionManager.Get(r.Context(), "user_data")
	if userData == nil {
		public.Index(nil, posts).Render(r.Context(), w)
	} else {
		public.Index(userData.(UserInfo).User, posts).Render(r.Context(), w)
	}
}

func (app Application) Membership(w http.ResponseWriter, r *http.Request) {
	userData := app.SessionManager.Get(r.Context(), "user_data")
	if userData == nil {
		public.Membership(nil).Render(r.Context(), w)
	} else {
		public.Membership(userData.(UserInfo).User).Render(r.Context(), w)
	}
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

	userData := app.SessionManager.Get(r.Context(), "user_data")
	if userData == nil {
		public.Post(nil, *post, comments).Render(r.Context(), w)
	} else {
		public.Post(userData.(UserInfo).User, *post, comments).Render(r.Context(), w)
	}
}

func (app Application) LikePost(w http.ResponseWriter, r *http.Request) {
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

	newLikes := post.Likes + 1
	post, err = app.PostRepo.Update(repos.PostUpdateParameters{
		ID:    post.ID,
		Likes: &newLikes,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(fmt.Sprintf("%d", post.Likes-post.Dislikes)))
}

func (app Application) DislikePost(w http.ResponseWriter, r *http.Request) {
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

	newDislikes := post.Dislikes + 1
	post, err = app.PostRepo.Update(repos.PostUpdateParameters{
		ID:       post.ID,
		Dislikes: &newDislikes,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(fmt.Sprintf("%d", post.Likes-post.Dislikes)))
}

func (app Application) Admin(w http.ResponseWriter, r *http.Request) {
	userData := app.SessionManager.Get(r.Context(), "user_data")
	admin.Index(*userData.(UserInfo).User).Render(r.Context(), w)
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

	userData := app.SessionManager.Get(r.Context(), "user_data")
	admin.Posts(*userData.(UserInfo).User, posts).Render(r.Context(), w)
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

	userData := app.SessionManager.Get(r.Context(), "user_data")
	admin.Post(*userData.(UserInfo).User, *post).Render(r.Context(), w)
}

func (app Application) DeletePost(w http.ResponseWriter, r *http.Request) {
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

	err = app.PostRepo.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/posts", http.StatusSeeOther)
}

func (app Application) AdminBloggers(w http.ResponseWriter, r *http.Request) {
	bloggers, err := app.BloggerRepo.All()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userData := app.SessionManager.Get(r.Context(), "user_data")
	admin.Bloggers(*userData.(UserInfo).User, bloggers).Render(r.Context(), w)
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

	userData := app.SessionManager.Get(r.Context(), "user_data")
	admin.Blogger(*userData.(UserInfo).User, *blogger).Render(r.Context(), w)
}

func (app Application) DeleteBlogger(w http.ResponseWriter, r *http.Request) {
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

	err = app.BloggerRepo.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/bloggers", http.StatusSeeOther)
}

func (app Application) AdminUsers(w http.ResponseWriter, r *http.Request) {
	users, err := app.UserRepo.All()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userData := app.SessionManager.Get(r.Context(), "user_data")
	admin.Users(*userData.(UserInfo).User, users).Render(r.Context(), w)
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

	userData := app.SessionManager.Get(r.Context(), "user_data")
	admin.User(*userData.(UserInfo).User, *user).Render(r.Context(), w)
}

func (app Application) DeleteUser(w http.ResponseWriter, r *http.Request) {
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

	err = app.UserRepo.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
}

func (app Application) InitializeBloggers(w http.ResponseWriter, r *http.Request) {
	var wg sync.WaitGroup

	// Make 5 AI bloggers
	fmt.Println("Creating bloggers...")
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			firstNameResp, err := app.BloggerAI.SimpleSend(
				context.Background(),
				"In one word, give me a random first name.",
			)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			firstName := firstNameResp.Choices[0].Message.Content

			lastNameResp, err := app.BloggerAI.SimpleSend(
				context.Background(),
				"In one word, give me a random last name.",
			)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			lastName := lastNameResp.Choices[0].Message.Content

			emailResp, err := app.BloggerAI.SimpleSend(
				context.Background(),
				"I need to make a new email, give me a random email name.",
			)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			email := fmt.Sprintf("%s@mid.com", emailResp.Choices[0].Message.Content)

			age := rand.Int()%45 + 15

			var gender models.Gender
			randGenderInt := rand.Intn(2)
			if randGenderInt == 0 {
				gender = models.Male
			} else {
				gender = models.Female
			}

			bioResp, err := app.BloggerAI.SimpleSend(
				context.Background(),
				"Write me a random bio that is AT MAX 255 characters long.",
			)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			bio := bioResp.Choices[0].Message.Content

			_, err = app.BloggerRepo.Insert(repos.BloggerInsertParameters{
				FirstName: firstName,
				LastName:  lastName,
				Email:     email,
				Age:       age,
				Gender:    gender,
				Bio:       bio,
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println("Created bloggers!")

	// Make 1 post per blogger
	bloggers, err := app.BloggerRepo.All()
	if err != nil {
		return
	}

	fmt.Println("Creating posts...")
	for _, blogger := range bloggers {
		wg.Add(1)
		go func(b models.Blogger) {
			titleResp, err := app.BloggerAI.SimpleSend(
				context.Background(),
				"Generate me the title of an article about absolutely anything (Tech, science, literary, political, etc.). Try to keep it around 5-8 words.",
			)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			title := strings.ReplaceAll(titleResp.Choices[0].Message.Content, "\"", "")

			contentResp, err := app.BloggerAI.SimpleSend(
				context.Background(),
				"Write me a 5 paragraph article about this title \""+title+"\"",
			)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			content := contentResp.Choices[0].Message.Content

			_, err = app.PostRepo.Insert(repos.PostInsertParameters{
				BloggerID: b.ID,
				Title:     title,
				Content:   content,
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			wg.Done()
		}(blogger)
	}
	wg.Wait()
	fmt.Println("Created posts!")

	w.Write([]byte("Created stuff!"))
}

func (app Application) CreatePostRandom(w http.ResponseWriter, r *http.Request) {
	// Get the bloggers
	bloggers, err := app.BloggerRepo.All()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get a random blogger
	blogger := bloggers[rand.Intn(len(bloggers))]

	// Get a random title
	titleResp, err := app.BloggerAI.SimpleSend(
		context.Background(),
		"Generate me the title of an article about absolutely anything (Tech, science, literary, political, etc.). Try to keep it around 5-8 words.",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	title := strings.ReplaceAll(titleResp.Choices[0].Message.Content, "\"", "")

	contentResp, err := app.BloggerAI.SimpleSend(
		context.Background(),
		"Write me a 5 paragraph article about this title \""+title+"\"",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	content := contentResp.Choices[0].Message.Content

	post, err := app.PostRepo.Insert(repos.PostInsertParameters{
		BloggerID: blogger.ID,
		Title:     title,
		Content:   content,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	admin.RandomPostLink(post.ID, post.Title).Render(r.Context(), w)
}

func (app Application) LoginPage(w http.ResponseWriter, r *http.Request) {
	public.Login().Render(r.Context(), w)
}

func (app Application) RegisterPage(w http.ResponseWriter, r *http.Request) {
	public.Register().Render(r.Context(), w)
}
