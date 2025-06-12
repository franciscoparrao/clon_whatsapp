import { ref, onUnmounted } from 'vue'

let socket: WebSocket | null = null
const connected = ref(false)
const listeners = new Map<string, Set<Function>>()

export const useSocket = () => {
  const connect = (token?: string) => {
    if (socket && socket.readyState === WebSocket.OPEN) {
      return
    }

    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const wsUrl = `${protocol}//${window.location.host}/ws`
    
    socket = new WebSocket(wsUrl)
    
    socket.onopen = () => {
      connected.value = true
      console.log('Connected to WebSocket server')
      
      // Send authentication if token is provided
      if (token) {
        socket!.send(JSON.stringify({
          type: 'auth',
          token
        }))
      }
      
      emit('connect')
    }
    
    socket.onclose = () => {
      connected.value = false
      console.log('Disconnected from WebSocket server')
      emit('disconnect')
    }
    
    socket.onerror = (error) => {
      console.error('WebSocket error:', error)
      emit('error', error)
    }
    
    socket.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data)
        emit('message', data)
        
        // Emit specific event types if they exist
        if (data.type) {
          emit(data.type, data)
        }
      } catch (error) {
        console.error('Failed to parse WebSocket message:', error)
      }
    }
  }
  
  const disconnect = () => {
    if (socket) {
      socket.close()
      socket = null
    }
  }
  
  const send = (data: any) => {
    if (socket && socket.readyState === WebSocket.OPEN) {
      socket.send(JSON.stringify(data))
    } else {
      console.error('WebSocket is not connected')
    }
  }
  
  const on = (event: string, callback: Function) => {
    if (!listeners.has(event)) {
      listeners.set(event, new Set())
    }
    listeners.get(event)!.add(callback)
    
    // Return cleanup function
    return () => {
      const callbacks = listeners.get(event)
      if (callbacks) {
        callbacks.delete(callback)
      }
    }
  }
  
  const off = (event: string, callback: Function) => {
    const callbacks = listeners.get(event)
    if (callbacks) {
      callbacks.delete(callback)
    }
  }
  
  const emit = (event: string, ...args: any[]) => {
    const callbacks = listeners.get(event)
    if (callbacks) {
      callbacks.forEach(callback => {
        try {
          callback(...args)
        } catch (error) {
          console.error(`Error in event listener for ${event}:`, error)
        }
      })
    }
  }
  
  // Auto cleanup on component unmount
  onUnmounted(() => {
    listeners.clear()
  })
  
  return {
    connect,
    disconnect,
    send,
    on,
    off,
    connected,
    isConnected: () => socket?.readyState === WebSocket.OPEN
  }
}