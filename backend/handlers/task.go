package handlers

import (
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func TaskHandler(c *gin.Context) {
	workerID := c.Param("worker_id")
	log.Printf("Recibida petición de tareas para %s", workerID)

	targetURL := pythonApiBaseUrl + "/tareas/" + workerID
	resp, err := http.Get(targetURL)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "No se pudo conectar con servicio de tareas"})
		return
	}
	defer resp.Body.Close()

	c.Status(resp.StatusCode)
	c.Header("Content-Type", resp.Header.Get("Content-Type"))
	io.Copy(c.Writer, resp.Body)
}
