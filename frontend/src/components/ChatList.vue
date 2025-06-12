<template>
  <div class="chat-list">
    <div
      v-for="chat in chats"
      :key="chat.id"
      @click="$emit('select-chat', chat.id)"
      class="chat-item"
    >
      <!-- Avatar -->
      <div class="chat-avatar">
        <img v-if="chat.avatar" :src="chat.avatar" class="avatar-img" />
        <div v-else class="avatar-placeholder">
          {{ getChatAvatar(chat) }}
        </div>
      </div>
      
      <!-- Chat Info -->
      <div class="chat-info">
        <div class="chat-header">
          <h3 class="chat-name">
            {{ getChatName(chat) }}
          </h3>
          <span class="chat-time">
            {{ formatTime(chat.lastMessage?.timestamp) }}
          </span>
        </div>
        <div class="chat-preview">
          <p class="chat-message">
            {{ chat.lastMessage?.content || 'No messages yet' }}
          </p>
          <span
            v-if="chat.unreadCount > 0"
            class="unread-badge"
          >
            {{ chat.unreadCount }}
          </span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Chat } from '@/types'

const props = defineProps<{
  chats: Chat[]
}>()

defineEmits<{
  'select-chat': [chatId: string]
}>()


const getChatName = (chat: Chat): string => {
  // If chat has a name, use it (this includes bot names)
  if (chat.name) {
    return chat.name
  }
  
  // For group chats without name
  if (chat.isGroup) {
    return 'Group Chat'
  }
  
  // For direct chats, try to get participant name
  if (typeof chat.participants[0] === 'object') {
    return chat.participants[0]?.name || 'Unknown'
  }
  
  return 'Unknown'
}

const formatTime = (timestamp?: Date): string => {
  if (!timestamp) return ''
  
  const now = new Date()
  const date = new Date(timestamp)
  const diff = now.getTime() - date.getTime()
  const days = Math.floor(diff / (1000 * 60 * 60 * 24))
  
  if (days === 0) {
    return date.toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit' })
  } else if (days === 1) {
    return 'Yesterday'
  } else if (days < 7) {
    return date.toLocaleDateString('en-US', { weekday: 'short' })
  } else {
    return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' })
  }
}

const getChatAvatar = (chat: Chat): string => {
  const name = getChatName(chat)
  // Check if it's a bot
  if (name.includes('Bot') || name.includes('Soporte') || name.includes('Ventas') || name.includes('Template')) {
    return '🤖'
  }
  return name.charAt(0).toUpperCase()
}
</script>

<style scoped>
.chat-list {
  flex: 1;
  overflow-y: auto;
}

.chat-item {
  display: flex;
  align-items: center;
  padding: 0.75rem 1rem;
  cursor: pointer;
  border-bottom: 1px solid #f3f4f6;
  transition: background-color 0.2s;
}

.chat-item:hover {
  background-color: #f3f4f6;
}

.chat-avatar {
  width: 3rem;
  height: 3rem;
  border-radius: 50%;
  background-color: #d3d3d3;
  flex-shrink: 0;
  margin-right: 0.75rem;
}

.avatar-img {
  width: 100%;
  height: 100%;
  border-radius: 50%;
  object-fit: cover;
}

.avatar-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.5rem;
  background-color: #25D366;
  color: white;
  border-radius: 50%;
  font-weight: bold;
}

.chat-info {
  flex: 1;
  min-width: 0;
}

.chat-header {
  display: flex;
  justify-content: space-between;
  align-items: baseline;
}

.chat-name {
  font-weight: 500;
  color: #111827;
  text-overflow: ellipsis;
  overflow: hidden;
  white-space: nowrap;
}

.chat-time {
  font-size: 0.75rem;
  color: #6b7280;
}

.chat-preview {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 0.25rem;
}

.chat-message {
  font-size: 0.875rem;
  color: #4b5563;
  text-overflow: ellipsis;
  overflow: hidden;
  white-space: nowrap;
}

.unread-badge {
  background-color: #25d366;
  color: white;
  font-size: 0.75rem;
  border-radius: 9999px;
  padding: 0.125rem 0.5rem;
  min-width: 20px;
  text-align: center;
}
</style>