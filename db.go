package main

import (
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
)

func getHome() string {
	homePath, err := homedir.Dir()
	if err != nil {
		Log.Panic(err)
	}

	return homePath
}

func openDB() *gorm.DB {
	var db *gorm.DB

	if homePath == "" {
		homePath = getHome()
	}

	pathArr := []string{homePath, ".arma-admin/arma_admin_db.sqlite"}
	dbPath := strings.Join(pathArr, "/")
	Log.WithField("dbPath", dbPath).Debug("DB Path")
	db, err := gorm.Open("sqlite3", dbPath)
	if err != nil {
		Log.WithFields(logrus.Fields{
			"err": err,
		}).Panic("Failed to connect to database")
		return db
	}
	return db
}
