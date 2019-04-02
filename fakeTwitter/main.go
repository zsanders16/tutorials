package main

import (
	"fakeTwitter/controllers"
	"fakeTwitter/models"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	port = ":3000"

	host     = "localhost"
	dbport   = 5432
	user     = "postgres"
	password = "password"
	dbname   = "fakeTwitter_dev"
)

func main() {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, dbport, user, password, dbname)

	userService, err := models.NewUserService(psqlInfo)
	if err != nil {
		panic(err)
	}
	defer userService.CloseDB()
	userService.AutoMigrate()

	staticC := controllers.NewStatic()
	userC := controllers.NewUsers(userService)

	r := mux.NewRouter()
	r.Handle("/", staticC.Home)
	r.Handle("/contact", staticC.Contact)

	r.HandleFunc("/cookie", userC.CookieTest)

	r.HandleFunc("/users/new", userC.New).Methods("GET")
	r.HandleFunc("/login", userC.UserLoginView).Methods("GET")
	r.HandleFunc("/login", userC.UserLogin).Methods("POST")
	r.HandleFunc("/users", userC.Create).Methods("POST")

	http.ListenAndServe(port, r)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
