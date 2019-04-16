package main

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
)

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
		// armaRoutes(router, "/servers")

		// mRouter is declared in mRoutes.go
		mRouter := melody.New()
		setupMRoutes(mRouter)
		router.GET("/servers/ws", func(c *gin.Context) {
			mRouter.HandleRequest(c.Writer, c.Request)
		})
	}
}
