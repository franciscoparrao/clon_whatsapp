package handlers

import (
	"net/http"
	"time"
	
	"whatsapp-clone/config"
	
	"github.com/gin-gonic/gin"
)

// Global variables for handlers
var (
	hub         *Hub
	authHandler *AuthHandler
	chatHandler *ChatHandler
)

func init() {
	// Initialize hub
	hub = NewHub()
	go hub.Run()
	
	// Initialize handlers with the JWT secret from config
	cfg := config.LoadConfig()
	authHandler = NewAuthHandler(cfg.JWTSecret)
	chatHandler = &ChatHandler{hub: hub}
}

// HealthCheck handles health check requests
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now(),
		"service":   "whatsapp-clone-backend",
	})
}

// Register handles user registration
func Register(c *gin.Context) {
	authHandler.Register(c)
}

// Login handles user login
func Login(c *gin.Context) {
	authHandler.Login(c)
}

// GetChats handles getting user's chats
func GetChats(c *gin.Context) {
	chatHandler.GetChats(c)
}

// CreateChat handles creating a new chat
func CreateChat(c *gin.Context) {
	chatHandler.CreateChat(c)
}

// GetMessages handles getting messages from a chat
func GetMessages(c *gin.Context) {
	chatHandler.GetMessages(c)
}

// SendMessage handles sending a message to a chat
func SendMessage(c *gin.Context) {
	chatHandler.SendMessage(c)
}

// GetUsers handles getting all users
func GetUsers(c *gin.Context) {
	chatHandler.GetUsers(c)
}

// GetCurrentUser handles getting current user info
func GetCurrentUser(c *gin.Context) {
	authHandler.GetCurrentUser(c)
}

// WebSocketHandler handles WebSocket connections
func WebSocketHandler(c *gin.Context) {
	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to upgrade connection"})
		return
	}
	
	// Get user info from context
	userID := c.GetString("userID")
	username := c.GetString("username")
	
	if userID == "" {
		userID = "anonymous"
		username = "Anonymous User"
	}
	
	// Create and register client
	client := &Client{
		userID:   userID,
		username: username,
		hub:      hub,
		conn:     conn,
		send:     make(chan []byte, 256),
	}
	
	hub.register <- client
	
	// Start goroutines for reading and writing
	go client.readPump()
	go client.writePump()
}