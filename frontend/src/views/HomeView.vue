<template>
  <div class="home-container">
    <!-- Sidebar -->
    <div class="sidebar">
      <!-- Header -->
      <div class="header">
        <div class="header-profile">
          <div class="avatar"></div>
          <span class="username">My Account</span>
        </div>
        <div class="header-actions">
          <button class="icon-button" @click="logout" title="Logout">
            <svg class="icon" fill="currentColor" viewBox="0 0 24 24">
              <path d="M17 7l-1.41 1.41L18.17 11H8v2h10.17l-2.58 2.58L17 17l5-5zM4 5h8V3H4c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h8v-2H4V5z"/>
            </svg>
          </button>
        </div>
      </div>
      
      <!-- Search -->
      <div class="search-container">
        <input
          type="text"
          placeholder="Search or start new chat"
          class="search-input"
          v-model="searchQuery"
        />
        <button
          @click="showNewChatModal = true"
          class="new-chat-button"
          title="New Chat"
        >
          <svg class="icon" fill="currentColor" viewBox="0 0 24 24">
            <path d="M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z"/>
          </svg>
        </button>
      </div>
      
      <!-- Chat List -->
      <ChatList :chats="chats" @select-chat="selectChat" />
    </div>
    
    <!-- Main Chat Area -->
    <div class="main-area">
      <ChatWindow 
        v-if="selectedChat"
        :chat="selectedChat"
        :messages="messages"
        @send-message="sendMessage"
      />
      <div v-else class="welcome-screen">
        <div class="welcome-content">
          <h2 class="welcome-title">WhatsApp Web Clone</h2>
          <p class="welcome-subtitle">Select a chat to start messaging</p>
        </div>
      </div>
    </div>
    
    <!-- New Chat Modal -->
    <div v-if="showNewChatModal" class="modal-overlay" @click="showNewChatModal = false">
      <div class="modal-content" @click.stop>
        <h3 class="modal-title">Start New Chat</h3>
        <div class="modal-body">
          <select v-model="selectedUserId" class="user-select">
            <option value="">Select a user</option>
            <option v-for="user in users" :key="user.id" :value="user.id">
              {{ user.name }} ({{ user.email }})
            </option>
          </select>
        </div>
        <div class="modal-actions">
          <button @click="showNewChatModal = false" class="btn-cancel">
            Cancel
          </button>
          <button @click="createNewChat" :disabled="!selectedUserId" class="btn-primary">
            Start Chat
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import ChatList from '@/components/ChatList.vue'
import ChatWindow from '@/components/ChatWindow.vue'
import type { Chat, Message } from '@/types'
import { useSocket } from '@/composables/useSocket'
import { api } from '@/services/api'

const route = useRoute()
const router = useRouter()
const { send, on, connect, connected } = useSocket()

const chats = ref<Chat[]>([])
const messages = ref<Message[]>([])
const selectedChatId = ref<string | null>(null)
const searchQuery = ref('')
const showNewChatModal = ref(false)
const users = ref<any[]>([])
const selectedUserId = ref('')

const selectedChat = computed(() => 
  chats.value.find(chat => chat.id === selectedChatId.value)
)

const selectChat = async (chatId: string) => {
  selectedChatId.value = chatId
  // Load messages for this chat
  await loadMessages(chatId)
}

const loadMessages = async (chatId: string) => {
  try {
    const data = await api.getMessages(chatId)
    messages.value = data.map(msg => ({
      ...msg,
      timestamp: new Date(msg.createdAt),
      status: 'read' as const,
      type: 'text' as const
    }))
  } catch (error) {
    console.error('Failed to load messages:', error)
  }
}

const sendMessage = async (content: string) => {
  if (!selectedChatId.value || !content.trim()) return
  
  try {
    const response = await api.sendMessage(selectedChatId.value, {
      content: content.trim()
    })
    
    const message: Message = {
      ...response,
      timestamp: new Date(response.created_at),
      status: 'sent' as const,
      type: 'text' as const
    }
    
    messages.value = [...messages.value, message]
    
    // Send via WebSocket for real-time updates
    send({
      type: 'message',
      chatId: selectedChatId.value,
      content: content.trim()
    })
  } catch (error) {
    console.error('Failed to send message:', error)
  }
}

const loadChats = async () => {
  try {
    const chatsData = await api.getChats()
    
    chats.value = chatsData.map(chat => {
      // Map last_message to lastMessage for consistency
      const mappedChat = {
        ...chat,
        unreadCount: 0,
        isGroup: chat.participants?.length > 2,
        lastMessage: undefined as any
      }
      
      // Handle both last_message and lastMessage
      const lastMsg = chat.last_message || chat.lastMessage
      if (lastMsg) {
        mappedChat.lastMessage = {
          ...lastMsg,
          timestamp: new Date(lastMsg.created_at || lastMsg.timestamp),
          status: 'read' as const,
          type: 'text' as const,
          content: lastMsg.content
        }
      }
      
      return mappedChat
    })
  } catch (error) {
    console.error('Failed to load chats:', error)
  }
}

const logout = async () => {
  try {
    await api.logout()
    router.push('/login')
  } catch (error) {
    console.error('Failed to logout:', error)
  }
}

const loadUsers = async () => {
  try {
    const allUsers = await api.getUsers()
    // For now, use the stored user info instead of calling getCurrentUser
    const storedUser = localStorage.getItem('current_user')
    if (storedUser) {
      const currentUser = JSON.parse(storedUser)
      users.value = allUsers.filter(user => user.id !== currentUser.id)
    } else {
      users.value = allUsers
    }
  } catch (error) {
    console.error('Failed to load users:', error)
  }
}

const createNewChat = async () => {
  if (!selectedUserId.value) return
  
  try {
    const chat = await api.createChat({
      participant_ids: [selectedUserId.value]
    })
    
    // Add to chats list
    chats.value.unshift({
      ...chat,
      unreadCount: 0,
      isGroup: chat.participants.length > 2,
      lastMessage: undefined
    })
    
    // Select the new chat
    selectedChatId.value = chat.id
    await loadMessages(chat.id)
    
    // Close modal and reset
    showNewChatModal.value = false
    selectedUserId.value = ''
  } catch (error) {
    console.error('Failed to create chat:', error)
  }
}

// Socket event listeners
onMounted(async () => {
  // Connect WebSocket
  const token = api.getToken()
  if (token) {
    connect(token)
  }
  
  // Load initial data
  await loadChats()
  await loadUsers()
  
  if (route.params.id) {
    await selectChat(route.params.id as string)
  }
  
  // Set up WebSocket listeners
  on('new-message', (data: any) => {
    const message: Message = {
      id: data.id,
      chatId: data.chat_id,
      senderId: data.sender_id,
      content: data.content,
      timestamp: new Date(data.created_at),
      status: 'delivered' as const,
      type: 'text' as const
    }
    
    if (selectedChatId.value === data.chat_id) {
      messages.value = [...messages.value, message]
    }
    
    // Update last message in chat list
    const chat = chats.value.find(c => c.id === data.chat_id)
    if (chat) {
      chat.lastMessage = message
    }
  })
  
  on('message-status-update', ({ messageId, status }: any) => {
    const message = messages.value.find(m => m.id === messageId)
    if (message) {
      message.status = status
    }
  })
})
</script>

<style scoped>
.home-container {
  display: flex;
  height: 100vh;
  background-color: #f0f2f5;
}

.sidebar {
  width: 33.333333%;
  border-right: 1px solid #d3d3d3;
  background-color: white;
  display: flex;
  flex-direction: column;
}

.header {
  background-color: #128c7e;
  color: white;
  padding: 1rem;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.header-profile {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.avatar {
  width: 2.5rem;
  height: 2.5rem;
  border-radius: 50%;
  background-color: #d3d3d3;
}

.username {
  font-weight: 500;
}

.header-actions {
  display: flex;
  gap: 1rem;
}

.icon-button {
  background: none;
  border: none;
  cursor: pointer;
  transition: opacity 0.2s;
}

.icon-button:hover {
  opacity: 0.8;
}

.icon {
  width: 1.5rem;
  height: 1.5rem;
}

.search-container {
  padding: 0.5rem;
  background-color: #f9fafb;
  display: flex;
  gap: 0.5rem;
  align-items: center;
}

.search-input {
  flex: 1;
  padding: 0.5rem 1rem;
  border-radius: 9999px;
  background-color: white;
  border: 1px solid #e5e7eb;
  outline: none;
  transition: border-color 0.2s;
}

.search-input:focus {
  border-color: #25d366;
}

.new-chat-button {
  background-color: #25d366;
  color: white;
  border: none;
  border-radius: 50%;
  width: 2.5rem;
  height: 2.5rem;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: background-color 0.2s;
}

.new-chat-button:hover {
  background-color: #128c7e;
}

.main-area {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.welcome-screen {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f0f2f5;
}

.welcome-content {
  text-align: center;
}

.welcome-title {
  font-size: 1.5rem;
  color: #4b5563;
  margin-bottom: 0.5rem;
}

.welcome-subtitle {
  color: #6b7280;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  background-color: white;
  border-radius: 0.5rem;
  padding: 1.5rem;
  width: 90%;
  max-width: 400px;
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1);
}

.modal-title {
  font-size: 1.25rem;
  font-weight: 600;
  margin-bottom: 1rem;
}

.modal-body {
  margin-bottom: 1.5rem;
}

.user-select {
  width: 100%;
  padding: 0.5rem;
  border: 1px solid #e5e7eb;
  border-radius: 0.375rem;
  font-size: 1rem;
  outline: none;
  transition: border-color 0.2s;
}

.user-select:focus {
  border-color: #25d366;
}

.modal-actions {
  display: flex;
  gap: 0.75rem;
  justify-content: flex-end;
}

.btn-cancel {
  padding: 0.5rem 1rem;
  border: 1px solid #e5e7eb;
  border-radius: 0.375rem;
  background-color: white;
  color: #374151;
  cursor: pointer;
  transition: background-color 0.2s;
}

.btn-cancel:hover {
  background-color: #f9fafb;
}

.btn-primary {
  padding: 0.5rem 1rem;
  border: none;
  border-radius: 0.375rem;
  background-color: #25d366;
  color: white;
  cursor: pointer;
  transition: background-color 0.2s;
}

.btn-primary:hover {
  background-color: #128c7e;
}

.btn-primary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>