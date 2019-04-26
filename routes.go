package main

import (
	"net/http"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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

	// Going to probably need below in the future to verify origin sources
	// could be customizable through environment variable
	// Reference: https://github.com/olahol/melody#faq
	socketRouter.Upgrader.CheckOrigin = func(req *http.Request) bool {
		if Log.Level == logrus.DebugLevel {
			Log.Debug("CheckOrigin is true")
			return true
		}
		Log.Debug("CheckOrigin is false")
		return false
	}
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
