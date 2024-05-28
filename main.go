package main

import (
	"github.com/gin-gonic/gin"
	"github.com/o1egl/paseto"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// load secret key
var secretKey []byte

// init function to load the secret key from the environment variable
func init() {
	key := os.Getenv("PASETO_SECRET_KEY")

	println("key: ", key)

	if len(key) != 32 {
		log.Fatal("PASETO_SECRET_KEY must be 32 bytes long")
	}

	if key == "" {
		log.Fatal("PASETO_SECRET_KEY environment variable is required")
	}

	secretKey = []byte(key)
}

// middleware to verify Paseto token
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			c.Abort()
			return
		}

		token := parts[1]

		var payload paseto.JSONToken
		var footer string

		err := paseto.NewV2().Decrypt(token, secretKey, &payload, &footer)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if time.Now().After(payload.Expiration) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has expired"})
			c.Abort()
			return
		}

		c.Set("username", payload.Subject)
		c.Next()
	}
}

func generateToken(username string) (string, error) {
	now := time.Now()
	exp := now.Add(time.Hour)

	jsonToken := paseto.JSONToken{
		Subject:    username,
		IssuedAt:   now,
		Expiration: exp,
	}

	token, err := paseto.NewV2().Encrypt(secretKey, jsonToken, "")

	if err != nil {
		log.Printf("Failed to generate token: %v", err)
		return "", err
	}

	return token, nil
}

func (server *Server) loginUser(c *gin.Context) {
	var loginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// implement your own logic here to retrieve the users
	if loginRequest.Username == "username123" && loginRequest.Password == "password123" {
		token, err := generateToken(loginRequest.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": token})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
	}
}

func (server *Server) privateRoute(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello from a private route"})
}

func (server *Server) setupRouter() {
	router := gin.Default()

	// public route
	router.POST("/login", server.loginUser)

	// apply middleware to secure routes
	protected := router.Group("/")
	protected.Use(authMiddleware())
	{
		protected.GET("/private", server.privateRoute)
	}

	server.router = router
}

type Server struct {
	router *gin.Engine
}

/*
 How to call this endpoint:

 - login
 curl -X POST http://localhost:8080/login -H "Content-Type: application/json" -d '{"username": "username123", "password": "password123"}'

 - private route
 curl -X GET http://localhost:8080/private -H "Authorization : Bearer <token from login response>"
*/

func main() {
	server := &Server{}
	server.setupRouter()

	port := ":2000"

	// run the server
	if err := server.router.Run(port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
