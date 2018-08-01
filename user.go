package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

const (
	// Superuser is above admin
	Superuser = "SUPERUSER"
	// Admin is below superuser and controls everything
	Admin = "ADMIN"
	// Manager is below admin and is managed by admins or superusers
	Manager = "MANAGER"
)

// User model
type User struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Username  string
	Password  string
	Role      string
}

// Create method for a User model
func (u *User) Create(username string, password string, role string) error {

	Log.WithFields(logrus.Fields{
		"username": username,
		"password": password,
		"role":     role,
	}).Debug("Creating user...")

	db, err := gorm.Open("sqlite3", dbName)
	if err != nil {
		return err
	}
	defer db.Close()

	hashedPassword := hashAndSalt([]byte(password))

	db.Create(&User{
		Username: username,
		Password: hashedPassword,
		Role:     role,
	})

	return nil
}

// Update method for User model
func (u *User) Update(user User) error {
	db, err := gorm.Open("sqlite3", dbName)
	if err != nil {
		return err
	}
	defer db.Close()

	db.Model(&user).Update(user)

	return nil
}

// Delete method for User model
func (u *User) Delete(user User) error {
	db, err := gorm.Open("sqlite3", dbName)
	if err != nil {
		return err
	}
	defer db.Close()

	db.Delete(&user)

	return nil
}

func getAllUsers() ([]User, error) {
	var users []User

	db, err := gorm.Open("sqlite3", dbName)
	if err != nil {
		return users, err
	}
	defer db.Close()

	db.Find(&users)

	return users, nil
	// c.JSON(200, gin.H{"data": users})
}

func getUser(id string) (*User, error) {
	var user *User

	db, err := gorm.Open("sqlite3", dbName)
	if err != nil {
		return user, nil
	}
	defer db.Close()

	db.Where("id = ?", id).First(&user)

	return user, nil
}

func findUserByUsername(username string) *User {
	db, err := gorm.Open("sqlite3", dbName)
	if err != nil {
		Log.WithFields(logrus.Fields{
			"username": username,
		}).Panic("Failed to connect to database")
	}
	defer db.Close()

	var user *User
	db.Where(&User{Username: username}).First(&user)

	Log.WithFields(logrus.Fields{
		"user": user,
	}).Debug("findUserByUsername")

	return user
}

// SetupRoutesUser Sets up routes for user model
func SetupRoutesUser(router *gin.Engine, uri string) *gin.RouterGroup {
	usersRoute := router.Group(uri)

	usersRoute.GET("", func(c *gin.Context) {
		allUsers, err := getAllUsers()
		if err != nil {
			c.JSON(500, gin.H{"data": nil})
		}
		c.JSON(200, gin.H{"data": allUsers})
	})
	usersRoute.GET("/:id", func(c *gin.Context) {
		id := c.Param("id")
		user, err := getUser(id)
		if err != nil {
			c.JSON(500, gin.H{"data": nil})
		}
		c.JSON(200, gin.H{"data": user})
	})

	return usersRoute
}

func hashAndSalt(pwd []byte) string {

	// Use GenerateFromPassword to hash & salt pwd.
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}

func comparePasswords(hashedPwd string, plainPwd []byte) bool {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}
