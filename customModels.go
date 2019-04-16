package main

import "time"

// CustomGormModel is struct for custom Gorm model
type CustomGormModel struct {
	ID        int       `gorm:"primary_key"`
	CreatedAt time.Time `gorm:"column:createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt"`
}
