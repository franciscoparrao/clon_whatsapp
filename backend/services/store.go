package services

import (
	"fmt"
	"sync"
	"time"
)

type User struct {
	ID       string    `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Name     string    `json:"name"`
	Password string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

type Chat struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Participants []string  `json:"participants"`
	CreatedAt    time.Time `json:"created_at"`
}

type Message struct {
	ID        string    `json:"id"`
	ChatID    string    `json:"chat_id"`
	SenderID  string    `json:"sender_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// MemoryStore is a simple in-memory store for demo purposes
type MemoryStore struct {
	users    map[string]*User
	chats    map[string]*Chat
	messages map[string][]*Message
	mu       sync.RWMutex
}

// Store se inicializa en store_init.go
// var Store = &MemoryStore{
// 	users:    make(map[string]*User),
// 	chats:    make(map[string]*Chat),
// 	messages: make(map[string][]*Message),
// }

// User methods
func (s *MemoryStore) CreateUser(user *User) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.users[user.ID] = user
	return nil
}

func (s *MemoryStore) GetUserByEmail(email string) (*User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	for _, user := range s.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, nil
}

func (s *MemoryStore) GetUserByUsername(username string) (*User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	for _, user := range s.users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, nil
}

func (s *MemoryStore) GetUserByID(id string) (*User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	user, exists := s.users[id]
	if !exists {
		return nil, nil
	}
	return user, nil
}

func (s *MemoryStore) GetUser(id string) (*User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	user, ok := s.users[id]
	if !ok {
		return nil, nil
	}
	return user, nil
}

func (s *MemoryStore) GetUsers() ([]*User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	users := make([]*User, 0, len(s.users))
	for _, user := range s.users {
		users = append(users, user)
	}
	return users, nil
}

// Chat methods
func (s *MemoryStore) CreateChat(chat *Chat) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.chats[chat.ID] = chat
	return nil
}

func (s *MemoryStore) GetChat(id string) (*Chat, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	chat, ok := s.chats[id]
	if !ok {
		return nil, nil
	}
	return chat, nil
}

func (s *MemoryStore) GetUserChats(userID string) ([]*Chat, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	fmt.Printf("GetUserChats: Looking for chats for user %s\n", userID)
	fmt.Printf("GetUserChats: Total chats in store: %d\n", len(s.chats))
	
	chats := make([]*Chat, 0)
	for _, chat := range s.chats {
		for _, participant := range chat.Participants {
			if participant == userID {
				fmt.Printf("GetUserChats: Found chat %s for user\n", chat.ID)
				chats = append(chats, chat)
				break
			}
		}
	}
	fmt.Printf("GetUserChats: Returning %d chats for user %s\n", len(chats), userID)
	return chats, nil
}

// Message methods
func (s *MemoryStore) CreateMessage(message *Message) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.messages[message.ChatID] = append(s.messages[message.ChatID], message)
	return nil
}

func (s *MemoryStore) GetChatMessages(chatID string) ([]*Message, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	messages, ok := s.messages[chatID]
	if !ok {
		return []*Message{}, nil
	}
	return messages, nil
}

// Helper to find or create a direct chat between two users
func (s *MemoryStore) FindOrCreateDirectChat(user1ID, user2ID string) (*Chat, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// Check if chat already exists
	for _, chat := range s.chats {
		if len(chat.Participants) == 2 {
			hasUser1 := false
			hasUser2 := false
			for _, p := range chat.Participants {
				if p == user1ID {
					hasUser1 = true
				}
				if p == user2ID {
					hasUser2 = true
				}
			}
			if hasUser1 && hasUser2 {
				return chat, nil
			}
		}
	}
	
	// Create new chat
	chat := &Chat{
		ID:           "chat_" + time.Now().Format("20060102150405"),
		Name:         "Direct Chat",
		Participants: []string{user1ID, user2ID},
		CreatedAt:    time.Now(),
	}
	s.chats[chat.ID] = chat
	return chat, nil
}

func (s *MemoryStore) IsUserInChat(userID, chatID string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	chat, exists := s.chats[chatID]
	if !exists {
		return false
	}
	
	for _, participant := range chat.Participants {
		if participant == userID {
			return true
		}
	}
	return false
}