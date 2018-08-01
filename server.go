package main

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	router := gin.Default()

	// Ping test
	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	return router
}

func runServer(portPtr int) {
	router := setupRouter()
	setupAuth(router)

	router.Use(jwtMiddleware.MiddlewareFunc())
	{
		router.GET("/refreshToken", jwtMiddleware.RefreshHandler)
		SetupRoutesUser(router, "/users")
	}

	port := strings.Join([]string{":", strconv.Itoa(portPtr)}, "")
	// Listen and Server in 0.0.0.0:8080
	router.Run(port)
}
