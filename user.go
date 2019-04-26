package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/ssh/terminal"
)

const (
	// Superuser is above admin
	Superuser = 0
	// Admin is below superuser and controls everything
	Admin = 1
	// Manager is below admin and is managed by admins or superusers
	Manager = 2
)

// User model
type User struct {
	CustomGormModel
	Username string `gorm:"unique;not null"`
	Password []byte
	Role     int
}

type userNoPassword struct {
	ID        int
	CreatedAt time.Time
	UpdatedAt time.Time
	Username  string
	Role      int
}

var user *User

// UserNoPassword is a helper method that removes password stored in DB
func (u *User) UserNoPassword() userNoPassword {
	return userNoPassword{
		ID:        u.ID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		Username:  u.Username,
		Role:      u.Role,
	}
}

func createsuperuser() error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter Username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSuffix(username, "\n")
	username = strings.TrimSuffix(username, "\r")

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
		fmt.Println()
		CreateUser(username, password, Superuser)
		return nil
	}

	return errors.New("Failed to create user. Passwords don't match")
}

// CreateUser creates a user model in database
func CreateUser(username string, password string, role int) error {
	db := openDB()
	defer db.Close()

	Log.WithFields(logrus.Fields{
		"username": username,
		"password": password,
		"role":     role,
	}).Debug("Creating user...")

	hashedPassword := hashAndSalt(password)

	user := &User{
		Username: username,
		Password: []byte(hashedPassword),
		Role:     role,
	}

	err := db.Create(user).Error
	if err != nil {
		return err
	}

	return nil
}

// Update method for User model
func (u *User) Update(user User) error {
	db := openDB()
	defer db.Close()

	db.Model(&user).Update(user)

	return nil
}

// Delete method for User model
func (u *User) Delete(user User) error {
	db := openDB()
	defer db.Close()

	db.Delete(&user)

	return nil
}

func getAllUsers() ([]userNoPassword, error) {
	db := openDB()
	defer db.Close()

	var users []userNoPassword
	// db.Find(&users)
	// db.Find(&users)
	err := db.Table("users").Select("id, username, role, createdAt, updatedAt").Find(&users).Error
	if err != nil {
		Log.Error(err)
	}

	return users, nil
}

func getUserByID(id int) (*User, error) {
	db := openDB()
	defer db.Close()

	var user User
	err := db.First(&user, id).RecordNotFound()
	Log.WithField("user", user).Debug("Getting user by id")
	Log.WithField("&user", &user).Debug("Getting user by id")
	if err {
		Log.Warn(err)
	}

	return &user, nil
}

func findUserByUsername(username string) (*User, error) {
	db := openDB()
	defer db.Close()

	var user User
	err := db.Table("users").Where(&User{Username: username}).First(&user).Error

	Log.WithFields(logrus.Fields{
		"user": user,
	}).Debug("findUserByUsername")

	if err != nil {
		return &user, err
	}

	return &user, nil
}

func findUserAuthenticate(username string) (User, error) {
	db := openDB()
	defer db.Close()

	var user User
	err := db.Where(&User{Username: username}).First(&user).Error

	Log.WithFields(logrus.Fields{
		"user": user,
	}).Debug("findUserAuthenticate")

	if err != nil {
		return user, err
	}

	return user, nil
}

func getSelf(c *gin.Context) (*User, error) {
	var err error
	jwtClaims := jwt.ExtractClaims(c)
	id := int(jwtClaims["id"].(float64))
	Log.WithField("id", id).Debug("Getting myself")
	user, err = getUserByID(id)
	if err != nil {
		Log.Warn(err)
		return user, err
	}
	return user, nil
}

// checkIfManager checks to see if user has role manager
func checkIfManager(c *gin.Context) {
	user, err := getSelf(c)
	if err != nil {
		c.JSON(500, gin.H{
			"data":    nil,
			"message": err,
		})
	}
	if user.Role == Manager {
		c.JSON(403, gin.H{
			"data":    nil,
			"message": "You are a manager and not authorized to request this information",
		})
	}
}

// userRoutes Sets up routes for user model
func userRoutes(router *gin.Engine, uri string) *gin.RouterGroup {

	usersRoute := router.Group(uri)

	usersRoute.GET("", func(c *gin.Context) {

		checkIfManager(c)
		var users []userNoPassword
		users, err := getAllUsers()
		if err != nil {
			c.JSON(500, gin.H{"data": nil})
		}

		c.JSON(200, gin.H{"data": users})

	})

	usersRoute.GET("/:id", func(c *gin.Context) {

		checkIfManager(c)
		id := c.Param("id")
		intID, err := strconv.Atoi(id)
		user, err := getUserByID(intID)
		if err != nil {
			c.JSON(500, gin.H{
				"data":    nil,
				"message": err,
			})
		}

		c.JSON(200, gin.H{
			"data":    user.UserNoPassword(),
			"message": "ok",
		})

	})

	router.GET("/me", func(c *gin.Context) {

		user, err := getSelf(c)
		if err != nil {
			c.JSON(500, gin.H{"data": err})
		}

		fmt.Println(user.UserNoPassword())
		c.JSON(200, gin.H{"data": user.UserNoPassword()})
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

	err := bcrypt.CompareHashAndPassword(byteHash, bytePlainPwd)
	if err != nil {
		Log.Error(err)
		return false
	}

	return true
}
