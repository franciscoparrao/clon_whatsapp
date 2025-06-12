package main

import (
	"log"
	"net/http"
	"whatsapp-clone/config"
	"whatsapp-clone/handlers"
	"whatsapp-clone/middleware"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

func main() {
	cfg := config.LoadConfig()
	
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			// Permitir cualquier origen (desarrollo)
			return true
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * 3600,
	}))

	// Add custom middleware to handle ngrok headers
	router.Use(func(c *gin.Context) {
		// Log incoming request for debugging
		log.Printf("Request: %s %s from %s", c.Request.Method, c.Request.URL.Path, c.Request.Header.Get("Origin"))
		
		// Handle ngrok headers
		if c.Request.Header.Get("X-Forwarded-Proto") != "" {
			c.Request.URL.Scheme = c.Request.Header.Get("X-Forwarded-Proto")
		}
		
		c.Next()
	})
	
	// Create logger middleware
	logger := log.New(log.Writer(), "[WHATSAPP] ", log.LstdFlags)
	loggerMiddleware := middleware.NewLoggerMiddleware(logger, middleware.INFO)
	
	router.Use(func(c *gin.Context) {
		loggerMiddleware.Log(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c.Next()
		})).ServeHTTP(c.Writer, c.Request)
	})

	api := router.Group("/api")
	{
		api.GET("/health", handlers.HealthCheck)
		api.POST("/auth/register", handlers.Register)
		api.POST("/auth/login", handlers.Login)
		
		protected := api.Group("/")
		protected.Use(middleware.GinAuthMiddleware(cfg.JWTSecret))
		{
			protected.GET("/user", handlers.GetCurrentUser)
			protected.GET("/users", handlers.GetUsers)
			protected.GET("/chats", handlers.GetChats)
			protected.POST("/chats", handlers.CreateChat)
			protected.GET("/chats/:id/messages", handlers.GetMessages)
			protected.POST("/chats/:id/messages", handlers.SendMessage)
		}
	}

	router.GET("/ws", handlers.WebSocketHandler)

	log.Printf("Server starting on port %s", cfg.ServerPort)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}