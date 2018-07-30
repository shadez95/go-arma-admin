package main

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/jinzhu/gorm"
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

	db.Create(&User{
		Username: username,
		Password: password,
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

func getAllUsers(c *gin.Context) {
	db, err := gorm.Open("sqlite3", "db.sqlite3")
	if err != nil {
		panic("Failed to connect to database")
	}
	defer db.Close()

	var users []User
	db.Find(&users)

	c.JSON(200, gin.H{"data": users})
}

func getUser(c *gin.Context) {
	db, err := gorm.Open("sqlite3", "db.sqlite3")
	if err != nil {
		panic("Failed to connect to database")
	}
	defer db.Close()

	var user User
	db.Where("id = ?", c.Param("id")).First(&user)

	c.JSON(200, gin.H{"data": user})
}

// SetupRoutesUser Sets up routes for user model
func SetupRoutesUser(router *gin.Engine, uri string) *gin.RouterGroup {
	usersRoute := router.Group(uri)

	usersRoute.GET("/", getAllUsers)
	usersRoute.GET("/:id", getUser)

	return usersRoute
}
