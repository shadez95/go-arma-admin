package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/gofrs/uuid"

	"gopkg.in/olahol/melody.v1"
)

type gopherInfo struct {
	ID int
}

func isAuthorized(id uuid.UUID) bool {
	// Check if ticket is expired
	// Tickets are stored in server memory
	// for the life of the running instance
	now := time.Now()
	isExpired := Tickets[id].expire.Before(now)
	if isExpired {
		// Delete ticket from memory if ticket is expired
		Log.WithField("Tickets[id]", Tickets[id]).Debug("Deleting Ticket")
		delete(Tickets, id)
		return false
	}

	Log.Debug("Ticket is authorized")
	return true
}

func handleAuth(s *melody.Session) {
	// Get key from URL parameter
	key := s.Request.URL.Query().Get("key")

	// convertKeyToUUID is found in wsTicket.go
	id := convertKeyToUUID(key)

	// If not authorized, close with message
	if !isAuthorized(id) {
		Log.Warn("Unauthorized access")
		s.CloseWithMsg([]byte("Unauthorized access"))
	}
}

// Setup mRoutes/websocket
func setupWebsocketRoute(mRouter *melody.Melody) {
	gophers := make(map[*melody.Session]*gopherInfo)
	lock := new(sync.Mutex)
	counter := 0

	mRouter.HandleConnect(func(s *melody.Session) {
		lock.Lock()

		// Check if authorized
		handleAuth(s)

		gophers[s] = &gopherInfo{ID: counter}
		counter++
		lock.Unlock()
	})

	mRouter.HandleDisconnect(func(s *melody.Session) {
		lock.Lock()

		// Get key from URL parameter
		key := s.Request.URL.Query().Get("key")

		// convertKeyToUUID is found in wsTicket.go
		id := convertKeyToUUID(key)

		// Delete ticket stored in Tickets if it exists.
		// This is here so if user disconnects, and not the server
		// forcing a disconnect, this will force the user to get
		// another ticket.
		if _, exists := Tickets[id]; exists {
			delete(Tickets, id)
		}

		// Delete session stored in gophers
		delete(gophers, s)
		lock.Unlock()
	})

	mRouter.HandleMessage(func(s *melody.Session, msg []byte) {
		lock.Lock()

		// Check if authorized
		handleAuth(s)

		// Anything happening to server gets broadcasted from here to all clients
		output := fmt.Sprintf("Server message: %s", string(msg[:]))
		mRouter.Broadcast([]byte(output))
		lock.Unlock()
	})
}
