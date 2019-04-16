package main

import (
	"sync"

	"gopkg.in/olahol/melody.v1"
)

type gopherInfo struct {
	ID int
}

func setupMRoutes(mRouter *melody.Melody) {
	gophers := make(map[*melody.Session]*gopherInfo)
	lock := new(sync.Mutex)
	counter := 0
	mRouter.HandleConnect(func(s *melody.Session) {
		lock.Lock()
		gophers[s] = &gopherInfo{ID: counter}
		counter++
		lock.Unlock()
	})

	mRouter.HandleDisconnect(func(s *melody.Session) {
		lock.Lock()
		delete(gophers, s)
		lock.Unlock()
	})

	mRouter.HandleMessage(func(s *melody.Session, msg []byte) {
		lock.Lock()
		// Anything happening to server gets broadcasted from here to all clients
		lock.Unlock()
	})
}
