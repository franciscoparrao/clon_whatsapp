package repositories

import (
	"database/sql"
	"fmt"
	"time"
	"whatsapp-clone/models"
)

// GigWorkerRepository maneja las operaciones de base de datos para gig workers
type GigWorkerRepository struct {
	db *sql.DB
}

// NewGigWorkerRepository crea una nueva instancia del repositorio
func NewGigWorkerRepository(db *sql.DB) *GigWorkerRepository {
	return &GigWorkerRepository{db: db}
}

// CreateGigWorker crea un nuevo gig worker
func (r *GigWorkerRepository) CreateGigWorker(worker *models.GigWorker) error {
	query := `
		INSERT INTO gig_workers (user_id, worker_code, status, rating, total_tasks_completed)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`
	
	err := r.db.QueryRow(
		query,
		worker.UserID,
		worker.WorkerCode,
		worker.Status,
		worker.Rating,
		worker.TotalTasksCompleted,
	).Scan(&worker.ID, &worker.CreatedAt, &worker.UpdatedAt)
	
	return err
}

// GetGigWorkerByID obtiene un gig worker por su ID
func (r *GigWorkerRepository) GetGigWorkerByID(id string) (*models.GigWorker, error) {
	worker := &models.GigWorker{}
	query := `
		SELECT id, user_id, worker_code, status, rating, total_tasks_completed, created_at, updated_at
		FROM gig_workers
		WHERE id = $1
	`
	
	err := r.db.QueryRow(query, id).Scan(
		&worker.ID,
		&worker.UserID,
		&worker.WorkerCode,
		&worker.Status,
		&worker.Rating,
		&worker.TotalTasksCompleted,
		&worker.CreatedAt,
		&worker.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("gig worker not found")
	}
	
	return worker, err
}

// GetGigWorkerByUserID obtiene un gig worker por el ID de usuario
func (r *GigWorkerRepository) GetGigWorkerByUserID(userID string) (*models.GigWorker, error) {
	worker := &models.GigWorker{}
	query := `
		SELECT id, user_id, worker_code, status, rating, total_tasks_completed, created_at, updated_at
		FROM gig_workers
		WHERE user_id = $1
	`
	
	err := r.db.QueryRow(query, userID).Scan(
		&worker.ID,
		&worker.UserID,
		&worker.WorkerCode,
		&worker.Status,
		&worker.Rating,
		&worker.TotalTasksCompleted,
		&worker.CreatedAt,
		&worker.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("gig worker not found")
	}
	
	return worker, err
}

// UpdateGigWorker actualiza un gig worker
func (r *GigWorkerRepository) UpdateGigWorker(worker *models.GigWorker) error {
	query := `
		UPDATE gig_workers
		SET status = $2, rating = $3, total_tasks_completed = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING updated_at
	`
	
	err := r.db.QueryRow(
		query,
		worker.ID,
		worker.Status,
		worker.Rating,
		worker.TotalTasksCompleted,
	).Scan(&worker.UpdatedAt)
	
	return err
}

// CreateAvailability crea una nueva disponibilidad
func (r *GigWorkerRepository) CreateAvailability(availability *models.WorkerAvailability) error {
	query := `
		INSERT INTO worker_availability (worker_id, date, start_time, end_time, is_confirmed, confirmed_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`
	
	err := r.db.QueryRow(
		query,
		availability.WorkerID,
		availability.Date,
		availability.StartTime,
		availability.EndTime,
		availability.IsConfirmed,
		availability.ConfirmedAt,
	).Scan(&availability.ID, &availability.CreatedAt)
	
	return err
}

// GetAvailabilitiesByWorkerAndDate obtiene las disponibilidades de un trabajador para una fecha
func (r *GigWorkerRepository) GetAvailabilitiesByWorkerAndDate(workerID string, date time.Time) ([]models.WorkerAvailability, error) {
	query := `
		SELECT id, worker_id, date, start_time, end_time, is_confirmed, confirmed_at, created_at
		FROM worker_availability
		WHERE worker_id = $1 AND date = $2
		ORDER BY start_time
	`
	
	rows, err := r.db.Query(query, workerID, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var availabilities []models.WorkerAvailability
	for rows.Next() {
		var a models.WorkerAvailability
		err := rows.Scan(
			&a.ID,
			&a.WorkerID,
			&a.Date,
			&a.StartTime,
			&a.EndTime,
			&a.IsConfirmed,
			&a.ConfirmedAt,
			&a.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		availabilities = append(availabilities, a)
	}
	
	return availabilities, nil
}

// ConfirmAvailability confirma una disponibilidad
func (r *GigWorkerRepository) ConfirmAvailability(availabilityID string) error {
	query := `
		UPDATE worker_availability
		SET is_confirmed = true, confirmed_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`
	
	_, err := r.db.Exec(query, availabilityID)
	return err
}

// CreateTask crea una nueva tarea
func (r *GigWorkerRepository) CreateTask(task *models.Task) error {
	query := `
		INSERT INTO tasks (title, description, location, scheduled_date, scheduled_start_time, 
		                  scheduled_end_time, status, priority)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at
	`
	
	err := r.db.QueryRow(
		query,
		task.Title,
		task.Description,
		task.Location,
		task.ScheduledDate,
		task.ScheduledStartTime,
		task.ScheduledEndTime,
		task.Status,
		task.Priority,
	).Scan(&task.ID, &task.CreatedAt, &task.UpdatedAt)
	
	return err
}

// GetTaskByID obtiene una tarea por su ID
func (r *GigWorkerRepository) GetTaskByID(id string) (*models.Task, error) {
	task := &models.Task{}
	query := `
		SELECT id, title, description, location, scheduled_date, scheduled_start_time,
		       scheduled_end_time, status, priority, created_at, updated_at
		FROM tasks
		WHERE id = $1
	`
	
	err := r.db.QueryRow(query, id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Location,
		&task.ScheduledDate,
		&task.ScheduledStartTime,
		&task.ScheduledEndTime,
		&task.Status,
		&task.Priority,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("task not found")
	}
	
	return task, err
}

// GetPendingTasks obtiene todas las tareas pendientes
func (r *GigWorkerRepository) GetPendingTasks() ([]models.Task, error) {
	query := `
		SELECT id, title, description, location, scheduled_date, scheduled_start_time,
		       scheduled_end_time, status, priority, created_at, updated_at
		FROM tasks
		WHERE status = 'pending'
		ORDER BY scheduled_date, scheduled_start_time
	`
	
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var tasks []models.Task
	for rows.Next() {
		var t models.Task
		err := rows.Scan(
			&t.ID,
			&t.Title,
			&t.Description,
			&t.Location,
			&t.ScheduledDate,
			&t.ScheduledStartTime,
			&t.ScheduledEndTime,
			&t.Status,
			&t.Priority,
			&t.CreatedAt,
			&t.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	
	return tasks, nil
}

// AssignTask asigna una tarea a un trabajador
func (r *GigWorkerRepository) AssignTask(assignment *models.TaskAssignment) error {
	// Iniciar transacción
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	
	// Insertar asignación
	query := `
		INSERT INTO task_assignments (task_id, worker_id, status)
		VALUES ($1, $2, $3)
		RETURNING id, assigned_at
	`
	
	err = tx.QueryRow(
		query,
		assignment.TaskID,
		assignment.WorkerID,
		assignment.Status,
	).Scan(&assignment.ID, &assignment.AssignedAt)
	
	if err != nil {
		return err
	}
	
	// Actualizar estado de la tarea
	updateQuery := `
		UPDATE tasks
		SET status = 'assigned', updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`
	
	_, err = tx.Exec(updateQuery, assignment.TaskID)
	if err != nil {
		return err
	}
	
	return tx.Commit()
}

// UpdateTaskAssignment actualiza una asignación de tarea
func (r *GigWorkerRepository) UpdateTaskAssignment(assignment *models.TaskAssignment) error {
	query := `
		UPDATE task_assignments
		SET status = $2, accepted_at = $3, started_at = $4, completed_at = $5, 
		    rating = $6, comments = $7
		WHERE id = $1
	`
	
	_, err := r.db.Exec(
		query,
		assignment.ID,
		assignment.Status,
		assignment.AcceptedAt,
		assignment.StartedAt,
		assignment.CompletedAt,
		assignment.Rating,
		assignment.Comments,
	)
	
	return err
}

// GetWorkerAssignments obtiene las asignaciones de un trabajador
func (r *GigWorkerRepository) GetWorkerAssignments(workerID string, status string) ([]models.TaskAssignment, error) {
	query := `
		SELECT ta.id, ta.task_id, ta.worker_id, ta.assigned_at, ta.accepted_at, 
		       ta.started_at, ta.completed_at, ta.status, ta.rating, ta.comments
		FROM task_assignments ta
		WHERE ta.worker_id = $1
	`
	
	args := []interface{}{workerID}
	if status != "" {
		query += " AND ta.status = $2"
		args = append(args, status)
	}
	
	query += " ORDER BY ta.assigned_at DESC"
	
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var assignments []models.TaskAssignment
	for rows.Next() {
		var a models.TaskAssignment
		err := rows.Scan(
			&a.ID,
			&a.TaskID,
			&a.WorkerID,
			&a.AssignedAt,
			&a.AcceptedAt,
			&a.StartedAt,
			&a.CompletedAt,
			&a.Status,
			&a.Rating,
			&a.Comments,
		)
		if err != nil {
			return nil, err
		}
		assignments = append(assignments, a)
	}
	
	return assignments, nil
}

// RecordAttendance registra la asistencia de un trabajador
func (r *GigWorkerRepository) RecordAttendance(attendance *models.AttendanceHistory) error {
	query := `
		INSERT INTO attendance_history (worker_id, scheduled_date, scheduled_time, 
		                               did_confirm, did_attend, absence_reason)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`
	
	err := r.db.QueryRow(
		query,
		attendance.WorkerID,
		attendance.ScheduledDate,
		attendance.ScheduledTime,
		attendance.DidConfirm,
		attendance.DidAttend,
		attendance.AbsenceReason,
	).Scan(&attendance.ID, &attendance.CreatedAt)
	
	return err
}

// GetAttendanceHistory obtiene el historial de asistencia de un trabajador
func (r *GigWorkerRepository) GetAttendanceHistory(workerID string, limit int) ([]models.AttendanceHistory, error) {
	query := `
		SELECT id, worker_id, scheduled_date, scheduled_time, did_confirm, 
		       did_attend, absence_reason, created_at
		FROM attendance_history
		WHERE worker_id = $1
		ORDER BY scheduled_date DESC, scheduled_time DESC
		LIMIT $2
	`
	
	rows, err := r.db.Query(query, workerID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var history []models.AttendanceHistory
	for rows.Next() {
		var h models.AttendanceHistory
		err := rows.Scan(
			&h.ID,
			&h.WorkerID,
			&h.ScheduledDate,
			&h.ScheduledTime,
			&h.DidConfirm,
			&h.DidAttend,
			&h.AbsenceReason,
			&h.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		history = append(history, h)
	}
	
	return history, nil
}

// CreateNotification crea una nueva notificación
func (r *GigWorkerRepository) CreateNotification(notification *models.Notification) error {
	query := `
		INSERT INTO notifications (user_id, type, title, message)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`
	
	err := r.db.QueryRow(
		query,
		notification.UserID,
		notification.Type,
		notification.Title,
		notification.Message,
	).Scan(&notification.ID, &notification.CreatedAt)
	
	return err
}

// GetUnreadNotifications obtiene las notificaciones no leídas de un usuario
func (r *GigWorkerRepository) GetUnreadNotifications(userID string) ([]models.Notification, error) {
	query := `
		SELECT id, user_id, type, title, message, is_read, read_at, created_at
		FROM notifications
		WHERE user_id = $1 AND is_read = false
		ORDER BY created_at DESC
	`
	
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var notifications []models.Notification
	for rows.Next() {
		var n models.Notification
		err := rows.Scan(
			&n.ID,
			&n.UserID,
			&n.Type,
			&n.Title,
			&n.Message,
			&n.IsRead,
			&n.ReadAt,
			&n.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, n)
	}
	
	return notifications, nil
}

// MarkNotificationAsRead marca una notificación como leída
func (r *GigWorkerRepository) MarkNotificationAsRead(notificationID string) error {
	query := `
		UPDATE notifications
		SET is_read = true, read_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`
	
	_, err := r.db.Exec(query, notificationID)
	return err
}