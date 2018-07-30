package main

import (
	"time"

	"github.com/jinzhu/gorm"
)

const (
	// Superuser is above admin
	Superuser = "SUPERUSER"
	// Admin is below superuser and controls everything
	Admin = "ADMIN"
	// Manager is below admin and is managed by admins or superusers
	Manager = "MANAGER"
)

// User model
type User struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Username  string
	Password  string
	Role      string
}

// Create method for a User model
func (u *User) Create(username string, password string, role string) {
	db, err := gorm.Open("sqlite3", "db.sqlite3")
	if err != nil {
		panic("Failed to connect to database")
	}
	defer db.Close()

	db.Create(&User{
		Username: username,
		Password: password,
		Role:     role,
	})
}

// Update method for User model
func (u *User) Update(user User) {
	db, err := gorm.Open("sqlite3", "db.sqlite3")
	if err != nil {
		panic("Failed to connect to database")
	}
	defer db.Close()

	db.Model(&user).Update(user)
}

// Delete method for User model
func (u *User) Delete(user User) {
	db, err := gorm.Open("sqlite3", "db.sqlite3")
	if err != nil {
		panic("Failed to connect to database")
	}
	defer db.Close()

	db.Delete(&user)
}
