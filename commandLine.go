package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

type command string

func (s command) String() string {
	return string(s)
}

const (
	// RunCommand is a subcommand type that represents the command run
	RunCommand command = "run"
	// CreateSuperUserCommand is a subcommand type that represents the command createsuperuser
	CreateSuperUserCommand command = "createsuperuser"
	// MakeMigrationsCommand is a subcommand type that represents the command makemigrations
	MakeMigrationsCommand command = "makemigrations"
)

var portPtr *int

func initCommands() (command, error) {
	// Subcommands
	runCommand := flag.NewFlagSet(RunCommand.String(), flag.ExitOnError)
	createSuperUserCommand := flag.NewFlagSet(CreateSuperUserCommand.String(), flag.ExitOnError)
	makeMigrationsCommand := flag.NewFlagSet(MakeMigrationsCommand.String(), flag.ExitOnError)

	flag.Parse()

	// Attempt to get port from environment variable, if not default to 8000
	port, err := strconv.Atoi(os.Getenv("APP_PORT"))
	if err != nil {
		port = 8000
	}

	// Get SteamCMD from environment variable
	steamCMD := os.Getenv("STEAM_CMD")

	// Setup flags for subcommands
	portPtr = runCommand.Int("port", port, "Set the port.")
	steamCMDptr := runCommand.String("steamcmd", steamCMD, "Set the path to the SteamCMD binary\n")
	if len(*steamCMDptr) == 0 {
		log.Fatal("Need steamcmd path set. Set as environment variable or set value to option. Run -h to see all options.")
		os.Exit(1)
	}

	// Overwrite flag.Usage to print out pretty usage help
	flag.Usage = func() {
		fmt.Printf("Usage: go-arma-admin <command> [options]\n")
		fmt.Println()
		fmt.Println("\tList of commands include: run, createsuperuser, makemigrations")
		fmt.Println()
		fmt.Println("\tgo-arma-admin <command> -h : quick help on <command>")
		fmt.Println()

		// Lookup loglevel flag and display as well in flag usage
		loglevelPtr := flag.Lookup("loglevel")
		fmt.Printf("\t-%v\n\t    %v (Default: %v)\n", loglevelPtr.Name, loglevelPtr.Usage, loglevelPtr.DefValue)

		// flag.VisitAll(func(flag) {
		// 	flagPtr := flag.Lookup("loglevel")
		// 	fmt.Printf("\t-%v\n\t    %v (Default: %v)\n", loglevelPtr.Name, loglevelPtr.Usage, loglevelPtr.DefValue)
		// })
	}

	if len(os.Args) < 2 {
		flag.Usage()
		return "", errors.New("No command. A command is required")
	}

	var commandFlags = os.Args[2:]
	switch os.Args[1] {
	case RunCommand.String():
		runCommand.Usage = func() {
			fmt.Printf("Usage: go-arma-admin run\n")
			fmt.Println("\tPrompt user to enter username and password to create a superuser")
		}
		runCommand.Parse(commandFlags)
		if runCommand.Lookup("help") != nil || runCommand.Lookup("h") != nil {
			runCommand.Usage()
			os.Exit(1)
		}
		return RunCommand, nil
	case CreateSuperUserCommand.String():
		createSuperUserCommand.Parse(commandFlags)
		return CreateSuperUserCommand, nil
	case MakeMigrationsCommand.String():
		makeMigrationsCommand.Parse(commandFlags)
		return MakeMigrationsCommand, nil
	default:
		return "", errors.New("Not a command")
	}
}
