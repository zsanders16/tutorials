package controllers

import (
	"fakeTwitter/models"
	"fakeTwitter/rand"
	"fakeTwitter/views"
	"fmt"
	"net/http"
)

// NewUsers instantiates all info for a user and returns a pointer to a user
func NewUsers(userService *models.UserService) *Users {
	return &Users{
		NewView:     views.NewView("bootstrap", "users/new"),
		LoginView:   views.NewView("bootstrap", "users/login"),
		UserService: userService,
	}
}

// Users hold all info pertaining to a User
type Users struct {
	NewView     *views.View
	LoginView   *views.View
	UserService *models.UserService
}

// New display the signup form to create a new user
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	if err := u.NewView.Render(w, nil); err != nil {
		panic(err)
	}
}

// UserLogin delivers the user login page
func (u *Users) UserLoginView(w http.ResponseWriter, r *http.Request) {
	if err := u.LoginView.Render(w, nil); err != nil {
		panic(err)
	}
}

// UserLogin authenticates a user login attempt
func (u *Users) UserLogin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		panic(err)
	}

	type loginForm struct {
		Email    string `schema:"email"`
		Password string `schema:"password"`
	}

	var lf loginForm
	if err := parseForm(r, &lf); err != nil {
		panic(err)
	}

	user, err := u.UserService.AuthenticateUser(lf.Email, lf.Password)
	if err != nil {
		switch err {
		case models.ErrNotFound:
			fmt.Fprintln(w, "Invalid email address.")
		case models.ErrInvalidPassword:
			fmt.Fprintln(w, "Invalid password provided")
		default:
			fmt.Fprintln(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	err = u.signIn(w, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	http.Redirect(w, r, "/cookie", http.StatusFound)
}

// Create creates a new user from the submitted signuup form delived from the New func
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintln(w, "Temp response from controllers/users.create")
	if err := r.ParseForm(); err != nil {
		panic(err)
	}

	type NewUserForm struct {
		Name     string `schema:"name"`
		Email    string `schema:"email"`
		Password string `schema:"password"`
	}

	var userForm NewUserForm
	if err := parseForm(r, &userForm); err != nil {
		panic(err)
	}

	user := models.User{
		Name:     userForm.Name,
		Email:    userForm.Email,
		Password: userForm.Password,
	}

	if err := u.UserService.CreateUser(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err := u.signIn(w, &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/cookie", http.StatusFound)
}

func (u *Users) signIn(w http.ResponseWriter, user *models.User) error {
	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
		err = u.UserService.UpdateUser(user)
		if err != nil {
			return err
		}
	}

	cookie := http.Cookie{
		Name:     "remember_token",
		Value:    user.Remember,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	return nil
}

func (u *Users) CookieTest(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("remember_token")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user, err := u.UserService.GetUserByRemember(cookie.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, user)
}
