package main

import (
	"time"
)

// Role type specifies a role details
type Role struct {
	ID            int       `gorm:"primary_key"`
	CreatedAt     time.Time `gorm:"column:createdAt"`
	UpdatedAt     time.Time `gorm:"column:updatedAt"`
	CreateServers bool
	DeleteServers bool
	UpdateServers bool
	CreateUsers   bool
	UpdateUsers   bool
	DeleteUsers   bool
}
