package main

import (
	"encoding/base64"
	"net/http"
	"time"

	ginJWT "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

// Ticket is a struct contains an id and expire time
type Ticket struct {
	id     uuid.UUID
	expire time.Time
	userID int
}

// Tickets contains all tickets that are valid
var Tickets map[uuid.UUID]Ticket

// TicketResponse is used to standardize responses to clients
type TicketResponse struct {
	Ok      bool   `json:"ok"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Key     string `json:"key"`
	Expire  int64  `json:"expire"`
}

func wsTicketHandler(c *gin.Context) {
	// Get user ID from context
	claims := ginJWT.ExtractClaims(c)
	userID := int(claims["id"].(float64))

	id, err := generateUUID()
	if err != nil {
		Log.Warning(err)
		c.JSON(http.StatusInternalServerError, &TicketResponse{
			Ok:      false,
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Key:     "",
			Expire:  0,
		})
	}

	if Tickets == nil {
		Log.Debug("Initializing Tickets global variable")
		Tickets = make(map[uuid.UUID]Ticket)
	}

	// Tickets[id] = make()
	// Tickets will only last for 2 minutes
	expire := time.Now().Add(time.Minute * 2)
	Tickets[id] = Ticket{
		id:     id,
		expire: expire,
		userID: userID,
	}

	Log.WithField("id.String()", id.String()).Debug("API key for the websocket connection")

	// Encode key to base64
	key := encodeKey(id)
	Log.WithField("key", key).Debug("Encoded key")

	c.JSON(http.StatusOK, &TicketResponse{
		Ok:      true,
		Code:    http.StatusOK,
		Message: "Success",
		Key:     key,
		Expire:  expire.Unix(),
	})
}

func encodeKey(key uuid.UUID) string {
	return base64.StdEncoding.EncodeToString(key.Bytes())
}

func decodeKey(key string) []byte {
	decodeKey, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		Log.Error(err)
		return []byte("")
	}
	Log.WithField("decodedKey", decodeKey).Debug("Decoded key in bytes")
	return decodeKey
}

func convertKeyToUUID(key string) uuid.UUID {
	id, err := uuid.FromBytes(decodeKey(key))
	if err != nil {
		Log.Error(err)
	}
	Log.WithField("id", id).Debug("uuid from bytes")
	return id
}
