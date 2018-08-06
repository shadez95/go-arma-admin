package main

import (
	"os"
	"strings"
	"time"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// the jwt middleware
var jwtMiddleware = jwt.GinJWTMiddleware{
	Realm: "armaadmin",
	// store this somewhere, if your server restarts and you're
	// generating random passwords, any valid JWTs will be invalid
	Key:           []byte(os.Getenv("APP_SECRET")),
	Timeout:       time.Hour,
	MaxRefresh:    time.Hour * 24,
	Authenticator: authenticate,
	Unauthorized: func(c *gin.Context, code int, message string) {
		c.JSON(code, gin.H{
			"message": message,
			"data":    nil,
		})
	},
	// this method allows you to jump in and set user information
	// JWTs aren't encrypted, so don't store any sensitive info
	PayloadFunc: payload,
}

func authenticate(userID string, password string, c *gin.Context) (string, bool) {
	// it goes without saying that you'd be going to some form
	// of persisted storage, rather than doing this

	Log.WithFields(logrus.Fields{
		"userID":   userID,
		"password": strings.Repeat("x", len(password)),
	}).Debug("Authenticating user...")

	user, err := findUserByUsername(userID)
	if err != nil {
		Log.Debug("User was not found")
		return "", false
	}
	Log.WithFields(logrus.Fields{
		"user":          user,
		"user.Username": user.Username,
	}).Debug("User that was retrieved")
	pwdMatch := comparePasswords(user.Password, password)

	if userID == user.Username && pwdMatch {
		Log.Debug("Passwords matched and returning userID")
		return userID, true
	}

	Log.Debug("Passwords do not match")
	return "", false
}

func payload(userID string) map[string]interface{} {
	// in this method, you'd want to fetch some user info
	// based on their email address (which is provided once
	// they've successfully logged in).  the information
	// you set here will be available the lifetime of the
	// user's sesion
	user, _ := findUserByUsername(userID)
	return map[string]interface{}{
		"id":       user.ID,
		"username": user.Username,
		"role":     user.Role,
	}
}
