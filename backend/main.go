package main

import (
	"log"
	"net/http"
	"whatsapp-clone/config"
	"whatsapp-clone/database"
	"whatsapp-clone/handlers"
	"whatsapp-clone/middleware"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

func main() {
	cfg := config.LoadConfig()
	
	// Conectar a la base de datos
	dbConfig := database.Config{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		DBName:   cfg.DBName,
		SSLMode:  cfg.DBSSLMode,
	}
	
	if err := database.Connect(dbConfig); err != nil {
		log.Fatal("Error conectando a la base de datos:", err)
	}
	defer database.Close()
	
	// Inicializar schema si es necesario
	if err := database.InitializeSchema(); err != nil {
		log.Printf("Error inicializando schema: %v", err)
	}
	
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

	// Inicializar handlers
	gigWorkerHandler := handlers.NewGigWorkerHandler()
	
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
			
			// Rutas de Gig Workers
			protected.POST("/workers/register", gigWorkerHandler.RegisterAsGigWorker)
			protected.GET("/workers/profile", gigWorkerHandler.GetMyWorkerProfile)
			protected.GET("/workers/:id", gigWorkerHandler.GetWorkerProfile)
			
			// Disponibilidad
			protected.POST("/workers/availability", gigWorkerHandler.UpdateAvailability)
			protected.PUT("/workers/availability/:id/confirm", gigWorkerHandler.ConfirmAvailability)
			protected.GET("/workers/availability", gigWorkerHandler.GetMyAvailability)
			
			// Tareas
			protected.GET("/tasks", gigWorkerHandler.GetAvailableTasks)
			protected.POST("/tasks", gigWorkerHandler.CreateTask) // Para admins
			protected.POST("/tasks/:id/accept", gigWorkerHandler.AcceptTask)
			protected.PUT("/assignments/:id/start", gigWorkerHandler.StartTask)
			protected.PUT("/assignments/:id/complete", gigWorkerHandler.CompleteTask)
			protected.GET("/workers/assignments", gigWorkerHandler.GetMyAssignments)
			
			// Historial y estadísticas
			protected.GET("/workers/attendance", gigWorkerHandler.GetAttendanceHistory)
			protected.POST("/workers/attendance", gigWorkerHandler.RecordAttendance) // Para admins
			
			// Notificaciones
			protected.GET("/notifications", gigWorkerHandler.GetNotifications)
			protected.PUT("/notifications/:id/read", gigWorkerHandler.MarkNotificationAsRead)
		}
	}

	router.GET("/ws", handlers.WebSocketHandler)

	log.Printf("Server starting on port %s", cfg.ServerPort)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}