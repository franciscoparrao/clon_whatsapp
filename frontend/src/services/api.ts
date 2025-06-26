interface LoginRequest {
  username: string
  password: string
}

interface RegisterRequest {
  username: string
  name: string
  email: string
  password: string
}

interface LoginResponse {
  token: string
  user: {
    id: string
    name: string
    email: string
  }
}

interface User {
  id: string
  name: string
  email: string
  created_at: string
}

interface Chat {
  id: string
  participants: User[]
  last_message?: Message
  created_at: string
}

interface Message {
  id: string
  chat_id: string
  sender_id: string
  content: string
  created_at: string
}

interface CreateChatRequest {
  participant_ids: string[]
}

interface SendMessageRequest {
  content: string
}

class ApiService {
  private baseUrl = '/api'
  private token: string | null = null

  constructor() {
    // Initialize token from localStorage
    this.token = localStorage.getItem('auth_token')
  }

  setToken(token: string | null) {
    this.token = token
    if (token) {
      localStorage.setItem('auth_token', token)
    } else {
      localStorage.removeItem('auth_token')
    }
  }

  getToken(): string | null {
    if (!this.token) {
      this.token = localStorage.getItem('auth_token')
    }
    return this.token
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const url = `${this.baseUrl}${endpoint}`
    const headers: any = {
      'Content-Type': 'application/json',
      'ngrok-skip-browser-warning': 'true',
      ...options.headers
    }

    const token = this.getToken()
    if (token) {
      headers['Authorization'] = `Bearer ${token}`
    }

    const response = await fetch(url, {
      ...options,
      headers
    })

    if (!response.ok) {
      let errorMessage = `HTTP error! status: ${response.status}`
      try {
        const errorData = await response.json()
        errorMessage = errorData.error || errorMessage
      } catch (e) {
        // Response might not be JSON or already consumed
        console.error('Error parsing error response:', e)
      }
      console.error('API Error:', errorMessage, 'Status:', response.status)
      throw new Error(errorMessage)
    }

    return response.json()
  }

  async login(data: LoginRequest): Promise<LoginResponse> {
    const response = await this.request<LoginResponse>('/auth/login', {
      method: 'POST',
      body: JSON.stringify(data)
    })
    
    if (response.token) {
      this.setToken(response.token)
    }
    
    return response
  }

  async register(data: RegisterRequest): Promise<LoginResponse> {
    const response = await this.request<LoginResponse>('/auth/register', {
      method: 'POST',
      body: JSON.stringify(data)
    })
    
    if (response.token) {
      this.setToken(response.token)
    }
    
    return response
  }

  async logout(): Promise<void> {
    try {
      await this.request('/auth/logout', {
        method: 'POST'
      })
    } finally {
      this.setToken(null)
    }
  }

  async getChats(): Promise<Chat[]> {
    const response = await this.request<{chats: Chat[]}>('/chats')
    return response.chats || []
  }

  async createChat(data: CreateChatRequest): Promise<Chat> {
    return this.request<Chat>('/chats', {
      method: 'POST',
      body: JSON.stringify(data)
    })
  }

  async getMessages(chatId: string): Promise<Message[]> {
    return this.request<Message[]>(`/chats/${chatId}/messages`)
  }

  async sendMessage(chatId: string, data: SendMessageRequest): Promise<Message> {
    return this.request<Message>(`/chats/${chatId}/messages`, {
      method: 'POST',
      body: JSON.stringify(data)
    })
  }

  async getUsers(): Promise<User[]> {
    return this.request<User[]>('/users')
  }

  async getCurrentUser(): Promise<User> {
    return this.request<User>('/user')
  }
}

// Export a getter to always get a fresh instance with current token
export const api = new ApiService()
export type { LoginRequest, RegisterRequest, LoginResponse, User, Chat, Message, CreateChatRequest, SendMessageRequest }