package main

import (
	"fmt"

	"github.com/gofrs/uuid"
)

func generateUUID() (uuid.UUID, error) {
	// Creating UUID Version 4
	// panic on error
	// or error handling
	uuid, err := uuid.NewV4()
	if err != nil {
		fmt.Printf("Failed to generate UUID: %s", err)
		return uuid, err
	}

	Log.WithField("uuid", uuid).Debug("Generated UUID")

	return uuid, err
}
