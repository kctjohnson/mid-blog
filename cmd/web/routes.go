package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/kctjohnson/mid-blog/internal/db/models"
	"github.com/kctjohnson/mid-blog/internal/db/repos"
	"github.com/kctjohnson/mid-blog/internal/templates/pages/admin"
	"github.com/kctjohnson/mid-blog/internal/templates/pages/public"
)

func (app Application) Index(w http.ResponseWriter, r *http.Request) {
	allPosts, err := app.PostRepo.All()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for i := range allPosts {
		blogger, err := app.BloggerRepo.FindByID(allPosts[i].BloggerID)
		if err != nil {
			return
		}
		allPosts[i].Blogger = blogger
	}

	trendingPosts, err := app.PostRepo.GetTrending(3)
	for i := range trendingPosts {
		blogger, err := app.BloggerRepo.FindByID(trendingPosts[i].BloggerID)
		if err != nil {
			return
		}
		trendingPosts[i].Blogger = blogger
	}

	userData := app.SessionManager.Get(r.Context(), "user_data")
	if userData == nil {
		public.Index(nil, trendingPosts, allPosts).Render(r.Context(), w)
	} else {
		public.Index(userData.(UserInfo).User, trendingPosts, allPosts).Render(r.Context(), w)
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

func (app Application) Comment(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Error: %s", err.Error())))
		return
	}

	postIDStr := r.FormValue("post_id")
	userIDStr := r.FormValue("user_id")
	content := r.FormValue("content")

	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = app.CommentRepo.Insert(repos.CommentInsertParameters{
		UserID:  userID,
		PostID:  postID,
		Content: content,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the user
	user, err := app.UserRepo.FindByID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the post
	post, err := app.PostRepo.FindByID(postID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the comments on the post
	comments, err := app.PostRepo.Comments(post.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Fill in the users on all of the comments
	for i := range comments {
		user, err := app.UserRepo.FindByID(comments[i].UserID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		comments[i].User = user
	}

	public.CommentsSection(user, *post, comments).Render(r.Context(), w)
}

func (app Application) LikeComment(w http.ResponseWriter, r *http.Request) {
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

	comment, err := app.CommentRepo.FindByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newLikes := comment.Likes + 1
	comment, err = app.CommentRepo.Update(repos.CommentUpdateParameters{
		ID:    comment.ID,
		Likes: &newLikes,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userData := app.SessionManager.Get(r.Context(), "user_data")
	public.CommentStats(userData.(UserInfo).User, *comment).Render(r.Context(), w)
}

func (app Application) DislikeComment(w http.ResponseWriter, r *http.Request) {
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

	comment, err := app.CommentRepo.FindByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newDislikes := comment.Dislikes + 1
	comment, err = app.CommentRepo.Update(repos.CommentUpdateParameters{
		ID:       comment.ID,
		Dislikes: &newDislikes,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userData := app.SessionManager.Get(r.Context(), "user_data")
	public.CommentStats(userData.(UserInfo).User, *comment).Render(r.Context(), w)
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

	userData := app.SessionManager.Get(r.Context(), "user_data")
	public.PostStats(userData.(UserInfo).User, *post).Render(r.Context(), w)
}

func (app Application) Blogger(w http.ResponseWriter, r *http.Request) {
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

	posts, err := app.BloggerRepo.Posts(blogger.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for i := range posts {
		posts[i].Blogger = blogger
	}

	userData := app.SessionManager.Get(r.Context(), "user_data")
	if userData == nil {
		public.Blogger(nil, *blogger, posts).Render(r.Context(), w)
	} else {
		public.Blogger(userData.(UserInfo).User, *blogger, posts).Render(r.Context(), w)
	}
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

	userData := app.SessionManager.Get(r.Context(), "user_data")
	public.PostStats(userData.(UserInfo).User, *post).Render(r.Context(), w)
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
			firstName, err := app.ContentCreator.GenerateFirstName(1.0)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			lastName, err := app.ContentCreator.GenerateLastName(1.0)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			emailName, err := app.ContentCreator.GenerateEmail(1.0)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			email := fmt.Sprintf("%s@mid.com", emailName)

			age := rand.Int()%45 + 15

			var gender models.Gender
			randGenderInt := rand.Intn(2)
			if randGenderInt == 0 {
				gender = models.Male
			} else {
				gender = models.Female
			}

			bio, err := app.ContentCreator.GenerateBio(1.5, firstName, lastName)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			_, err = app.BloggerRepo.Insert(repos.BloggerInsertParameters{
				FirstName: firstName,
				LastName:  lastName,
				Email:     email,
				Age:       age,
				Gender:    gender,
				Bio:       bio,
				Avatar:    rand.Intn(9) + 1,
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
			title, err := app.ContentCreator.GenerateTitle(1.5, blogger.Bio)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			content, err := app.ContentCreator.GeneratePost(1.0, title)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

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

func (app Application) CreateRandomBlogger(w http.ResponseWriter, r *http.Request) {
	firstName, err := app.ContentCreator.GenerateFirstName(1.0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	lastName, err := app.ContentCreator.GenerateLastName(1.0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	emailName, err := app.ContentCreator.GenerateEmail(1.0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	email := fmt.Sprintf("%s@mid.com", emailName)

	age := rand.Int()%45 + 15

	var gender models.Gender
	randGenderInt := rand.Intn(2)
	if randGenderInt == 0 {
		gender = models.Male
	} else {
		gender = models.Female
	}

	bio, err := app.ContentCreator.GenerateBio(1.5, firstName, lastName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = app.BloggerRepo.Insert(repos.BloggerInsertParameters{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Age:       age,
		Gender:    gender,
		Bio:       bio,
		Avatar:    rand.Intn(9) + 1,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Created blogger!"))
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
	title, err := app.ContentCreator.GenerateTitle(1.5, blogger.Bio)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	content, err := app.ContentCreator.GeneratePost(1.0, title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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
