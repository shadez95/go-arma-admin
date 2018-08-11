package main

import (
	"time"

	"github.com/gin-gonic/gin"
)

type Servers struct {
	ID        int    `gorm:"primary_key"`
	Name      string `gorm:"unique;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func armaRoutes(router *gin.Engine, uri string) *gin.RouterGroup {
	armaRoute := router.Group(uri)

	armaRoute.GET("", func(c *gin.Context) {
		c.JSON(200, gin.H{"data": map[string]interface{}{"hello": "world"}, "message": "It worked"})
	})

	return armaRoute
}
