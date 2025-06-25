<template>
  <div class="gig-worker-dashboard">
    <!-- Header con información del trabajador -->
    <div class="worker-header bg-white shadow-sm p-6 rounded-lg mb-6">
      <div v-if="workerProfile" class="flex justify-between items-center">
        <div>
          <h1 class="text-2xl font-bold text-gray-800">Panel de Trabajador</h1>
          <p class="text-gray-600 mt-1">Código: {{ workerProfile.workerCode }}</p>
        </div>
        <div class="text-right">
          <div class="flex items-center mb-2">
            <span class="text-yellow-500 mr-1">⭐</span>
            <span class="font-semibold">{{ workerProfile.rating.toFixed(1) }}</span>
          </div>
          <p class="text-sm text-gray-600">{{ workerProfile.totalTasksCompleted }} tareas completadas</p>
        </div>
      </div>
      
      <!-- Botón de registro si no es trabajador -->
      <div v-else-if="!loading" class="text-center py-8">
        <h2 class="text-xl font-semibold mb-4">No estás registrado como trabajador</h2>
        <button 
          @click="showRegisterModal = true"
          class="bg-blue-500 text-white px-6 py-2 rounded-lg hover:bg-blue-600"
        >
          Registrarse como Trabajador
        </button>
      </div>
    </div>

    <!-- Tabs de navegación -->
    <div v-if="workerProfile" class="bg-white shadow-sm rounded-lg mb-6">
      <div class="flex border-b">
        <button
          v-for="tab in tabs"
          :key="tab.id"
          @click="activeTab = tab.id"
          :class="[
            'px-6 py-3 font-medium transition-colors',
            activeTab === tab.id 
              ? 'text-blue-600 border-b-2 border-blue-600' 
              : 'text-gray-600 hover:text-gray-800'
          ]"
        >
          {{ tab.label }}
        </button>
      </div>
    </div>

    <!-- Contenido según tab activa -->
    <div v-if="workerProfile" class="content-area">
      <!-- Tab de Tareas Disponibles -->
      <div v-if="activeTab === 'tasks'" class="space-y-4">
        <h2 class="text-xl font-semibold mb-4">Tareas Disponibles</h2>
        <div v-if="availableTasks.length === 0" class="text-center py-8 bg-gray-50 rounded-lg">
          <p class="text-gray-600">No hay tareas disponibles en este momento</p>
        </div>
        <div v-else class="grid gap-4">
          <div 
            v-for="task in availableTasks" 
            :key="task.id"
            class="bg-white p-6 rounded-lg shadow-sm border hover:shadow-md transition-shadow"
          >
            <div class="flex justify-between items-start mb-3">
              <h3 class="text-lg font-semibold">{{ task.title }}</h3>
              <span :class="[
                'px-3 py-1 rounded-full text-sm font-medium',
                getPriorityClass(task.priority)
              ]">
                {{ task.priority }}
              </span>
            </div>
            <p v-if="task.description" class="text-gray-600 mb-3">{{ task.description }}</p>
            <div class="flex flex-wrap gap-4 text-sm text-gray-600 mb-4">
              <div class="flex items-center">
                <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"></path>
                </svg>
                {{ formatDate(task.scheduledDate) }}
              </div>
              <div class="flex items-center">
                <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                </svg>
                {{ task.scheduledStartTime }} - {{ task.scheduledEndTime }}
              </div>
              <div v-if="task.location" class="flex items-center">
                <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z"></path>
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z"></path>
                </svg>
                {{ task.location }}
              </div>
            </div>
            <button
              @click="acceptTask(task.id)"
              class="w-full bg-blue-500 text-white py-2 rounded-lg hover:bg-blue-600 transition-colors"
            >
              Aceptar Tarea
            </button>
          </div>
        </div>
      </div>

      <!-- Tab de Mis Asignaciones -->
      <div v-if="activeTab === 'assignments'" class="space-y-4">
        <h2 class="text-xl font-semibold mb-4">Mis Asignaciones</h2>
        <div v-if="myAssignments.length === 0" class="text-center py-8 bg-gray-50 rounded-lg">
          <p class="text-gray-600">No tienes asignaciones activas</p>
        </div>
        <div v-else class="grid gap-4">
          <div 
            v-for="assignment in myAssignments" 
            :key="assignment.id"
            class="bg-white p-6 rounded-lg shadow-sm border"
          >
            <div class="flex justify-between items-start mb-3">
              <h3 class="text-lg font-semibold">{{ assignment.task?.title || 'Tarea' }}</h3>
              <span :class="[
                'px-3 py-1 rounded-full text-sm font-medium',
                getStatusClass(assignment.status)
              ]">
                {{ getStatusLabel(assignment.status) }}
              </span>
            </div>
            <div class="space-y-2 text-sm text-gray-600 mb-4">
              <p>Asignado: {{ formatDateTime(assignment.assignedAt) }}</p>
              <p v-if="assignment.acceptedAt">Aceptado: {{ formatDateTime(assignment.acceptedAt) }}</p>
              <p v-if="assignment.startedAt">Iniciado: {{ formatDateTime(assignment.startedAt) }}</p>
            </div>
            <div class="flex gap-2">
              <button
                v-if="assignment.status === 'assigned'"
                @click="acceptAssignment(assignment)"
                class="flex-1 bg-green-500 text-white py-2 rounded-lg hover:bg-green-600"
              >
                Aceptar
              </button>
              <button
                v-if="assignment.status === 'accepted'"
                @click="startAssignment(assignment.id)"
                class="flex-1 bg-blue-500 text-white py-2 rounded-lg hover:bg-blue-600"
              >
                Iniciar
              </button>
              <button
                v-if="assignment.status === 'in_progress'"
                @click="completeAssignment(assignment.id)"
                class="flex-1 bg-purple-500 text-white py-2 rounded-lg hover:bg-purple-600"
              >
                Completar
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Tab de Disponibilidad -->
      <div v-if="activeTab === 'availability'" class="space-y-6">
        <h2 class="text-xl font-semibold mb-4">Mi Disponibilidad</h2>
        
        <!-- Formulario para agregar disponibilidad -->
        <div class="bg-white p-6 rounded-lg shadow-sm">
          <h3 class="text-lg font-medium mb-4">Agregar Disponibilidad</h3>
          <form @submit.prevent="addAvailability" class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Fecha</label>
              <input
                v-model="availabilityForm.date"
                type="date"
                required
                class="w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
            </div>
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Hora inicio</label>
                <input
                  v-model="availabilityForm.startTime"
                  type="time"
                  required
                  class="w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                >
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Hora fin</label>
                <input
                  v-model="availabilityForm.endTime"
                  type="time"
                  required
                  class="w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                >
              </div>
            </div>
            <button
              type="submit"
              class="w-full bg-blue-500 text-white py-2 rounded-lg hover:bg-blue-600"
            >
              Agregar Disponibilidad
            </button>
          </form>
        </div>

        <!-- Lista de disponibilidades -->
        <div class="bg-white p-6 rounded-lg shadow-sm">
          <h3 class="text-lg font-medium mb-4">Próximas Disponibilidades</h3>
          <div v-if="myAvailabilities.length === 0" class="text-center py-4 text-gray-600">
            No has registrado disponibilidad
          </div>
          <div v-else class="space-y-3">
            <div 
              v-for="availability in myAvailabilities" 
              :key="availability.id"
              class="flex justify-between items-center p-4 border rounded-lg"
            >
              <div>
                <p class="font-medium">{{ formatDate(availability.date) }}</p>
                <p class="text-sm text-gray-600">{{ availability.startTime }} - {{ availability.endTime }}</p>
              </div>
              <div class="flex items-center gap-2">
                <span v-if="availability.isConfirmed" class="text-green-600 text-sm">
                  ✓ Confirmado
                </span>
                <button
                  v-else
                  @click="confirmAvailability(availability.id)"
                  class="px-4 py-1 bg-green-500 text-white rounded hover:bg-green-600"
                >
                  Confirmar
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Tab de Estadísticas -->
      <div v-if="activeTab === 'stats'" class="space-y-6">
        <h2 class="text-xl font-semibold mb-4">Mis Estadísticas</h2>
        
        <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
          <div class="bg-white p-6 rounded-lg shadow-sm">
            <h3 class="text-lg font-medium text-gray-700 mb-2">Tasa de Asistencia</h3>
            <p class="text-3xl font-bold text-blue-600">{{ attendanceStats.attendanceRate.toFixed(1) }}%</p>
          </div>
          <div class="bg-white p-6 rounded-lg shadow-sm">
            <h3 class="text-lg font-medium text-gray-700 mb-2">Tasa de Confirmación</h3>
            <p class="text-3xl font-bold text-green-600">{{ attendanceStats.confirmationRate.toFixed(1) }}%</p>
          </div>
          <div class="bg-white p-6 rounded-lg shadow-sm">
            <h3 class="text-lg font-medium text-gray-700 mb-2">Tareas Completadas</h3>
            <p class="text-3xl font-bold text-purple-600">{{ workerProfile.totalTasksCompleted }}</p>
          </div>
        </div>

        <!-- Historial de asistencia -->
        <div class="bg-white p-6 rounded-lg shadow-sm">
          <h3 class="text-lg font-medium mb-4">Historial de Asistencia</h3>
          <div v-if="attendanceHistory.length === 0" class="text-center py-4 text-gray-600">
            No hay historial disponible
          </div>
          <div v-else class="overflow-x-auto">
            <table class="w-full">
              <thead>
                <tr class="border-b">
                  <th class="text-left py-2 px-4">Fecha</th>
                  <th class="text-left py-2 px-4">Hora</th>
                  <th class="text-left py-2 px-4">Confirmado</th>
                  <th class="text-left py-2 px-4">Asistió</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="record in attendanceHistory" :key="record.id" class="border-b">
                  <td class="py-2 px-4">{{ formatDate(record.scheduledDate) }}</td>
                  <td class="py-2 px-4">{{ record.scheduledTime }}</td>
                  <td class="py-2 px-4">
                    <span :class="record.didConfirm ? 'text-green-600' : 'text-red-600'">
                      {{ record.didConfirm ? 'Sí' : 'No' }}
                    </span>
                  </td>
                  <td class="py-2 px-4">
                    <span :class="record.didAttend ? 'text-green-600' : 'text-red-600'">
                      {{ record.didAttend ? 'Sí' : 'No' }}
                    </span>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>

    <!-- Modal de registro -->
    <div v-if="showRegisterModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white p-6 rounded-lg max-w-md w-full">
        <h2 class="text-xl font-bold mb-4">Registrarse como Trabajador</h2>
        <form @submit.prevent="registerAsWorker">
          <div class="mb-4">
            <label class="block text-sm font-medium text-gray-700 mb-1">Código de Trabajador</label>
            <input
              v-model="registerForm.workerCode"
              type="text"
              required
              placeholder="Ej: W12345"
              class="w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
          </div>
          <div class="flex gap-2">
            <button
              type="submit"
              class="flex-1 bg-blue-500 text-white py-2 rounded-lg hover:bg-blue-600"
            >
              Registrar
            </button>
            <button
              type="button"
              @click="showRegisterModal = false"
              class="flex-1 bg-gray-300 text-gray-700 py-2 rounded-lg hover:bg-gray-400"
            >
              Cancelar
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'
import { gigWorkerApi, type GigWorker, type Task, type TaskAssignment, type WorkerAvailability, type AttendanceHistory } from '../services/gigWorkerApi'

// Estado reactivo
const loading = ref(true)
const workerProfile = ref<GigWorker | null>(null)
const activeTab = ref('tasks')
const showRegisterModal = ref(false)

// Datos
const availableTasks = ref<Task[]>([])
const myAssignments = ref<TaskAssignment[]>([])
const myAvailabilities = ref<WorkerAvailability[]>([])
const attendanceHistory = ref<AttendanceHistory[]>([])
const attendanceStats = ref({
  attendanceRate: 0,
  confirmationRate: 0
})

// Formularios
const registerForm = ref({
  workerCode: ''
})

const availabilityForm = ref({
  date: '',
  startTime: '',
  endTime: ''
})

// Tabs de navegación
const tabs = [
  { id: 'tasks', label: 'Tareas Disponibles' },
  { id: 'assignments', label: 'Mis Asignaciones' },
  { id: 'availability', label: 'Disponibilidad' },
  { id: 'stats', label: 'Estadísticas' }
]

// Funciones de utilidad
const formatDate = (dateString: string) => {
  const date = new Date(dateString)
  return date.toLocaleDateString('es-ES', { 
    year: 'numeric', 
    month: 'long', 
    day: 'numeric' 
  })
}

const formatDateTime = (dateString: string) => {
  const date = new Date(dateString)
  return date.toLocaleString('es-ES', { 
    year: 'numeric', 
    month: 'short', 
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const getPriorityClass = (priority: string) => {
  const classes = {
    low: 'bg-gray-100 text-gray-600',
    normal: 'bg-blue-100 text-blue-600',
    high: 'bg-orange-100 text-orange-600',
    urgent: 'bg-red-100 text-red-600'
  }
  return classes[priority as keyof typeof classes] || classes.normal
}

const getStatusClass = (status: string) => {
  const classes = {
    assigned: 'bg-yellow-100 text-yellow-600',
    accepted: 'bg-blue-100 text-blue-600',
    in_progress: 'bg-purple-100 text-purple-600',
    completed: 'bg-green-100 text-green-600',
    rejected: 'bg-red-100 text-red-600'
  }
  return classes[status as keyof typeof classes] || 'bg-gray-100 text-gray-600'
}

const getStatusLabel = (status: string) => {
  const labels = {
    assigned: 'Asignado',
    accepted: 'Aceptado',
    in_progress: 'En Progreso',
    completed: 'Completado',
    rejected: 'Rechazado'
  }
  return labels[status as keyof typeof labels] || status
}

// Funciones de API
const loadWorkerProfile = async () => {
  try {
    workerProfile.value = await gigWorkerApi.getMyProfile()
    await loadDashboardData()
  } catch (error) {
    console.log('Usuario no registrado como trabajador')
  } finally {
    loading.value = false
  }
}

const loadDashboardData = async () => {
  if (!workerProfile.value) return

  // Cargar datos según la tab activa
  if (activeTab.value === 'tasks') {
    availableTasks.value = await gigWorkerApi.getAvailableTasks()
  } else if (activeTab.value === 'assignments') {
    myAssignments.value = await gigWorkerApi.getMyAssignments()
  } else if (activeTab.value === 'availability') {
    const today = new Date().toISOString().split('T')[0]
    myAvailabilities.value = await gigWorkerApi.getMyAvailability(today)
  } else if (activeTab.value === 'stats') {
    const stats = await gigWorkerApi.getAttendanceHistory()
    attendanceHistory.value = stats.history
    attendanceStats.value = {
      attendanceRate: stats.attendanceRate,
      confirmationRate: stats.confirmationRate
    }
  }
}

const registerAsWorker = async () => {
  try {
    workerProfile.value = await gigWorkerApi.registerAsWorker(registerForm.value.workerCode)
    showRegisterModal.value = false
    await loadDashboardData()
  } catch (error) {
    alert('Error al registrarse como trabajador')
  }
}

const acceptTask = async (taskId: string) => {
  try {
    await gigWorkerApi.acceptTask(taskId)
    alert('Tarea aceptada exitosamente')
    // Recargar tareas y asignaciones
    availableTasks.value = await gigWorkerApi.getAvailableTasks()
    activeTab.value = 'assignments'
    myAssignments.value = await gigWorkerApi.getMyAssignments()
  } catch (error) {
    alert('Error al aceptar la tarea')
  }
}

const acceptAssignment = async (assignment: TaskAssignment) => {
  // La tarea ya fue aceptada al hacer click en "Aceptar Tarea"
  // Esta función es para cuando la tarea fue asignada por un admin
}

const startAssignment = async (assignmentId: string) => {
  try {
    await gigWorkerApi.startTask(assignmentId)
    alert('Tarea iniciada')
    myAssignments.value = await gigWorkerApi.getMyAssignments()
  } catch (error) {
    alert('Error al iniciar la tarea')
  }
}

const completeAssignment = async (assignmentId: string) => {
  const comments = prompt('Comentarios (opcional):')
  try {
    await gigWorkerApi.completeTask(assignmentId, comments || undefined)
    alert('Tarea completada')
    myAssignments.value = await gigWorkerApi.getMyAssignments()
  } catch (error) {
    alert('Error al completar la tarea')
  }
}

const addAvailability = async () => {
  try {
    await gigWorkerApi.updateAvailability(availabilityForm.value)
    alert('Disponibilidad agregada')
    // Recargar disponibilidades
    const date = availabilityForm.value.date
    myAvailabilities.value = await gigWorkerApi.getMyAvailability(date)
    // Limpiar formulario
    availabilityForm.value = { date: '', startTime: '', endTime: '' }
  } catch (error) {
    alert('Error al agregar disponibilidad')
  }
}

const confirmAvailability = async (availabilityId: string) => {
  try {
    await gigWorkerApi.confirmAvailability(availabilityId)
    alert('Disponibilidad confirmada')
    // Recargar disponibilidades
    const today = new Date().toISOString().split('T')[0]
    myAvailabilities.value = await gigWorkerApi.getMyAvailability(today)
  } catch (error) {
    alert('Error al confirmar disponibilidad')
  }
}

// Lifecycle
onMounted(() => {
  loadWorkerProfile()
})

// Watchers
// Recargar datos cuando cambia la tab
watch(activeTab, () => {
  loadDashboardData()
})
</script>

<style scoped>
.gig-worker-dashboard {
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem;
}

@media (max-width: 768px) {
  .gig-worker-dashboard {
    padding: 1rem;
  }
}
</style>