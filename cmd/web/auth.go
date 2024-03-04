package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/kctjohnson/mid-blog/internal/db/models"
	"github.com/kctjohnson/mid-blog/internal/db/repos"
	"golang.org/x/crypto/bcrypt"
)

type UserInfo struct {
	User *models.User
}

func (app Application) Register(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Error: %s", err.Error())))
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	email := r.FormValue("email")

	// Generate the hash and add the user to the users table
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Error: %s", err.Error())))
		return
	}

	newUser, err := app.UserRepo.Insert(repos.UserInsertParameters{
		Username: username,
		Password: string(hash),
		Email:    email,
	})
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Error: %s", err.Error())))
		return
	}

	// Authorize the user with a session auth token, then redirect them to the admin panel
	var info UserInfo
	info.User = newUser
	app.SessionManager.Put(r.Context(), "user_data", info)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app Application) Login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Error: %s", err.Error())))
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	// First, make sure that the user exists
	user, err := app.UserRepo.FindByUsername(username)
	if err != nil {
		if err == sql.ErrNoRows {
			w.Write([]byte(fmt.Sprintf("User \"%s\" Doesn't Exist", username)))
			return
		}
		w.Write([]byte(fmt.Sprintf("Error: %s", err.Error())))
		return
	}

	// Compare the entered password to the hash password using bcrypt
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Error: %s", err.Error())))
		return
	}

	// We good? Awesome, add the session auth token and redirect them to the admin panel
	app.SessionManager.Put(r.Context(), "user_data", UserInfo{
		User: user,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app Application) Logout(w http.ResponseWriter, r *http.Request) {
	// Delete all session stored values and redirect the user to the home page
	app.SessionManager.Destroy(r.Context())
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app Application) UnauthorizedPage(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
