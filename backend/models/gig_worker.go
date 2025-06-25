package models

import (
	"time"
)

// GigWorker representa un trabajador en el sistema
type GigWorker struct {
	ID                  string    `json:"id" db:"id"`
	UserID              string    `json:"userId" db:"user_id"`
	WorkerCode          string    `json:"workerCode" db:"worker_code"`
	Status              string    `json:"status" db:"status"` // active, inactive, suspended
	Rating              float64   `json:"rating" db:"rating"`
	TotalTasksCompleted int       `json:"totalTasksCompleted" db:"total_tasks_completed"`
	CreatedAt           time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt           time.Time `json:"updatedAt" db:"updated_at"`
	
	// Relaciones
	User *User `json:"user,omitempty"`
}

// WorkerAvailability representa la disponibilidad de un trabajador
type WorkerAvailability struct {
	ID          string    `json:"id" db:"id"`
	WorkerID    string    `json:"workerId" db:"worker_id"`
	Date        time.Time `json:"date" db:"date"`
	StartTime   string    `json:"startTime" db:"start_time"` // Formato HH:MM
	EndTime     string    `json:"endTime" db:"end_time"`     // Formato HH:MM
	IsConfirmed bool      `json:"isConfirmed" db:"is_confirmed"`
	ConfirmedAt *time.Time `json:"confirmedAt,omitempty" db:"confirmed_at"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	
	// Relaciones
	Worker *GigWorker `json:"worker,omitempty"`
}

// Task representa una tarea o turno
type Task struct {
	ID                string    `json:"id" db:"id"`
	Title             string    `json:"title" db:"title"`
	Description       string    `json:"description,omitempty" db:"description"`
	Location          string    `json:"location,omitempty" db:"location"`
	ScheduledDate     time.Time `json:"scheduledDate" db:"scheduled_date"`
	ScheduledStartTime string   `json:"scheduledStartTime" db:"scheduled_start_time"` // Formato HH:MM
	ScheduledEndTime   string   `json:"scheduledEndTime" db:"scheduled_end_time"`     // Formato HH:MM
	Status            string    `json:"status" db:"status"`           // pending, assigned, in_progress, completed, cancelled
	Priority          string    `json:"priority" db:"priority"`       // low, normal, high, urgent
	CreatedAt         time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt         time.Time `json:"updatedAt" db:"updated_at"`
	
	// Relaciones
	Assignments []TaskAssignment `json:"assignments,omitempty"`
}

// TaskAssignment representa la asignación de una tarea a un trabajador
type TaskAssignment struct {
	ID          string     `json:"id" db:"id"`
	TaskID      string     `json:"taskId" db:"task_id"`
	WorkerID    string     `json:"workerId" db:"worker_id"`
	AssignedAt  time.Time  `json:"assignedAt" db:"assigned_at"`
	AcceptedAt  *time.Time `json:"acceptedAt,omitempty" db:"accepted_at"`
	StartedAt   *time.Time `json:"startedAt,omitempty" db:"started_at"`
	CompletedAt *time.Time `json:"completedAt,omitempty" db:"completed_at"`
	Status      string     `json:"status" db:"status"` // assigned, accepted, rejected, in_progress, completed
	Rating      *int       `json:"rating,omitempty" db:"rating"`
	Comments    string     `json:"comments,omitempty" db:"comments"`
	
	// Relaciones
	Task   *Task      `json:"task,omitempty"`
	Worker *GigWorker `json:"worker,omitempty"`
}

// AttendanceHistory representa el historial de asistencia para predicciones
type AttendanceHistory struct {
	ID            string    `json:"id" db:"id"`
	WorkerID      string    `json:"workerId" db:"worker_id"`
	ScheduledDate time.Time `json:"scheduledDate" db:"scheduled_date"`
	ScheduledTime string    `json:"scheduledTime" db:"scheduled_time"` // Formato HH:MM
	DidConfirm    bool      `json:"didConfirm" db:"did_confirm"`
	DidAttend     bool      `json:"didAttend" db:"did_attend"`
	AbsenceReason string    `json:"absenceReason,omitempty" db:"absence_reason"`
	CreatedAt     time.Time `json:"createdAt" db:"created_at"`
	
	// Relaciones
	Worker *GigWorker `json:"worker,omitempty"`
}

// Notification representa una notificación para un usuario
type Notification struct {
	ID        string     `json:"id" db:"id"`
	UserID    string     `json:"userId" db:"user_id"`
	Type      string     `json:"type" db:"type"`           // task_assigned, reminder, system
	Title     string     `json:"title" db:"title"`
	Message   string     `json:"message" db:"message"`
	IsRead    bool       `json:"isRead" db:"is_read"`
	ReadAt    *time.Time `json:"readAt,omitempty" db:"read_at"`
	CreatedAt time.Time  `json:"createdAt" db:"created_at"`
	
	// Relaciones
	User *User `json:"user,omitempty"`
}

// Métodos de utilidad

// IsAvailable verifica si un trabajador está disponible en un momento específico
func (w *GigWorker) IsAvailable(date time.Time, startTime, endTime string) bool {
	return w.Status == "active"
}

// CanAcceptTask verifica si un trabajador puede aceptar una tarea
func (w *GigWorker) CanAcceptTask() bool {
	return w.Status == "active" && w.Rating >= 3.0
}

// IsOverdue verifica si una tarea está vencida
func (t *Task) IsOverdue() bool {
	if t.Status == "completed" || t.Status == "cancelled" {
		return false
	}
	
	now := time.Now()
	scheduledDateTime := t.ScheduledDate
	// Aquí podrías parsear la hora y compararla más precisamente
	return now.After(scheduledDateTime)
}

// CanBeAssigned verifica si una tarea puede ser asignada
func (t *Task) CanBeAssigned() bool {
	return t.Status == "pending" && !t.IsOverdue()
}

// IsActive verifica si una asignación está activa
func (ta *TaskAssignment) IsActive() bool {
	return ta.Status == "in_progress" || ta.Status == "accepted"
}

// AttendanceRate calcula la tasa de asistencia del trabajador
func CalculateAttendanceRate(history []AttendanceHistory) float64 {
	if len(history) == 0 {
		return 0.0
	}
	
	attended := 0
	for _, h := range history {
		if h.DidAttend {
			attended++
		}
	}
	
	return float64(attended) / float64(len(history)) * 100
}

// ConfirmationRate calcula la tasa de confirmación del trabajador
func CalculateConfirmationRate(history []AttendanceHistory) float64 {
	if len(history) == 0 {
		return 0.0
	}
	
	confirmed := 0
	for _, h := range history {
		if h.DidConfirm {
			confirmed++
		}
	}
	
	return float64(confirmed) / float64(len(history)) * 100
}