package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
)

var dbName = "db.sqlite3"

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
