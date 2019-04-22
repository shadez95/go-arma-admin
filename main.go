package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	"github.com/subosito/gotenv"
)

// Log is the main logger. Use this for logging
var Log = logrus.New()

// AppSecret is the variable that is used for encryption
// and is set through an environment variable
var AppSecret string

func init() {
	gotenv.Load()
	AppSecret = os.Getenv("APP_SECRET")
}

func runMigrations() error {
	db := openDB()
	defer db.Close()

	db.AutoMigrate(&User{})

	return nil
}

func runServer(portPtr int) {
	router := gin.Default()
	setupRoutes(router)
	port := strings.Join([]string{":", strconv.Itoa(portPtr)}, "")
	// Listen and Server in 0.0.0.0:8080
	router.Run(port)
}

func main() {
	// Logger output to stdout
	Log.Out = os.Stdout

	subcommand, err := initCommands()
	if err != nil {
		fmt.Println()
		fmt.Println(err)
		os.Exit(1)
	}

	// Display log level output
	fmt.Println("Log level output: ", Log.GetLevel())

	// Check for migrations flag
	if subcommand == "makemigrations" {

		Log.Info("Making migrations if needed...")
		err := runMigrations()
		if err != nil {
			Log.Panic("Failed to make migrations")
		}
		Log.Info("Successfully made migrations if needed")

		os.Exit(0)

	} else if subcommand == "createsuperuser" { // Check for createsuperuser flag

		Log.Info("Attempting to create super user...")
		err = createsuperuser()
		if err != nil {
			Log.Panic(err)
		}
		Log.Info("Superuser created successfully")
		os.Exit(0)

	} else if subcommand == "run" {
		Log.Info("Starting server...")

		// *portPtr is a var declared in commandLine.go
		runServer(*portPtr)
	}
}
