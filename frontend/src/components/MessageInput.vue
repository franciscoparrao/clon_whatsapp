<template>
  <div class="input-container">
    <!-- Emoji Button -->
    <button class="icon-button">
      <svg class="icon" fill="currentColor" viewBox="0 0 24 24">
        <path d="M11.99 2C6.47 2 2 6.48 2 12s4.47 10 9.99 10C17.52 22 22 17.52 22 12S17.52 2 11.99 2zM12 20c-4.42 0-8-3.58-8-8s3.58-8 8-8 8 3.58 8 8-3.58 8-8 8zm3.5-9c.83 0 1.5-.67 1.5-1.5S16.33 8 15.5 8 14 8.67 14 9.5s.67 1.5 1.5 1.5zm-7 0c.83 0 1.5-.67 1.5-1.5S9.33 8 8.5 8 7 8.67 7 9.5 7.67 11 8.5 11zm3.5 6.5c2.33 0 4.31-1.46 5.11-3.5H6.89c.8 2.04 2.78 3.5 5.11 3.5z"/>
      </svg>
    </button>
    
    <!-- Attachment Button -->
    <button class="icon-button">
      <svg class="icon" fill="currentColor" viewBox="0 0 24 24">
        <path d="M16.5 6v11.5c0 2.21-1.79 4-4 4s-4-1.79-4-4V5c0-1.38 1.12-2.5 2.5-2.5s2.5 1.12 2.5 2.5v10.5c0 .55-.45 1-1 1s-1-.45-1-1V6H10v9.5c0 1.38 1.12 2.5 2.5 2.5s2.5-1.12 2.5-2.5V5c0-2.21-1.79-4-4-4S7 2.79 7 5v12.5c0 3.04 2.46 5.5 5.5 5.5s5.5-2.46 5.5-5.5V6h-1.5z"/>
      </svg>
    </button>
    
    <!-- Message Input -->
    <input
      v-model="messageText"
      @keyup.enter="sendMessage"
      type="text"
      placeholder="Type a message"
      class="message-input"
    />
    
    <!-- Send/Voice Button -->
    <button
      @click="sendMessage"
      :class="['icon-button', messageText ? 'send-active' : '']"
    >
      <svg v-if="messageText" class="icon" fill="currentColor" viewBox="0 0 24 24">
        <path d="M2.01 21L23 12 2.01 3 2 10l15 2-15 2z"/>
      </svg>
      <svg v-else class="icon" fill="currentColor" viewBox="0 0 24 24">
        <path d="M12 15c1.66 0 2.99-1.34 2.99-3L15 5c0-1.66-1.34-3-3-3S9 3.34 9 5v7c0 1.66 1.34 3 3 3zm5.3-3c0 3-2.54 5.1-5.3 5.1S6.7 15 6.7 12H5c0 3.42 2.72 6.23 6 6.72V22h2v-3.28c3.28-.48 6-3.3 6-6.72h-1.7z"/>
      </svg>
    </button>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const emit = defineEmits<{
  'send-message': [content: string]
}>()

const messageText = ref('')

const sendMessage = () => {
  if (messageText.value.trim()) {
    emit('send-message', messageText.value)
    messageText.value = ''
  }
}
</script>

<style scoped>
.input-container {
  background-color: #f3f4f6;
  padding: 0.75rem 1rem;
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.icon-button {
  color: #6b7280;
  background: none;
  border: none;
  cursor: pointer;
  transition: color 0.2s;
}

.icon-button:hover {
  color: #374151;
}

.icon-button.send-active {
  color: #25d366;
}

.icon-button.send-active:hover {
  color: #128c7e;
}

.icon {
  width: 1.5rem;
  height: 1.5rem;
}

.message-input {
  flex: 1;
  padding: 0.5rem 1rem;
  border-radius: 9999px;
  background-color: white;
  border: 1px solid #e5e7eb;
  outline: none;
  transition: border-color 0.2s;
}

.message-input:focus {
  border-color: #25d366;
}
</style>