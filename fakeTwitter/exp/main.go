package main

import (
	"fakeTwitter/models"
	"fmt"
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

	us, err := models.NewUserService(psqlInfo)
	if err != nil {
		panic(err)
	}
	defer us.CloseDB()
	us.DestructiveReset()

	user := models.User{
		Name:     "Michael Scott",
		Email:    "michael@dundermifflin.com",
		Password: "bestboss",
	}
	err = us.CreateUser(&user)
	if err != nil {
		panic(err)
	}
	// Verify that the user has a Remember and RememberHash
	fmt.Printf("%+v\n", user)
	if user.Remember == "" {
		panic("Invalid remember token")
	}

	// Now verify that we can lookup a user with that remember
	// token
	user2, err := us.GetUserByRemember(user.Remember)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", *user2)
}
