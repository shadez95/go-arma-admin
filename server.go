package main

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func runServer(portPtr int) {
	router := gin.Default()
	setupRoutes(router)
	port := strings.Join([]string{":", strconv.Itoa(portPtr)}, "")
	// Listen and Server in 0.0.0.0:8080
	router.Run(port)
}
