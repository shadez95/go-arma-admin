package main

import "github.com/gin-gonic/gin"

func setupRoutes(router *gin.Engine) {
	// Ping test
	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	// Authenticate at this route
	router.POST("/login", jwtMiddleware.LoginHandler)

	// Auth required routes
	router.Use(jwtMiddleware.MiddlewareFunc())
	{
		router.GET("/refreshToken", jwtMiddleware.RefreshHandler)
		userRoutes(router, "/users")
	}
}
