package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
)

var portPtr *int

// Helper function to initialize new subcommands
func newSubcommand(name string, usage string, description string) *flag.FlagSet {
	cmd := flag.NewFlagSet(name, flag.ExitOnError)
	cmd.Usage = func() {
		fmt.Printf("Usage: %v\n", usage)
		fmt.Printf("\tDescription: %v\n", description)
		loglevelPtr := cmd.Lookup("loglevel")
		fmt.Printf("\t-%v\n\t    %v (Default: %v)\n", loglevelPtr.Name, loglevelPtr.Usage, loglevelPtr.DefValue)
	}

	return cmd
}

// Initialize all subcommands and return which command is ran and possible error
func initCommands() (string, error) {

	// Initialize subcommands
	runSubcommmand := newSubcommand("run", "go-arma-admin run", "Starts go-arma-admin server. SteamCMD must be in PATH environment variable.")
	csuSubcommand := newSubcommand("createsuperuser", "go-arma-admin createsuperuser", "Prompt user to enter username and password to create a superuser")
	makemigrationsSubcommand := newSubcommand("makemigrations", "go-arma-admin makemigrations", "Make migrations on current databased configured in settings")

	// Attempt to get port from environment variable, if not default to 8000
	port, err := strconv.Atoi(os.Getenv("APP_PORT"))
	if err != nil {
		port = 8000
	}

	// Setup flags for subcommands
	portPtr = runSubcommmand.Int("port", port, "Set the port.")

	// Overwrite default flag.Usage when executing no subcommands to print out usage help
	flag.Usage = func() {
		fmt.Printf("Usage: go-arma-admin <command> [options]\n")
		fmt.Println()
		fmt.Println("\tList of commands include: run, createsuperuser, makemigrations")
		fmt.Println()
		fmt.Println("\tgo-arma-admin <command> -h : quick help on <command>")
		fmt.Println()
	}

	// If no subcommands, print usage and error
	if len(os.Args) < 2 {
		flag.Usage()
		return "", errors.New("No command. A command is required")
	}

	// Figure out which subcommand was executed and run the binary based on command
	switch os.Args[1] {
	// If run command was executed
	case runSubcommmand.Name():
		commandWrapper(runSubcommmand)
		return "run", nil

	// If createsuperuser command was executed
	case csuSubcommand.Name():
		commandWrapper(csuSubcommand)
		return "createsuperuser", nil

	// If makemigrations command was executed
	case makemigrationsSubcommand.Name():
		commandWrapper(makemigrationsSubcommand)
		return "makemigrations", nil

	// If a command was executed that isn't a recognized command
	default:
		return "", errors.New("Not a recognized command")
	}
}

func commandWrapper(cmd *flag.FlagSet) {
	logLevelPtr := cmd.String("loglevel", "info", "Set logging level. Options: debug, info, warn, error, fatal, panic")
	cmd.Parse(os.Args[2:])
	configureLogger(*logLevelPtr)
	if cmd.Lookup("help") != nil || cmd.Lookup("h") != nil {
		cmd.Usage()
		os.Exit(1)
	}
}
