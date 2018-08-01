package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/subosito/gotenv"
	"golang.org/x/crypto/ssh/terminal"
)

var dbName = "db.sqlite3"

func init() {
	gotenv.Load()
}

func runMigrations() error {
	db, err := gorm.Open("sqlite3", dbName)
	if err != nil {
		return err
	}
	defer db.Close()

	db.AutoMigrate(&User{})

	return nil
}

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	router := gin.Default()

	// Ping test
	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	return router
}

func createsuperuser() (string, string) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter Username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSuffix(username, "\n")

	fmt.Print("Enter Password: ")
	bytePassword, _ := terminal.ReadPassword(int(syscall.Stdin))
	password := string(bytePassword)

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

func main() {
	port, err := strconv.Atoi(os.Getenv("APP_PORT"))
	if err != nil {
		port = 8000
	}
	superuserPtr := flag.Bool("createsuperuser", false, "Create a superuser")
	makemigrationsPtr := flag.Bool("makemigrations", false, "Make migrations")
	portPtr := flag.Int("port", port, "Set the port. Default: 8000")
	flag.Parse()

	if *makemigrationsPtr {

		fmt.Println("Making migrations if needed...")
		err := runMigrations()
		if err != nil {
			panic("Failed to make migrations")
		}
		fmt.Println("Successfully made migrations if needed")

		os.Exit(0)

	} else if *superuserPtr {

		username, password := createsuperuser()

		if username == "" {
			// fmt.Println()
			panic("Username cannot be blank")
			os.Exit(1)
		}

		fmt.Println("Attempting to create user...")
		var user *User
		err := user.Create(username, hashAndSalt([]byte(password)), Superuser)
		if err != nil {
			panic("Failed to create user in database")
			os.Exit(1)
		}
		fmt.Println("Superuser created successfully")
		os.Exit(0)

	} else {

		fmt.Println("Starting server...")

		router := setupRouter()
		setupAuth(router)

		router.Use(jwtMiddleware.MiddlewareFunc())
		{
			router.GET("/refreshToken", jwtMiddleware.RefreshHandler)
			SetupRoutesUser(router, "/users")
		}

		port := strings.Join([]string{":", strconv.Itoa(*portPtr)}, "")
		// Listen and Server in 0.0.0.0:8080
		router.Run(port)

	}
}
