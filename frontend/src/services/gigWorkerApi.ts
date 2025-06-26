import axios from 'axios'

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:5000/api'

// Crear instancia de axios con configuración base
const apiClient = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json'
  }
})

// Interceptor para agregar el token a todas las peticiones
apiClient.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('auth_token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Tipos para TypeScript
export interface GigWorker {
  id: string
  userId: string
  workerCode: string
  status: 'active' | 'inactive' | 'suspended'
  rating: number
  totalTasksCompleted: number
  createdAt: string
  updatedAt: string
}

export interface WorkerAvailability {
  id: string
  workerId: string
  date: string
  startTime: string
  endTime: string
  isConfirmed: boolean
  confirmedAt?: string
  createdAt: string
}

export interface Task {
  id: string
  title: string
  description?: string
  location?: string
  scheduledDate: string
  scheduledStartTime: string
  scheduledEndTime: string
  status: 'pending' | 'assigned' | 'in_progress' | 'completed' | 'cancelled'
  priority: 'low' | 'normal' | 'high' | 'urgent'
  createdAt: string
  updatedAt: string
}

export interface TaskAssignment {
  id: string
  taskId: string
  workerId: string
  assignedAt: string
  acceptedAt?: string
  startedAt?: string
  completedAt?: string
  status: 'assigned' | 'accepted' | 'rejected' | 'in_progress' | 'completed'
  rating?: number
  comments?: string
  task?: Task
}

export interface AttendanceHistory {
  id: string
  workerId: string
  scheduledDate: string
  scheduledTime: string
  didConfirm: boolean
  didAttend: boolean
  absenceReason?: string
  createdAt: string
}

export interface Notification {
  id: string
  userId: string
  type: 'task_assigned' | 'reminder' | 'system'
  title: string
  message: string
  isRead: boolean
  readAt?: string
  createdAt: string
}

// API Service
export const gigWorkerApi = {
  // Registro y perfil
  registerAsWorker: async (workerCode: string): Promise<GigWorker> => {
    const response = await apiClient.post('/workers/register', { workerCode })
    return response.data
  },

  getMyProfile: async (): Promise<GigWorker> => {
    const response = await apiClient.get('/workers/profile')
    return response.data
  },

  getWorkerProfile: async (workerId: string): Promise<GigWorker> => {
    const response = await apiClient.get(`/workers/${workerId}`)
    return response.data
  },

  // Disponibilidad
  updateAvailability: async (data: {
    date: string
    startTime: string
    endTime: string
  }): Promise<WorkerAvailability> => {
    const response = await apiClient.post('/workers/availability', data)
    return response.data
  },

  confirmAvailability: async (availabilityId: string): Promise<{ message: string }> => {
    const response = await apiClient.put(`/workers/availability/${availabilityId}/confirm`)
    return response.data
  },

  getMyAvailability: async (date: string): Promise<WorkerAvailability[]> => {
    const response = await apiClient.get('/workers/availability', { params: { date } })
    return response.data
  },

  // Tareas
  getAvailableTasks: async (): Promise<Task[]> => {
    const response = await apiClient.get('/tasks')
    return response.data
  },

  createTask: async (task: Omit<Task, 'id' | 'status' | 'createdAt' | 'updatedAt'>): Promise<Task> => {
    const response = await apiClient.post('/tasks', task)
    return response.data
  },

  acceptTask: async (taskId: string): Promise<TaskAssignment> => {
    const response = await apiClient.post(`/tasks/${taskId}/accept`)
    return response.data
  },

  startTask: async (assignmentId: string): Promise<{ message: string }> => {
    const response = await apiClient.put(`/assignments/${assignmentId}/start`)
    return response.data
  },

  completeTask: async (assignmentId: string, comments?: string): Promise<{ message: string }> => {
    const response = await apiClient.put(`/assignments/${assignmentId}/complete`, { comments })
    return response.data
  },

  getMyAssignments: async (status?: string): Promise<TaskAssignment[]> => {
    const response = await apiClient.get('/workers/assignments', { params: { status } })
    return response.data
  },

  // Historial y estadísticas
  getAttendanceHistory: async (limit: number = 30): Promise<{
    history: AttendanceHistory[]
    attendanceRate: number
    confirmationRate: number
  }> => {
    const response = await apiClient.get('/workers/attendance', { params: { limit } })
    return response.data
  },

  recordAttendance: async (data: {
    workerId: string
    scheduledDate: string
    scheduledTime: string
    didConfirm: boolean
    didAttend: boolean
    absenceReason?: string
  }): Promise<AttendanceHistory> => {
    const response = await apiClient.post('/workers/attendance', data)
    return response.data
  },

  // Notificaciones
  getNotifications: async (): Promise<Notification[]> => {
    const response = await apiClient.get('/notifications')
    return response.data
  },

  markNotificationAsRead: async (notificationId: string): Promise<{ message: string }> => {
    const response = await apiClient.put(`/notifications/${notificationId}/read`)
    return response.data
  }
}

// Hook para manejar errores de forma consistente
export const useApiError = () => {
  const handleError = (error: any) => {
    if (error.response) {
      // Error de respuesta del servidor
      const message = error.response.data.error || 'Error en el servidor'
      console.error('API Error:', message)
      return message
    } else if (error.request) {
      // La petición se hizo pero no hubo respuesta
      console.error('No response from server')
      return 'No se pudo conectar con el servidor'
    } else {
      // Algo pasó al configurar la petición
      console.error('Request error:', error.message)
      return 'Error al realizar la petición'
    }
  }

  return { handleError }
}