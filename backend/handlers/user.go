package handlers

import (
	"bytes"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

const pythonApiBaseUrl = "http://localhost:8000"

func ScheduleHandler(c *gin.Context) {
	workerID := c.Param("worker_id")
	log.Printf("Recibida petición de horario para %s", workerID)

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No se pudo leer el cuerpo de la petición"})
		return
	}
	defer c.Request.Body.Close()

	targetURL := pythonApiBaseUrl + "/horarios/" + workerID
	req, err := http.NewRequest(http.MethodPost, targetURL, bytes.NewBuffer(body))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear petición"})
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "No se pudo conectar con servicio de horarios"})
		return
	}
	defer resp.Body.Close()

	c.Status(resp.StatusCode)
	io.Copy(c.Writer, resp.Body)
}
