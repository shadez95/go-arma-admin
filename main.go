package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/subosito/gotenv"
	"golang.org/x/crypto/ssh/terminal"
)

var DB = make(map[string]string)

func init() {
	gotenv.Load()
}

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	router := gin.Default()

	// Ping test
	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	// Get user value
	router.GET("/user/:name", func(c *gin.Context) {
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
	authorized := router.Group("/", gin.BasicAuth(gin.Accounts{
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

	return router
}

func createsuperuser() (string, string) {
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
	port, err := strconv.Atoi(os.Getenv("APP_PORT"))
	if err != nil {
		port = 8000
	}
	superuserPtr := flag.Bool("createsuperuser", false, "Create a superuser")
	portPtr := flag.Int("port", port, "Set the port. Default: 8000")
	flag.Parse()

	if !*superuserPtr {

		fmt.Println("No superuser to create")
		fmt.Println("Starting server...")

		router := setupRouter()
		SetupRoutesUser(router, "/users")

		port := strings.Join([]string{":", strconv.Itoa(*portPtr)}, "")
		// Listen and Server in 0.0.0.0:8080
		router.Run(port)

	} else {

		username, password := createsuperuser()

		if username == "" {
			fmt.Println()
			panic("Failed to create username and password")
		}

		// TODO: Create user in database when superuser is created

		fmt.Println(username, password)
	}

}
