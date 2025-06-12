<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-50">
    <div class="max-w-md w-full space-y-8">
      <div>
        <h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900">
          WhatsApp Clone
        </h2>
        <p class="mt-2 text-center text-sm text-gray-600">
          {{ isLogin ? 'Sign in to your account' : 'Create a new account' }}
        </p>
      </div>
      
      <form class="mt-8 space-y-6" @submit.prevent="handleSubmit">
        <div class="rounded-md shadow-sm -space-y-px">
          <div v-if="!isLogin">
            <label for="name" class="sr-only">Name</label>
            <input
              id="name"
              v-model="formData.name"
              name="name"
              type="text"
              required
              class="appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-t-md focus:outline-none focus:ring-green-500 focus:border-green-500 focus:z-10 sm:text-sm"
              placeholder="Full Name"
            />
          </div>
          
          <div>
            <label for="email" class="sr-only">Email address</label>
            <input
              id="email"
              v-model="formData.email"
              name="email"
              type="email"
              autocomplete="email"
              required
              class="appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 focus:outline-none focus:ring-green-500 focus:border-green-500 focus:z-10 sm:text-sm"
              :class="{ 'rounded-t-md': isLogin }"
              placeholder="Email address"
            />
          </div>
          
          <div>
            <label for="password" class="sr-only">Password</label>
            <input
              id="password"
              v-model="formData.password"
              name="password"
              type="password"
              autocomplete="current-password"
              required
              class="appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-b-md focus:outline-none focus:ring-green-500 focus:border-green-500 focus:z-10 sm:text-sm"
              placeholder="Password"
            />
          </div>
        </div>

        <div v-if="error" class="text-red-600 text-sm text-center">
          {{ error }}
        </div>

        <div>
          <button
            type="submit"
            :disabled="loading"
            class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-green-600 hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {{ loading ? 'Loading...' : (isLogin ? 'Sign in' : 'Sign up') }}
          </button>
        </div>

        <div class="text-center">
          <button
            type="button"
            @click="toggleMode"
            class="text-sm text-green-600 hover:text-green-500"
          >
            {{ isLogin ? "Don't have an account? Sign up" : 'Already have an account? Sign in' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '@/services/api'
import { useSocket } from '@/composables/useSocket'

const router = useRouter()
const { connect } = useSocket()

const isLogin = ref(true)
const loading = ref(false)
const error = ref('')

const formData = reactive({
  name: '',
  email: '',
  password: ''
})

const toggleMode = () => {
  isLogin.value = !isLogin.value
  error.value = ''
}

const handleSubmit = async () => {
  error.value = ''
  loading.value = true
  
  try {
    let response
    
    if (isLogin.value) {
      response = await api.login({
        username: formData.email.split('@')[0], // Use email prefix as username
        password: formData.password
      })
    } else {
      response = await api.register({
        username: formData.email.split('@')[0], // Use email prefix as username
        name: formData.name,
        email: formData.email,
        password: formData.password
      })
    }
    
    // Store user info
    localStorage.setItem('current_user', JSON.stringify(response.user || {
      id: response.userId,
      username: response.username,
      name: response.name,
      email: response.email
    }))
    
    // Connect WebSocket with token
    connect(response.token)
    
    // Small delay to ensure everything is set up
    setTimeout(() => {
      router.push('/')
    }, 100)
  } catch (err: any) {
    error.value = err.message || 'An error occurred'
  } finally {
    loading.value = false
  }
}
</script>