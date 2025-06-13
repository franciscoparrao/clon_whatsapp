<template>
  <div class="chat-window">
    <!-- Chat Header -->
    <div class="chat-header">
      <div class="chat-header-info">
        <div class="chat-avatar">
          <img v-if="chat.avatar" :src="chat.avatar" class="avatar-img" />
        </div>
        <div>
          <h3 class="chat-name">{{ getChatName(chat) }}</h3>
          <p class="chat-status">{{ getStatusText() }}</p>
        </div>
      </div>
      <div class="chat-actions">
        <button class="icon-button">
          <svg class="icon" fill="currentColor" viewBox="0 0 24 24">
            <path d="M15.5 14h-.79l-.28-.27A6.471 6.471 0 0 0 16 9.5 6.5 6.5 0 1 0 9.5 16c1.61 0 3.09-.59 4.23-1.57l.27.28v.79l5 4.99L20.49 19l-4.99-5zm-6 0C7.01 14 5 11.99 5 9.5S7.01 5 9.5 5 14 7.01 14 9.5 11.99 14 9.5 14z"/>
          </svg>
        </button>
        <button class="icon-button">
          <svg class="icon" fill="currentColor" viewBox="0 0 24 24">
            <path d="M12 8c1.1 0 2-.9 2-2s-.9-2-2-2-2 .9-2 2 .9 2 2 2zm0 2c-1.1 0-2 .9-2 2s.9 2 2 2 2-.9 2-2-.9-2-2-2zm0 6c-1.1 0-2 .9-2 2s.9 2 2 2 2-.9 2-2-.9-2-2-2z"/>
          </svg>
        </button>
      </div>
    </div>
    
    <!-- Messages Area -->
    <div class="messages-container" ref="messagesContainer">
      <div v-for="message in messages" :key="message.id">
        <MessageItem :message="message" :is-own="isOwnMessage(message)" />
      </div>
    </div>
    
    <!-- Message Input -->
    <MessageInput @send-message="handleSendMessage" />
  </div>
</template>

<script setup lang="ts">
import { ref, nextTick, watch, onMounted } from 'vue'
import MessageItem from './MessageItem.vue'
import MessageInput from './MessageInput.vue'
import type { Chat, Message } from '@/types'
import { api } from '@/services/api'

const props = defineProps<{
  chat: Chat
  messages: Message[]
}>()

const emit = defineEmits<{
  'send-message': [content: string]
}>()

const messagesContainer = ref<HTMLElement>()
const currentUserId = ref<string>('')

const getChatName = (chat: Chat): string => {
  if (chat.isGroup) {
    return chat.name || 'Group Chat'
  }
  return chat.participants[0]?.name || 'Unknown'
}

const getStatusText = (): string => {
  if (props.chat.isGroup) {
    return `${props.chat.participants.length} participants`
  }
  return 'online'
}

const isOwnMessage = (message: Message): boolean => {
  return message.senderId === currentUserId.value
}

const handleSendMessage = (content: string) => {
  emit('send-message', content)
}


// Auto-scroll to bottom when new messages arrive
watch(() => props.messages.length, async () => {
    await nextTick()
    if (messagesContainer.value) {
      messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
    }
  },
  {deep: true, immediate: true}
)

// Load current user ID
onMounted(async () => {
  try {
    const user = await api.getCurrentUser()
    currentUserId.value = user.id
  } catch (error) {
    console.error('Failed to get current user:', error)
  }
})
</script>

<style scoped>
.chat-window {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.chat-header {
  background-color: #128c7e;
  color: white;
  padding: 0.75rem 1rem;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.chat-header-info {
  display: flex;
  align-items: center;
}

.chat-avatar {
  width: 2.5rem;
  height: 2.5rem;
  border-radius: 50%;
  background-color: #d3d3d3;
  margin-right: 0.75rem;
}

.avatar-img {
  width: 100%;
  height: 100%;
  border-radius: 50%;
  object-fit: cover;
}

.chat-name {
  font-weight: 500;
}

.chat-status {
  font-size: 0.75rem;
  opacity: 0.75;
}

.chat-actions {
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

.messages-container {
  flex: 1;
  overflow-y: auto;
  background-color: #f0f2f5;
  padding: 1rem;
}
</style>