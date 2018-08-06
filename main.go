package main

import (
	"flag"
	"os"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	"github.com/subosito/gotenv"
)

// Log is the main logger. Use this for logging
var Log = logrus.New()

func init() {
	gotenv.Load()
}

func runMigrations() error {
	db := openDB()
	defer db.Close()

	db.AutoMigrate(&User{})

	return nil
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
	superuserPtr := flag.Bool("createsuperuser", false, "Create a superuser\n")
	makemigrationsPtr := flag.Bool("makemigrations", false, "Make migrations\n")
	portPtr := flag.Int("port", port, "Set the port. Default: 8000\n")
	logleverPtr := flag.String("loglevel", "info", "Set logging level.\nOptions: debug, info, warn, error, fatal, panic\n")
	flag.Parse()

	// Setup logging level
	configureLogger(*logleverPtr)

	// Check for migrations flag
	if *makemigrationsPtr {

		Log.Info("Making migrations if needed...")
		err := runMigrations()
		if err != nil {
			Log.Panic("Failed to make migrations")
		}
		Log.Info("Successfully made migrations if needed")

		os.Exit(0)

	} else if *superuserPtr { // Check for createsuperuser flag

		username, password := createsuperuser()

		Log.Info("Attempting to create user...")

		var user *User
		err := user.Create(username, password, Superuser)
		if err != nil {
			Log.Panic("Failed to create user in database")
		}
		Log.Info("Superuser created successfully")
		os.Exit(0)

	} else {

		Log.Info("Starting server...")

		runServer(*portPtr)
	}
}
