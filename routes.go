package main

import (
	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
)

func setupRoutes(router *gin.Engine) {

	// Authenticate at this route
	authMiddleware, err := jwt.New(jwtMiddleware)
	if err != nil {
		Log.Fatalf("JWT Error: %v\n", err)
	}
	router.POST("/login", authMiddleware.LoginHandler)

	// mRouter is declared in mRoutes.go
	socketRouter := melody.New()
	setupWebsocketRoute(socketRouter)

	// Websocket server
	router.GET("/servers/ws", func(c *gin.Context) {
		socketRouter.HandleRequest(c.Writer, c.Request)
	})

	// Auth required routes
	router.Use(jwtMiddleware.MiddlewareFunc())
	{
		// Ticket system for connecting to websocket server
		router.GET("/wsTicket", wsTicketHandler)
		router.GET("/refreshToken", jwtMiddleware.RefreshHandler)
		userRoutes(router, "/users")
	}
}
