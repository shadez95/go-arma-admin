package main

import (
	"os"
	"strings"
	"time"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// DO NOT CHANGE
// Reference: https://github.com/appleboy/gin-jwt/issues/170
var identityKey = "id"

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// the jwt middleware
var jwtMiddleware = &jwt.GinJWTMiddleware{
	Realm: "armaadmin",
	// store this somewhere, if your server restarts and you're
	// generating random passwords, any valid JWTs will be invalid
	Key:             []byte(os.Getenv("APP_SECRET")),
	Timeout:         time.Hour,
	MaxRefresh:      time.Hour * 24,
	IdentityKey:     identityKey,
	IdentityHandler: idHandler,
	Authenticator:   authenticate,
	// this method allows you to jump in and set user information
	// JWTs aren't encrypted, so don't store any sensitive info
	PayloadFunc:      payload,
	SigningAlgorithm: "HS256",
}

func idHandler(c *gin.Context) interface{} {
	Log.Debug("IdentityHandler executing")
	claims := jwt.ExtractClaims(c)
	var user User
	user.ID = int(claims["id"].(float64))
	return &user
}

func authenticate(c *gin.Context) (interface{}, error) {
	var loginVals login
	err := c.ShouldBind(&loginVals)
	if err != nil {
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
	pwdMatch := comparePasswords(string(user.Password), password)

	if username == user.Username && pwdMatch {
		Log.Debug("Passwords matched and returning username")
		Log.Debug(user)
		return user, nil
	}

	Log.Debug("Passwords do not match")
	return nil, jwt.ErrFailedAuthentication
}

func payload(data interface{}) jwt.MapClaims {
	Log.Debug("Deploying payload")
	if v, ok := data.(User); ok {
		return jwt.MapClaims{
			"username":  v.Username,
			"role":      v.Role,
			identityKey: v.ID,
		}
	}
	Log.Warn("User data returned is corrupted")
	return jwt.MapClaims{"error": true}
}
