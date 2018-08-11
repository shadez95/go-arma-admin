package main

import (
	"flag"
	"log"
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

func init() {
	gotenv.Load()
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

	// Setup port stuff
	port, err := strconv.Atoi(os.Getenv("APP_PORT"))
	if err != nil {
		port = 8000
	}

	// SteamCMD
	steamCMD := os.Getenv("STEAM_CMD")

	// Setup flags
	superuserPtr := flag.Bool("createsuperuser", false, "Create a superuser\n")
	makemigrationsPtr := flag.Bool("makemigrations", false, "Make migrations\n")
	portPtr := flag.Int("port", port, "Set the port. Default: 8000\n")
	logleverPtr := flag.String("loglevel", "info", "Set logging level.\nOptions: debug, info, warn, error, fatal, panic\n")
	steamCMDptr := flag.String("steamcmd", steamCMD, "Set the path to the SteamCMD binary\n")
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

		Log.Info("Attempting to create super user...")
		err = createsuperuser()
		if err != nil {
			Log.Panic(err)
		}
		Log.Info("Superuser created successfully")
		os.Exit(0)

	} else {
		if len(*steamCMDptr) == 0 {
			log.Fatal("Need steamcmd path set. Set as environment variable or set value to option. Run -h to see all options.")
			os.Exit(1)
		}
		Log.Info("Starting server...")

		runServer(*portPtr)
	}
}
