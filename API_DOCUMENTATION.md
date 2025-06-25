# Documentación API - Clon WhatsApp con Gestión de Gig Workers

## Base URL
```
http://localhost:8080/api
```

## Autenticación
Todas las rutas protegidas requieren un token JWT en el header:
```
Authorization: Bearer <token>
```

## Endpoints de Gig Workers

### 1. Registrar como Gig Worker
**POST** `/api/workers/register`

Registra a un usuario autenticado como gig worker.

**Request Body:**
```json
{
  "workerCode": "W12345"
}
```

**Response:** `201 Created`
```json
{
  "id": "uuid",
  "userId": "uuid",
  "workerCode": "W12345",
  "status": "active",
  "rating": 5.0,
  "totalTasksCompleted": 0,
  "createdAt": "2024-01-20T10:00:00Z",
  "updatedAt": "2024-01-20T10:00:00Z"
}
```

### 2. Obtener Perfil del Trabajador
**GET** `/api/workers/profile`

Obtiene el perfil del trabajador actual.

**Response:** `200 OK`
```json
{
  "id": "uuid",
  "userId": "uuid",
  "workerCode": "W12345",
  "status": "active",
  "rating": 4.8,
  "totalTasksCompleted": 150,
  "createdAt": "2024-01-20T10:00:00Z",
  "updatedAt": "2024-01-20T10:00:00Z"
}
```

### 3. Actualizar Disponibilidad
**POST** `/api/workers/availability`

Actualiza la disponibilidad del trabajador para una fecha específica.

**Request Body:**
```json
{
  "date": "2024-01-25",
  "startTime": "09:00",
  "endTime": "18:00"
}
```

**Response:** `201 Created`
```json
{
  "id": "uuid",
  "workerId": "uuid",
  "date": "2024-01-25T00:00:00Z",
  "startTime": "09:00",
  "endTime": "18:00",
  "isConfirmed": false,
  "createdAt": "2024-01-20T10:00:00Z"
}
```

### 4. Confirmar Disponibilidad
**PUT** `/api/workers/availability/:id/confirm`

Confirma la disponibilidad para un día específico.

**Response:** `200 OK`
```json
{
  "message": "Disponibilidad confirmada"
}
```

### 5. Obtener Disponibilidad
**GET** `/api/workers/availability?date=2024-01-25`

Obtiene la disponibilidad del trabajador para una fecha específica.

**Response:** `200 OK`
```json
[
  {
    "id": "uuid",
    "workerId": "uuid",
    "date": "2024-01-25T00:00:00Z",
    "startTime": "09:00",
    "endTime": "18:00",
    "isConfirmed": true,
    "confirmedAt": "2024-01-20T11:00:00Z",
    "createdAt": "2024-01-20T10:00:00Z"
  }
]
```

### 6. Obtener Tareas Disponibles
**GET** `/api/tasks`

Obtiene las tareas pendientes disponibles para asignar.

**Response:** `200 OK`
```json
[
  {
    "id": "uuid",
    "title": "Entrega en Zona Norte",
    "description": "Recoger paquete en almacén y entregar",
    "location": "Zona Norte, Santiago",
    "scheduledDate": "2024-01-25T00:00:00Z",
    "scheduledStartTime": "10:00",
    "scheduledEndTime": "11:00",
    "status": "pending",
    "priority": "normal",
    "createdAt": "2024-01-20T08:00:00Z",
    "updatedAt": "2024-01-20T08:00:00Z"
  }
]
```

### 7. Aceptar Tarea
**POST** `/api/tasks/:id/accept`

Acepta una tarea disponible.

**Response:** `200 OK`
```json
{
  "id": "uuid",
  "taskId": "uuid",
  "workerId": "uuid",
  "assignedAt": "2024-01-20T10:30:00Z",
  "acceptedAt": "2024-01-20T10:30:00Z",
  "status": "accepted"
}
```

### 8. Iniciar Tarea
**PUT** `/api/assignments/:id/start`

Marca el inicio de una tarea asignada.

**Response:** `200 OK`
```json
{
  "message": "Tarea iniciada"
}
```

### 9. Completar Tarea
**PUT** `/api/assignments/:id/complete`

Marca una tarea como completada.

**Request Body (opcional):**
```json
{
  "comments": "Entrega realizada sin problemas"
}
```

**Response:** `200 OK`
```json
{
  "message": "Tarea completada"
}
```

### 10. Obtener Mis Asignaciones
**GET** `/api/workers/assignments?status=active`

Obtiene las asignaciones del trabajador actual.

**Query Parameters:**
- `status` (opcional): Filtrar por estado (assigned, accepted, in_progress, completed)

**Response:** `200 OK`
```json
[
  {
    "id": "uuid",
    "taskId": "uuid",
    "workerId": "uuid",
    "assignedAt": "2024-01-20T10:30:00Z",
    "acceptedAt": "2024-01-20T10:31:00Z",
    "startedAt": "2024-01-20T11:00:00Z",
    "status": "in_progress"
  }
]
```

### 11. Obtener Historial de Asistencia
**GET** `/api/workers/attendance?limit=30`

Obtiene el historial de asistencia con estadísticas.

**Response:** `200 OK`
```json
{
  "history": [
    {
      "id": "uuid",
      "workerId": "uuid",
      "scheduledDate": "2024-01-19T00:00:00Z",
      "scheduledTime": "09:00",
      "didConfirm": true,
      "didAttend": true,
      "createdAt": "2024-01-19T18:00:00Z"
    }
  ],
  "attendanceRate": 95.5,
  "confirmationRate": 98.0
}
```

### 12. Obtener Notificaciones
**GET** `/api/notifications`

Obtiene las notificaciones no leídas del usuario.

**Response:** `200 OK`
```json
[
  {
    "id": "uuid",
    "userId": "uuid",
    "type": "task_assigned",
    "title": "Nueva tarea asignada",
    "message": "Se te ha asignado una nueva tarea para mañana",
    "isRead": false,
    "createdAt": "2024-01-20T10:00:00Z"
  }
]
```

### 13. Marcar Notificación como Leída
**PUT** `/api/notifications/:id/read`

Marca una notificación como leída.

**Response:** `200 OK`
```json
{
  "message": "Notificación marcada como leída"
}
```

## Endpoints Administrativos

### 14. Crear Tarea (Admin)
**POST** `/api/tasks`

Crea una nueva tarea.

**Request Body:**
```json
{
  "title": "Entrega urgente Centro",
  "description": "Entregar documentos importantes",
  "location": "Santiago Centro",
  "scheduledDate": "2024-01-26",
  "scheduledStartTime": "14:00",
  "scheduledEndTime": "15:00",
  "priority": "high"
}
```

**Response:** `201 Created`

### 15. Registrar Asistencia (Admin)
**POST** `/api/workers/attendance`

Registra la asistencia de un trabajador.

**Request Body:**
```json
{
  "workerId": "uuid",
  "scheduledDate": "2024-01-20",
  "scheduledTime": "09:00",
  "didConfirm": true,
  "didAttend": true
}
```

**Response:** `201 Created`

## Códigos de Estado

- `200 OK`: Solicitud exitosa
- `201 Created`: Recurso creado exitosamente
- `400 Bad Request`: Datos de entrada inválidos
- `401 Unauthorized`: Token no válido o faltante
- `403 Forbidden`: Sin permisos para realizar la acción
- `404 Not Found`: Recurso no encontrado
- `409 Conflict`: Conflicto (ej: usuario ya registrado como trabajador)
- `500 Internal Server Error`: Error del servidor