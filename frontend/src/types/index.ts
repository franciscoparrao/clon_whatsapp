export interface User {
  id: string
  name: string
  avatar?: string
  status?: string
  lastSeen?: Date
}

export interface Message {
  id: string
  chatId: string
  senderId: string
  content: string
  timestamp: Date
  status: 'sent' | 'delivered' | 'read'
  type: 'text' | 'image' | 'audio' | 'video'
}

export interface Chat {
  id: string
  participants: string[] | User[]
  lastMessage?: Message
  last_message?: any
  unreadCount: number
  isGroup: boolean
  name?: string
  avatar?: string
  created_at?: string
  updated_at?: string
}