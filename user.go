package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
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
func (u *User) Create(username string, password string, role string) {
	db, err := gorm.Open("sqlite3", "db.sqlite3")
	if err != nil {
		panic("Failed to connect to database")
	}
	defer db.Close()

	hashedPassword := hashAndSalt([]byte(password))

	db.Create(&User{
		Username: username,
		Password: hashedPassword,
		Role:     role,
	})
}

// Update method for User model
func (u *User) Update(user User) {
	db, err := gorm.Open("sqlite3", "db.sqlite3")
	if err != nil {
		panic("Failed to connect to database")
	}
	defer db.Close()

	db.Model(&user).Update(user)
}

// Delete method for User model
func (u *User) Delete(user User) {
	db, err := gorm.Open("sqlite3", "db.sqlite3")
	if err != nil {
		panic("Failed to connect to database")
	}
	defer db.Close()

	db.Delete(&user)
}

func getAllUsers() []User {
	db, err := gorm.Open("sqlite3", "db.sqlite3")
	if err != nil {
		panic("Failed to connect to database")
	}
	defer db.Close()

	var users []User
	db.Find(&users)

	return users
	// c.JSON(200, gin.H{"data": users})
}

func getUser(id string) *User {
	db, err := gorm.Open("sqlite3", "db.sqlite3")
	if err != nil {
		panic("Failed to connect to database")
	}
	defer db.Close()

	var user *User
	db.Where("id = ?", id).First(&user)

	return user
}

func findUserByUsername(username string) *User {
	db, err := gorm.Open("sqlite3", "db.sqlite3")
	if err != nil {
		panic("Failed to connect to database")
	}
	defer db.Close()

	var user *User
	db.Where(&User{Username: username}).First(&user)

	return user
}

// SetupRoutesUser Sets up routes for user model
func SetupRoutesUser(router *gin.Engine, uri string) *gin.RouterGroup {
	usersRoute := router.Group(uri)

	usersRoute.GET("", func(c *gin.Context) {
		allUsers := getAllUsers()
		c.JSON(200, gin.H{"data": allUsers})
	})
	usersRoute.GET("/:id", func(c *gin.Context) {
		id := c.Param("id")
		user := getUser(id)
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
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MaxCost)
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
