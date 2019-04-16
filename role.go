package main

import (
	"time"
)

type Role struct {
	ID            int `gorm:"primary_key"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	CreateServers bool
	DeleteServers bool
	UpdateServers bool
	CreateUsers   bool
	UpdateUsers   bool
	DeleteUsers   bool
}
