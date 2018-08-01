package main

import (
	"os"
	"time"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
)

var jwtMiddleware jwt.GinJWTMiddleware

func helloHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	c.JSON(200, gin.H{
		"userID": claims["id"],
		"text":   "Hello World.",
	})
}

func setupAuth(r *gin.Engine) {
	// the jwt middleware
	jwtMiddleware = jwt.GinJWTMiddleware{
		Realm: "armaadmin",
		// store this somewhere, if your server restarts and you're
		// generating random passwords, any valid JWTs will be invalid
		Key:           []byte(os.Getenv("APP_SECRET")),
		Timeout:       time.Hour,
		MaxRefresh:    time.Hour * 24,
		Authenticator: authenticate,
		// this method allows you to jump in and set user information
		// JWTs aren't encrypted, so don't store any sensitive info
		PayloadFunc: payload,
	}

	r.POST("/login", jwtMiddleware.LoginHandler)
}

func authenticate(userID string, password string, c *gin.Context) (string, bool) {
	// it goes without saying that you'd be going to some form
	// of persisted storage, rather than doing this

	user := findUserByUsername(userID)
	pwdMatch := comparePasswords(user.Password, []byte(password))

	if userID == user.Username && pwdMatch {
		return userID, true
	}

	return "", false
}

func payload(userID string) map[string]interface{} {
	// in this method, you'd want to fetch some user info
	// based on their email address (which is provided once
	// they've successfully logged in).  the information
	// you set here will be available the lifetime of the
	// user's sesion
	return map[string]interface{}{
		"id":   "1349",
		"role": "admin",
	}
}
