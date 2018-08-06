package main

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
)

func runServer(portPtr int) {
	router := gin.Default()
	setupRoutes(router)
	port := strings.Join([]string{":", strconv.Itoa(portPtr)}, "")
	// Listen and Server in 0.0.0.0:8080
	router.Run(port)
}

func openDB() *gorm.DB {
	var db *gorm.DB
	db, err := gorm.Open("sqlite3", dbName)
	if err != nil {
		Log.WithFields(logrus.Fields{
			"err": err,
		}).Panic("Failed to connect to database")
		return db
	}
	return db
}
