package models

import (
	"errors"
	"fakeTwitter/hash"
	"fakeTwitter/rand"

	"golang.org/x/crypto/bcrypt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

const hmacSecretKey = "super-secrect-hmac-key"

var (
	//ErrNotFound returned when unable to find resource in the database
	ErrNotFound = errors.New("models: resouces not found")
	// ErrInvalidID is returned when a provide ID is not valid
	ErrInvalidID = errors.New("models: ID provided was invalid")
	// ErrInvalidPassword is returned when an invalid password is used when attempting to authenticate a user.
	ErrInvalidPassword = errors.New("models: incorrect password provided")

	userPwPepper = "super-secrect-pepper"
)

// User defines all fields related to a user
type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"not null;unique_index"`
	Password     string `gorm:"-"`
	PasswordHash string `gorm:"not nul"`
	Remember     string `gorm:"-"`
	RememberHash string `gorm:"not null;uniqe_index"`
}

// NewUserService opens a connection to the database and returns a UserService to witch appends the db connection
func NewUserService(psqlInfo string) (*UserService, error) {

	db, err := gorm.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	hmac := hash.NewHMAC(hmacSecretKey)

	return &UserService{
		db:   db,
		hmac: hmac,
	}, err
}

// UserService is the abstraction layer between the User and database
type UserService struct {
	db   *gorm.DB
	hmac hash.HMAC
}

// CreateUser creates a user in the database
func (us *UserService) CreateUser(user *User) error {

	hashedBytes, err := us.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes)
	user.Password = ""

	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
	}
	user.RememberHash = us.hmac.Hash(user.Remember)

	return us.db.Create(user).Error
}

// GetUserByID retrieves a user by ID
func (us *UserService) GetUserByID(id string) (*User, error) {
	var user User
	db := us.db.Where("id = ?", id)
	err := first(db, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByEmail retrieves a user by email
func (us *UserService) GetUserByEmail(email string) (*User, error) {
	var user User
	db := us.db.Where("email = ?", email)
	err := first(db, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByRemember retrieves a user by the remember token
func (us *UserService) GetUserByRemember(remember string) (*User, error) {
	var user User
	ht := us.hmac.Hash(remember)
	db := us.db.Where("remember_hash = ?", ht)
	err := first(db, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser will update a user
func (us *UserService) UpdateUser(user *User) error {
	if user.Remember != "" {
		user.RememberHash = us.hmac.Hash(user.Remember)
	}

	return us.db.Save(user).Error
}

// DeleteUser is remove a user from the database
func (us *UserService) DeleteUser(id uint) error {
	if id == 0 {
		return ErrInvalidID
	}
	user := &User{Model: gorm.Model{ID: id}}
	return us.db.Delete(user).Error
}

// CloseDB closes the db connection that it holds
func (us *UserService) CloseDB() {
	us.db.Close()
}

//  Helper methods

// DestructiveReset resets user table for developmemt purposes
func (us *UserService) DestructiveReset() error {
	err := us.db.DropTableIfExists(&User{}).Error
	if err != nil {
		return err
	}
	return us.AutoMigrate()
}

// AutoMigrate creates the user table in the database
func (us *UserService) AutoMigrate() error {
	if err := us.db.AutoMigrate(&User{}).Error; err != nil {
		return err
	}
	return nil
}

// HashPassword takes a User with a password and hashes it, setting the unhashed password to an empty string
func (us *UserService) HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password+userPwPepper), bcrypt.DefaultCost)
}

// AuthenticateUser verifies a useremail and password matches what is found in the database
func (us *UserService) AuthenticateUser(email, password string) (*User, error) {
	user := &User{
		Email: email,
	}
	u, err := us.GetUserByEmail(user.Email)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password+userPwPepper)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return nil, ErrInvalidPassword
		}
		return nil, err
	}

	return u, nil
}

func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}
	return err
}
