package main

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shadez95/go-arma"
)

const (
	// Active is a constant string that is used by Status
	Active = "active"

	// Inactive is a constant string that is used by Status
	Inactive = "inactive"
	// Starting is a constant string that is used by Status
	Starting = "starting"
)

// Status type represents status of a server
type Status string

// Server type represents a server
type Server struct {
	ID        int       `gorm:"primary_key"`
	CreatedAt time.Time `gorm:"column:createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt"`
	Name      string    `gorm:"unique;not null"`
	Status    Status
	arma.Server
}

var server *Server

func createServer(s *Server) error {
	db := openDB()
	defer db.Close()

	err := db.Create(s).Error
	if err != nil {
		return err
	}

	return nil
}

func updateServer(s *Server) error {
	db := openDB()
	defer db.Close()

	err := db.Model(&s).Update(s).Error
	if err != nil {
		return err
	}

	return nil
}

func deleteServer(s *Server) error {
	db := openDB()
	defer db.Close()

	err := db.Delete(&s).Error
	if err != nil {
		return err
	}
	return nil
}

func findServerByID(id int) (*Server, error) {
	db := openDB()
	defer db.Close()

	err := db.Where("id = ?", id).First(&server).Error
	if err != nil {
		return nil, err
	}

	return server, nil
}

func findServerByName(name string) (*Server, error) {
	db := openDB()
	defer db.Close()

	err := db.Where("name = ?", name).First(&server).Error
	if err != nil {
		return nil, err
	}
	return server, nil
}

func getServer(c *gin.Context) (*Server, error) {
	id := c.Param("id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	db := openDB()
	defer db.Close()

	db.Where("id = ?", intID).First(&server)
	if err != nil {
		return nil, err
	}

	return server, nil
}
