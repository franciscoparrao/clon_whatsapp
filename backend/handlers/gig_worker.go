package handlers

import (
	"net/http"
	"strconv"
	"time"
	"whatsapp-clone/database"
	"whatsapp-clone/models"
	"whatsapp-clone/repositories"

	"github.com/gin-gonic/gin"
)

// GigWorkerHandler maneja las solicitudes relacionadas con gig workers
type GigWorkerHandler struct {
	repo *repositories.GigWorkerRepository
}

// NewGigWorkerHandler crea una nueva instancia del handler
func NewGigWorkerHandler() *GigWorkerHandler {
	return &GigWorkerHandler{
		repo: repositories.NewGigWorkerRepository(database.GetDB()),
	}
}

// RegisterAsGigWorker registra a un usuario como gig worker
func (h *GigWorkerHandler) RegisterAsGigWorker(c *gin.Context) {
	userID := c.GetString("userID") // Obtenido del middleware de autenticación
	
	var req struct {
		WorkerCode string `json:"workerCode" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Código de trabajador requerido"})
		return
	}
	
	// Verificar si el usuario ya es un gig worker
	existingWorker, _ := h.repo.GetGigWorkerByUserID(userID)
	if existingWorker != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "El usuario ya está registrado como trabajador"})
		return
	}
	
	worker := &models.GigWorker{
		UserID:              userID,
		WorkerCode:          req.WorkerCode,
		Status:              "active",
		Rating:              5.0,
		TotalTasksCompleted: 0,
	}
	
	if err := h.repo.CreateGigWorker(worker); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al registrar trabajador"})
		return
	}
	
	c.JSON(http.StatusCreated, worker)
}

// GetWorkerProfile obtiene el perfil de un gig worker
func (h *GigWorkerHandler) GetWorkerProfile(c *gin.Context) {
	workerID := c.Param("id")
	
	worker, err := h.repo.GetGigWorkerByID(workerID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Trabajador no encontrado"})
		return
	}
	
	c.JSON(http.StatusOK, worker)
}

// GetMyWorkerProfile obtiene el perfil del trabajador actual
func (h *GigWorkerHandler) GetMyWorkerProfile(c *gin.Context) {
	userID := c.GetString("userID")
	
	worker, err := h.repo.GetGigWorkerByUserID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No está registrado como trabajador"})
		return
	}
	
	c.JSON(http.StatusOK, worker)
}

// UpdateAvailability actualiza la disponibilidad de un trabajador
func (h *GigWorkerHandler) UpdateAvailability(c *gin.Context) {
	userID := c.GetString("userID")
	
	// Obtener el worker ID del usuario
	worker, err := h.repo.GetGigWorkerByUserID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No está registrado como trabajador"})
		return
	}
	
	var req struct {
		Date      string `json:"date" binding:"required"`      // Formato YYYY-MM-DD
		StartTime string `json:"startTime" binding:"required"` // Formato HH:MM
		EndTime   string `json:"endTime" binding:"required"`   // Formato HH:MM
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}
	
	// Parsear fecha
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de fecha inválido"})
		return
	}
	
	availability := &models.WorkerAvailability{
		WorkerID:    worker.ID,
		Date:        date,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		IsConfirmed: false,
	}
	
	if err := h.repo.CreateAvailability(availability); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar disponibilidad"})
		return
	}
	
	c.JSON(http.StatusCreated, availability)
}

// ConfirmAvailability confirma la disponibilidad para un día
func (h *GigWorkerHandler) ConfirmAvailability(c *gin.Context) {
	availabilityID := c.Param("id")
	userID := c.GetString("userID")
	
	// Verificar que el usuario es el dueño de esta disponibilidad
	// (Esta verificación se podría mejorar con una consulta JOIN)
	
	if err := h.repo.ConfirmAvailability(availabilityID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al confirmar disponibilidad"})
		return
	}
	
	// Crear notificación
	notification := &models.Notification{
		UserID:  userID,
		Type:    "system",
		Title:   "Disponibilidad confirmada",
		Message: "Tu disponibilidad ha sido confirmada exitosamente",
	}
	h.repo.CreateNotification(notification)
	
	c.JSON(http.StatusOK, gin.H{"message": "Disponibilidad confirmada"})
}

// GetMyAvailability obtiene la disponibilidad del trabajador actual
func (h *GigWorkerHandler) GetMyAvailability(c *gin.Context) {
	userID := c.GetString("userID")
	dateStr := c.Query("date") // Formato YYYY-MM-DD
	
	worker, err := h.repo.GetGigWorkerByUserID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No está registrado como trabajador"})
		return
	}
	
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de fecha inválido"})
		return
	}
	
	availabilities, err := h.repo.GetAvailabilitiesByWorkerAndDate(worker.ID, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener disponibilidad"})
		return
	}
	
	c.JSON(http.StatusOK, availabilities)
}

// GetAvailableTasks obtiene las tareas disponibles para asignar
func (h *GigWorkerHandler) GetAvailableTasks(c *gin.Context) {
	tasks, err := h.repo.GetPendingTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener tareas"})
		return
	}
	
	c.JSON(http.StatusOK, tasks)
}

// AcceptTask acepta una tarea asignada
func (h *GigWorkerHandler) AcceptTask(c *gin.Context) {
	taskID := c.Param("id")
	userID := c.GetString("userID")
	
	worker, err := h.repo.GetGigWorkerByUserID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No está registrado como trabajador"})
		return
	}
	
	// Verificar que el trabajador puede aceptar tareas
	if !worker.CanAcceptTask() {
		c.JSON(http.StatusForbidden, gin.H{"error": "No puede aceptar tareas en este momento"})
		return
	}
	
	// Crear asignación
	assignment := &models.TaskAssignment{
		TaskID:   taskID,
		WorkerID: worker.ID,
		Status:   "assigned",
	}
	
	if err := h.repo.AssignTask(assignment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al asignar tarea"})
		return
	}
	
	// Actualizar estado a aceptado
	now := time.Now()
	assignment.AcceptedAt = &now
	assignment.Status = "accepted"
	h.repo.UpdateTaskAssignment(assignment)
	
	// Crear notificación
	notification := &models.Notification{
		UserID:  userID,
		Type:    "task_assigned",
		Title:   "Tarea asignada",
		Message: "Se te ha asignado una nueva tarea",
	}
	h.repo.CreateNotification(notification)
	
	c.JSON(http.StatusOK, assignment)
}

// StartTask marca el inicio de una tarea
func (h *GigWorkerHandler) StartTask(c *gin.Context) {
	assignmentID := c.Param("id")
	_ = c.GetString("userID") // TODO: Verificar que el usuario es el trabajador asignado
	
	assignment := &models.TaskAssignment{
		ID:     assignmentID,
		Status: "in_progress",
	}
	now := time.Now()
	assignment.StartedAt = &now
	
	if err := h.repo.UpdateTaskAssignment(assignment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al iniciar tarea"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Tarea iniciada"})
}

// CompleteTask marca una tarea como completada
func (h *GigWorkerHandler) CompleteTask(c *gin.Context) {
	assignmentID := c.Param("id")
	
	var req struct {
		Comments string `json:"comments"`
	}
	c.ShouldBindJSON(&req)
	
	assignment := &models.TaskAssignment{
		ID:       assignmentID,
		Status:   "completed",
		Comments: req.Comments,
	}
	now := time.Now()
	assignment.CompletedAt = &now
	
	if err := h.repo.UpdateTaskAssignment(assignment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al completar tarea"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Tarea completada"})
}

// GetMyAssignments obtiene las asignaciones del trabajador actual
func (h *GigWorkerHandler) GetMyAssignments(c *gin.Context) {
	userID := c.GetString("userID")
	status := c.Query("status") // Filtro opcional por estado
	
	worker, err := h.repo.GetGigWorkerByUserID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No está registrado como trabajador"})
		return
	}
	
	assignments, err := h.repo.GetWorkerAssignments(worker.ID, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener asignaciones"})
		return
	}
	
	c.JSON(http.StatusOK, assignments)
}

// GetAttendanceHistory obtiene el historial de asistencia
func (h *GigWorkerHandler) GetAttendanceHistory(c *gin.Context) {
	userID := c.GetString("userID")
	limitStr := c.DefaultQuery("limit", "30")
	
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 30
	}
	
	worker, err := h.repo.GetGigWorkerByUserID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No está registrado como trabajador"})
		return
	}
	
	history, err := h.repo.GetAttendanceHistory(worker.ID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener historial"})
		return
	}
	
	// Calcular estadísticas
	attendanceRate := models.CalculateAttendanceRate(history)
	confirmationRate := models.CalculateConfirmationRate(history)
	
	c.JSON(http.StatusOK, gin.H{
		"history":          history,
		"attendanceRate":   attendanceRate,
		"confirmationRate": confirmationRate,
	})
}

// GetNotifications obtiene las notificaciones del usuario
func (h *GigWorkerHandler) GetNotifications(c *gin.Context) {
	userID := c.GetString("userID")
	
	notifications, err := h.repo.GetUnreadNotifications(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener notificaciones"})
		return
	}
	
	c.JSON(http.StatusOK, notifications)
}

// MarkNotificationAsRead marca una notificación como leída
func (h *GigWorkerHandler) MarkNotificationAsRead(c *gin.Context) {
	notificationID := c.Param("id")
	
	if err := h.repo.MarkNotificationAsRead(notificationID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al marcar notificación"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Notificación marcada como leída"})
}

// CreateTask crea una nueva tarea (para administradores)
func (h *GigWorkerHandler) CreateTask(c *gin.Context) {
	var task models.Task
	
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}
	
	task.Status = "pending"
	
	if err := h.repo.CreateTask(&task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear tarea"})
		return
	}
	
	c.JSON(http.StatusCreated, task)
}

// RecordAttendance registra la asistencia de un trabajador
func (h *GigWorkerHandler) RecordAttendance(c *gin.Context) {
	var req struct {
		WorkerID      string `json:"workerId" binding:"required"`
		ScheduledDate string `json:"scheduledDate" binding:"required"` // YYYY-MM-DD
		ScheduledTime string `json:"scheduledTime" binding:"required"` // HH:MM
		DidConfirm    bool   `json:"didConfirm"`
		DidAttend     bool   `json:"didAttend"`
		AbsenceReason string `json:"absenceReason"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}
	
	date, err := time.Parse("2006-01-02", req.ScheduledDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de fecha inválido"})
		return
	}
	
	attendance := &models.AttendanceHistory{
		WorkerID:      req.WorkerID,
		ScheduledDate: date,
		ScheduledTime: req.ScheduledTime,
		DidConfirm:    req.DidConfirm,
		DidAttend:     req.DidAttend,
		AbsenceReason: req.AbsenceReason,
	}
	
	if err := h.repo.RecordAttendance(attendance); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al registrar asistencia"})
		return
	}
	
	c.JSON(http.StatusCreated, attendance)
}