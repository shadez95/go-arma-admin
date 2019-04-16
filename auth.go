package main

import (
	"os"
	"strings"
	"time"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type userData map[string]interface{}

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

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

func authenticate(c *gin.Context) (interface{}, error) {
	var loginVals login
	if err := c.ShouldBind(&loginVals); err != nil {
		return "", jwt.ErrMissingLoginValues
	}

	username := loginVals.Username
	password := loginVals.Password

	Log.WithFields(logrus.Fields{
		"username": username,
		"password": strings.Repeat("x", len(password)),
	}).Debug("Authenticating user...")

	user, err := findUserAuthenticate(username)
	if err != nil {
		Log.Debug("User was not found")
		return nil, jwt.ErrFailedAuthentication
	}
	Log.WithFields(logrus.Fields{
		"user":          user,
		"user.Username": user.Username,
	}).Debug("User that was retrieved")
	pwdMatch := comparePasswords(user.Password, password)

	if username == user.Username && pwdMatch {
		Log.Debug("Passwords matched and returning username")
		return &user, nil
	}

	Log.Debug("Passwords do not match")
	return nil, jwt.ErrFailedAuthentication
}

func payload(data interface{}) jwt.MapClaims {
	if v, ok := data.(*userNoPassword); ok {
		return jwt.MapClaims{
			"id":       v.ID,
			"username": v.Username,
			"rold":     v.Role,
		}
	}
	return jwt.MapClaims{}
}
