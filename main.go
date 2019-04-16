package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	"github.com/subosito/gotenv"
)

var dbName = "db.sqlite3"

// Log is the main logger. Use this for logging
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

func main() {
	// Logger output to stdout
	Log.Out = os.Stdout

	// Setup flags
	logleverPtr := flag.String("loglevel", "info", "Set logging level. Options: debug, info, warn, error, fatal, panic")

	command, err := initCommands()
	if err != nil {
		fmt.Println()
		fmt.Println(err)
		os.Exit(1)
	}

	// Setup logging level
	configureLogger(*logleverPtr)

	// Check for migrations flag
	if command == MakeMigrationsCommand {

		Log.Info("Making migrations if needed...")
		err := runMigrations()
		if err != nil {
			Log.Panic("Failed to make migrations")
		}
		Log.Info("Successfully made migrations if needed")

		os.Exit(0)

	} else if command == CreateSuperUserCommand { // Check for createsuperuser flag

		username, password := createsuperuser()

		Log.Info("Attempting to create user...")

		var user *User
		err := user.Create(username, password, Superuser)
		if err != nil {
			Log.Panic("Failed to create user in database")
		}
		Log.Info("Superuser created successfully")
		os.Exit(0)

	} else if command == RunCommand {

		Log.Info("Starting server...")

		// *portPtr is a var declared in commandLine.go
		runServer(*portPtr)
	}
}
