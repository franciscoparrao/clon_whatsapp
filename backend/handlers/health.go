package handlers

import (
	"net/http"
	"runtime"
	"time"
	"whatsapp-clone/database"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	startTime time.Time
}

type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Uptime    string            `json:"uptime"`
	Version   string            `json:"version"`
	Services  map[string]string `json:"services"`
	System    SystemInfo        `json:"system"`
}

type SystemInfo struct {
	GoVersion    string `json:"goVersion"`
	NumGoroutine int    `json:"numGoroutine"`
	NumCPU       int    `json:"numCPU"`
	MemoryUsage  MemoryStats `json:"memoryUsage"`
}

type MemoryStats struct {
	Alloc      uint64 `json:"alloc"`      // bytes allocated and in use
	TotalAlloc uint64 `json:"totalAlloc"` // bytes allocated (even if freed)
	Sys        uint64 `json:"sys"`        // bytes obtained from system
	NumGC      uint32 `json:"numGC"`      // number of garbage collections
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{
		startTime: time.Now(),
	}
}

// Health returns the health status of the application
func (h *HealthHandler) Health(c *gin.Context) {
	// Get memory statistics
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Check service dependencies
	services := h.checkServices()

	// Determine overall health status
	status := "healthy"
	for _, serviceStatus := range services {
		if serviceStatus != "up" {
			status = "degraded"
			break
		}
	}

	response := HealthResponse{
		Status:    status,
		Timestamp: time.Now(),
		Uptime:    time.Since(h.startTime).String(),
		Version:   "1.0.0", // TODO: Get from build info
		Services:  services,
		System: SystemInfo{
			GoVersion:    runtime.Version(),
			NumGoroutine: runtime.NumGoroutine(),
			NumCPU:       runtime.NumCPU(),
			MemoryUsage: MemoryStats{
				Alloc:      m.Alloc,
				TotalAlloc: m.TotalAlloc,
				Sys:        m.Sys,
				NumGC:      m.NumGC,
			},
		},
	}

	if status != "healthy" {
		c.JSON(http.StatusServiceUnavailable, response)
	} else {
		c.JSON(http.StatusOK, response)
	}
}

// Liveness is a simple check to verify the service is alive
func (h *HealthHandler) Liveness(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "alive",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// Readiness checks if the service is ready to accept traffic
func (h *HealthHandler) Readiness(c *gin.Context) {
	// Check critical dependencies
	services := h.checkServices()
	
	ready := true
	for service, status := range services {
		// Only database and cache are critical for readiness
		if (service == "database" || service == "cache") && status != "up" {
			ready = false
			break
		}
	}

	if !ready {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "not_ready",
			"services": services,
			"timestamp": time.Now().Format(time.RFC3339),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
		"services": services,
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// checkServices checks the status of dependent services
func (h *HealthHandler) checkServices() map[string]string {
	services := make(map[string]string)

	// TODO: Implement actual health checks for these services
	// For now, return mock statuses
	
	// Check database connection
	services["database"] = h.checkDatabase()
	
	// Check Redis/cache connection
	services["cache"] = h.checkCache()
	
	// Check WebSocket service
	services["websocket"] = h.checkWebSocket()
	
	// Check file storage service
	services["storage"] = h.checkStorage()

	return services
}

func (h *HealthHandler) checkDatabase() string {
	// Verificar conexión a la base de datos
	if err := database.HealthCheck(); err != nil {
		return "down"
	}
	return "up"
}

func (h *HealthHandler) checkCache() string {
	// TODO: Implement actual Redis health check
	// Example: ping Redis connection
	return "up"
}

func (h *HealthHandler) checkWebSocket() string {
	// TODO: Check if WebSocket hub is running
	return "up"
}

func (h *HealthHandler) checkStorage() string {
	// TODO: Check if file storage is accessible
	// Example: try to write/read a test file
	return "up"
}