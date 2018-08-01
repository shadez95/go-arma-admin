package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/ssh/terminal"
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

func createsuperuser() (string, string) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter Username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSuffix(username, "\n")

	if username == "" {
		fmt.Println()
		fmt.Println()
		fmt.Println("Username cannot be blank")
		os.Exit(1)
	}

	fmt.Print("Enter Password: ")
	bytePassword, _ := terminal.ReadPassword(int(syscall.Stdin))
	password := string(bytePassword)

	if len(password) < 6 {
		fmt.Println()
		fmt.Println()
		fmt.Println("Password must be at least 6 characters long")
		os.Exit(1)
	}

	fmt.Println()

	fmt.Print("Confirm Password: ")
	bytePasswordConfirm, _ := terminal.ReadPassword(int(syscall.Stdin))
	passwordConfirm := string(bytePasswordConfirm)

	if password == passwordConfirm {
		fmt.Println("")
		return username, password
	}

	return "", ""
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

	hashedPassword := hashAndSalt(password)

	Log.WithFields(logrus.Fields{
		"hashedPassword": hashedPassword,
	}).Debug("Hashed password")

	user := &User{
		Username: username,
		Password: hashedPassword,
		Role:     role,
	}

	Log.WithFields(logrus.Fields{
		"user.Username": user.Username,
		"user.Password": user.Password,
		"user.Role":     user.Role,
	}).Debug("New user information")

	db.Create(user)

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

func findUserByUsername(username string) User {
	db, err := gorm.Open("sqlite3", dbName)
	if err != nil {
		Log.WithFields(logrus.Fields{
			"username": username,
		}).Panic("Failed to connect to database")
	}
	defer db.Close()

	var user User
	Log.Info("About to run query...")
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

func hashAndSalt(pwd string) string {

	// Use GenerateFromPassword to hash & salt pwd.
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	bytePwd := []byte(pwd)
	hash, err := bcrypt.GenerateFromPassword(bytePwd, bcrypt.DefaultCost)
	if err != nil {
		Log.Error(err)
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}

func comparePasswords(hashedPwd string, plainPwd string) bool {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	bytePlainPwd := []byte(plainPwd)

	Log.WithFields(logrus.Fields{
		"byteHash": byteHash,
		"plainPwd": bytePlainPwd,
	}).Debug("Comparing hash password and plain password")
	err := bcrypt.CompareHashAndPassword(byteHash, bytePlainPwd)
	if err != nil {
		Log.Error(err)
		return false
	}

	return true
}
