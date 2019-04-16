package main

import (
	"strconv"
	"time"

	"github.com/shadez95/go-arma"

	"github.com/gin-gonic/gin"
)

type Server struct {
	ID        int    `gorm:"primary_key"`
	Name      string `gorm:"unique;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	arma.Server
}

func createServer(s *Server) error {
	db := openDB()
	defer db.Close()

	err := db.Create(s).Error
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) Update() error {
	db := openDB()
	defer db.Close()

	err := db.Model(s).Update(s).Error
	if err != nil {
		return err
	}

	return nil
}

func armaRoutes(router *gin.Engine, uri string) *gin.RouterGroup {
	armaRoute := router.Group(uri)

	armaRoute.GET("", func(c *gin.Context) {
		c.JSON(200, gin.H{"data": map[string]interface{}{"hello": "world"}, "message": "It worked"})
	})
	armaRoute.GET("/:id", getServer)
	armaRoute.POST("", newServer)

	return armaRoute
}

func getServer(c *gin.Context) {
	id := c.Param("id")
	intID, err := strconv.Atoi(id)

}

func newServer(c *gin.Context) {

}
