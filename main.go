package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"syscall"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/ssh/terminal"
)

var DB = make(map[string]string)

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	// Get user value
	r.GET("/user/:name", func(c *gin.Context) {
		user := c.Params.ByName("name")
		value, ok := DB[user]
		if ok {
			c.JSON(200, gin.H{"user": user, "value": value})
		} else {
			c.JSON(200, gin.H{"user": user, "status": "no value"})
		}
	})

	// Authorized group (uses gin.BasicAuth() middleware)
	// Same than:
	// authorized := r.Group("/")
	// authorized.Use(gin.BasicAuth(gin.Credentials{
	//	  "foo":  "bar",
	//	  "manu": "123",
	//}))
	authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
		"foo":  "bar", // user:foo password:bar
		"manu": "123", // user:manu password:123
	}))

	authorized.POST("admin", func(c *gin.Context) {
		user := c.MustGet(gin.AuthUserKey).(string)

		// Parse JSON
		var json struct {
			Value string `json:"value" binding:"required"`
		}

		if c.Bind(&json) == nil {
			DB[user] = json.Value
			c.JSON(200, gin.H{"status": "ok"})
		}
	})

	return r
}

func credentials() (string, string) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter Username: ")
	username, _ := reader.ReadString('\n')

	fmt.Print("Enter Password: ")
	bytePassword, _ := terminal.ReadPassword(int(syscall.Stdin))
	password := string(bytePassword)

	fmt.Println()

	fmt.Print("Confirm Password: ")
	bytePasswordConfirm, _ := terminal.ReadPassword(int(syscall.Stdin))
	passwordConfirm := string(bytePasswordConfirm)

	if password == passwordConfirm {
		return username, password
	}

	return "", ""
}

func main() {

	superuserPtr := flag.Bool("createsuperuser", false, "Create a superuser")
	flag.Parse()

	if !*superuserPtr {

		fmt.Println("No superuser to create")
		fmt.Println("Starting server...")

		r := setupRouter()
		// Listen and Server in 0.0.0.0:8080
		r.Run(":8080")

	} else {

		username, password := credentials()

		if username == "" {
			fmt.Println()
			panic("Failed to create username and password")
		}

		// TODO: Create user in database when superuser is created

		fmt.Println(username, password)
	}

}
