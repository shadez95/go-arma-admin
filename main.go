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
	"github.com/sirupsen/logrus"
	"github.com/subosito/gotenv"
	"golang.org/x/crypto/ssh/terminal"
)

var dbName = "db.sqlite3"
var Log = logrus.New()

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

	if username == "" {
		fmt.Println()
		fmt.Println()
		fmt.Println("Username cannot be blank")
		os.Exit(1)
	}

	fmt.Print("Enter Password: ")
	bytePassword, _ := terminal.ReadPassword(int(syscall.Stdin))
	password := string(bytePassword)

	if len(password) <= 6 {
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

func main() {
	// Logger output to stdout
	Log.Out = os.Stdout

	// Setup port stuff
	port, err := strconv.Atoi(os.Getenv("APP_PORT"))
	if err != nil {
		port = 8000
	}

	// Setup flags
	superuserPtr := flag.Bool("createsuperuser", false, "Create a superuser")
	makemigrationsPtr := flag.Bool("makemigrations", false, "Make migrations")
	portPtr := flag.Int("port", port, "Set the port. Default: 8000")
	logleverPtr := flag.String("loglevel", "info", "Set logging level. Default: warning")
	flag.Parse()

	// Setup logging level
	switch *logleverPtr {
	case "debug":
		Log.Level = logrus.DebugLevel
	case "info":
		Log.Level = logrus.InfoLevel
	case "warn":
		Log.Level = logrus.WarnLevel
	case "error":
		Log.Level = logrus.ErrorLevel
	case "fatal":
		Log.Level = logrus.FatalLevel
	case "panic":
		Log.Level = logrus.PanicLevel
	default:
		Log.Level = logrus.InfoLevel
	}

	// Check for migrations flag
	if *makemigrationsPtr {

		fmt.Println("Making migrations if needed...")
		err := runMigrations()
		if err != nil {
			panic("Failed to make migrations")
		}
		fmt.Println("Successfully made migrations if needed")

		os.Exit(0)

	} else if *superuserPtr { // Check for createsuperuser flag

		username, password := createsuperuser()

		fmt.Println("Attempting to create user...")
		var user *User
		err := user.Create(username, hashAndSalt([]byte(password)), Superuser)
		if err != nil {
			panic("Failed to create user in database")
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
