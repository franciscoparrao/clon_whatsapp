<template>
  <div :class="['message-wrapper', isOwn ? 'own-message' : 'other-message']">
    <div
      :class="['message-bubble', isOwn ? 'own-bubble' : 'other-bubble']"
    >
      <p class="message-text">{{ message.content }}</p>
      <div class="message-meta">
        <span class="message-time">
          {{ formatTime(message.timestamp) }}
        </span>
        <span v-if="isOwn">
          <svg
            v-if="message.status === 'sent'"
            class="status-icon sent"
            fill="currentColor"
            viewBox="0 0 24 24"
          >
            <path d="M9 16.17L4.83 12l-1.42 1.41L9 19 21 7l-1.41-1.41z"/>
          </svg>
          <svg
            v-else-if="message.status === 'delivered'"
            class="status-icon delivered"
            fill="currentColor"
            viewBox="0 0 24 24"
          >
            <path d="M18 7l-1.41-1.41-6.34 6.34 1.41 1.41L18 7zm4.24-1.41L11.66 16.17 7.48 12l-1.41 1.41L11.66 19l12-12-1.42-1.41zM.41 13.41L6 19l1.41-1.41L1.83 12 .41 13.41z"/>
          </svg>
          <svg
            v-else-if="message.status === 'read'"
            class="status-icon read"
            fill="currentColor"
            viewBox="0 0 24 24"
          >
            <path d="M18 7l-1.41-1.41-6.34 6.34 1.41 1.41L18 7zm4.24-1.41L11.66 16.17 7.48 12l-1.41 1.41L11.66 19l12-12-1.42-1.41zM.41 13.41L6 19l1.41-1.41L1.83 12 .41 13.41z"/>
          </svg>
        </span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Message } from '@/types'

defineProps<{
  message: Message
  isOwn: boolean
}>()

const formatTime = (timestamp: Date): string => {
  const date = new Date(timestamp)
  return date.toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit' })
}
</script>

<style scoped>
.message-wrapper {
  display: flex;
  margin-bottom: 0.5rem;
}

.own-message {
  justify-content: flex-end;
}

.other-message {
  justify-content: flex-start;
}

.message-bubble {
  max-width: 20rem;
  padding: 0.5rem 1rem;
  border-radius: 0.5rem;
}

@media (min-width: 1024px) {
  .message-bubble {
    max-width: 28rem;
  }
}

.own-bubble {
  background-color: #dcf8c6;
}

.other-bubble {
  background-color: white;
}

.message-text {
  font-size: 0.875rem;
  color: #111827;
}

.message-meta {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  margin-top: 0.25rem;
  gap: 0.25rem;
}

.message-time {
  font-size: 0.75rem;
  color: #6b7280;
}

.status-icon {
  width: 1rem;
  height: 1rem;
}

.status-icon.sent,
.status-icon.delivered {
  color: #6b7280;
}

.status-icon.read {
  color: #3b82f6;
}
</style>