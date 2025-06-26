package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
	
	"whatsapp-clone/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ChatHandler struct {
	hub *Hub
}

type CreateChatRequest struct {
	Name         string   `json:"name"`
	Type         string   `json:"type"` // private, group
	Participants []string `json:"participants"`
}

type SendMessageRequest struct {
	Content     string `json:"content"`
	Type        string `json:"type"` // text, image, video, audio, document
	ReplyTo     string `json:"replyTo,omitempty"`
	MediaURL    string `json:"mediaUrl,omitempty"`
	ThumbnailURL string `json:"thumbnailUrl,omitempty"`
}

type ChatResponse struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Type         string    `json:"type"`
	Participants []string  `json:"participants"`
	LastMessage  *MessageResponse `json:"lastMessage,omitempty"`
	UnreadCount  int       `json:"unreadCount"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type MessageResponse struct {
	ID          string    `json:"id"`
	ChatID      string    `json:"chatId"`
	SenderID    string    `json:"senderId"`
	Content     string    `json:"content"`
	Type        string    `json:"type"`
	Status      string    `json:"status"` // sent, delivered, read
	ReplyTo     string    `json:"replyTo,omitempty"`
	MediaURL    string    `json:"mediaUrl,omitempty"`
	ThumbnailURL string   `json:"thumbnailUrl,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func NewChatHandler() *ChatHandler {
	return &ChatHandler{}
}

// GetChats returns all chats for the authenticated user
func (h *ChatHandler) GetChats(c *gin.Context) {
	userID := c.GetString("userID")
	c.Writer.Header().Set("X-Debug-UserID", userID) // Debug header
	
	if userID == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No user ID in context"})
		return
	}

	// Get all chats for the user from the store
	chats, err := services.Store.GetUserChats(userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch chats"})
		return
	}
	
	// Debug log
	log.Printf("GetChats: Found %d chats for user %s", len(chats), userID)

	// Convert to response format
	chatResponses := make([]ChatResponse, 0, len(chats))
	for _, chat := range chats {
		// Get messages for this chat to find the last message
		messages, _ := services.Store.GetChatMessages(chat.ID)
		
		var lastMessage *MessageResponse
		if len(messages) > 0 {
			lastMsg := messages[len(messages)-1]
			lastMessage = &MessageResponse{
				ID:        lastMsg.ID,
				ChatID:    lastMsg.ChatID,
				SenderID:  lastMsg.SenderID,
				Content:   lastMsg.Content,
				Type:      "text",
				Status:    "delivered",
				CreatedAt: lastMsg.CreatedAt,
				UpdatedAt: lastMsg.CreatedAt,
			}
		}

		// For private chats, set the name to the other participant's name
		chatName := chat.Name
		if chatName == "Direct Chat" && len(chat.Participants) == 2 {
			for _, participantID := range chat.Participants {
				if participantID != userID {
					if user, _ := services.Store.GetUser(participantID); user != nil {
						chatName = user.Name
					}
					break
				}
			}
		}

		chatResponses = append(chatResponses, ChatResponse{
			ID:           chat.ID,
			Name:         chatName,
			Type:         "private",
			Participants: chat.Participants,
			LastMessage:  lastMessage,
			UnreadCount:  0, // TODO: Implement unread count
			CreatedAt:    chat.CreatedAt,
			UpdatedAt:    chat.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"chats": chatResponses,
	})
}

// GetChat returns a specific chat by ID
func (h *ChatHandler) GetChat(c *gin.Context) {
	chatID := c.Param("chatId")
	userID := c.GetString("userID")

	// TODO: Fetch chat from database and verify user is participant
	// For now, return mock data
	chat := ChatResponse{
		ID:   chatID,
		Name: "John Doe",
		Type: "private",
		Participants: []string{userID, "user2"},
		CreatedAt:    time.Now().Add(-24 * time.Hour),
		UpdatedAt:    time.Now().Add(-5 * time.Minute),
	}

	c.JSON(http.StatusOK, chat)
}

// CreateChat creates a new chat
func (h *ChatHandler) CreateChat(c *gin.Context) {
	userID := c.GetString("userID")
	
	var req CreateChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Validate request
	if len(req.Participants) == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "At least one participant is required"})
		return
	}

	// Add current user to participants if not already included
	hasCurrentUser := false
	for _, p := range req.Participants {
		if p == userID {
			hasCurrentUser = true
			break
		}
	}
	if !hasCurrentUser {
		req.Participants = append(req.Participants, userID)
	}

	var chat *services.Chat
	var err error

	// If only one participant (besides current user), create or find direct chat
	if len(req.Participants) == 2 {
		// Find the other participant
		var otherUserID string
		for _, p := range req.Participants {
			if p != userID {
				otherUserID = p
				break
			}
		}
		chat, err = services.Store.FindOrCreateDirectChat(userID, otherUserID)
	} else {
		// Create group chat
		chat = &services.Chat{
			ID:           generateChatID(),
			Name:         req.Name,
			Participants: req.Participants,
			CreatedAt:    time.Now(),
		}
		err = services.Store.CreateChat(chat)
	}

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create chat"})
		return
	}

	// Get name for private chats
	chatName := chat.Name
	if chatName == "Direct Chat" && len(chat.Participants) == 2 {
		for _, participantID := range chat.Participants {
			if participantID != userID {
				if user, _ := services.Store.GetUser(participantID); user != nil {
					chatName = user.Name
				}
				break
			}
		}
	}

	chatResponse := ChatResponse{
		ID:           chat.ID,
		Name:         chatName,
		Type:         "private",
		Participants: chat.Participants,
		CreatedAt:    chat.CreatedAt,
		UpdatedAt:    chat.CreatedAt,
	}

	c.JSON(http.StatusCreated, chatResponse)
}

// GetMessages returns messages for a specific chat
func (h *ChatHandler) GetMessages(c *gin.Context) {
	chatID := c.Param("id")
	userID := c.GetString("userID")

	// Get query parameters for pagination
	limitStr := c.DefaultQuery("limit", "50")
	limit, _ := strconv.Atoi(limitStr)
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	// Verify user is participant of the chat
	chat, err := services.Store.GetChat(chatID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch chat"})
		return
	}

	if chat == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Chat not found"})
		return
	}

	// Check if user is participant
	isParticipant := false
	for _, p := range chat.Participants {
		if p == userID {
			isParticipant = true
			break
		}
	}

	if !isParticipant {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "You are not a participant of this chat"})
		return
	}

	// Fetch messages from store
	messages, err := services.Store.GetChatMessages(chatID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch messages"})
		return
	}

	// Convert to response format and apply limit
	messageResponses := make([]MessageResponse, 0)
	start := 0
	if len(messages) > limit {
		start = len(messages) - limit
	}

	for i := start; i < len(messages); i++ {
		msg := messages[i]
		messageResponses = append(messageResponses, MessageResponse{
			ID:        msg.ID,
			ChatID:    msg.ChatID,
			SenderID:  msg.SenderID,
			Content:   msg.Content,
			Type:      "text",
			Status:    "delivered",
			CreatedAt: msg.CreatedAt,
			UpdatedAt: msg.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, messageResponses)
}

// SendMessage sends a new message to a chat
func (h *ChatHandler) SendMessage(c *gin.Context) {
	chatID := c.Param("id")
	userID := c.GetString("userID")

	var req SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Validate request
	if req.Content == "" && req.MediaURL == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Message content or media is required"})
		return
	}

	// Default to text type if not specified
	if req.Type == "" {
		req.Type = "text"
	}

	// Verify user is participant of the chat
	chat, err := services.Store.GetChat(chatID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch chat"})
		return
	}

	if chat == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Chat not found"})
		return
	}

	// Check if user is participant
	isParticipant := false
	for _, p := range chat.Participants {
		if p == userID {
			isParticipant = true
			break
		}
	}

	if !isParticipant {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "You are not a participant of this chat"})
		return
	}

	// Create message in store
	message := &services.Message{
		ID:        generateMessageID(),
		ChatID:    chatID,
		SenderID:  userID,
		Content:   req.Content,
		CreatedAt: time.Now(),
	}

	err = services.Store.CreateMessage(message)
	if err != nil {
		// Log the actual error for debugging
		c.Writer.Header().Set("X-Debug-Error", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to save message", "details": err.Error()})
		return
	}

	// Check if any participant is a bot and generate response
	for _, participantID := range chat.Participants {
		if participantID != userID {
			// Check if participant is a bot by checking their email
			participantUser, _ := services.Store.GetUser(participantID)
			if participantUser != nil {
				// Check if this is a bot email
				isBotEmail := participantUser.Email == "support@whatsapp-clone.com" ||
					participantUser.Email == "sales@whatsapp-clone.com" ||
					participantUser.Email == "template@whatsapp-clone.com"
				
				if isBotEmail {
					// Map email to bot ID for processing
					var botID string
					switch participantUser.Email {
					case "support@whatsapp-clone.com":
						botID = "bot_support"
					case "sales@whatsapp-clone.com":
						botID = "bot_sales"
					case "template@whatsapp-clone.com":
						botID = "bot_template"
					}
					
					if botID != "" {
						// Generate bot response asynchronously
						go services.GetBotResponse(botID, req.Content, chatID, h.hub)
					}
				}
			}
		}
	}

	// Create response
	messageResponse := MessageResponse{
		ID:           message.ID,
		ChatID:       message.ChatID,
		SenderID:     message.SenderID,
		Content:      message.Content,
		Type:         req.Type,
		Status:       "sent",
		ReplyTo:      req.ReplyTo,
		MediaURL:     req.MediaURL,
		ThumbnailURL: req.ThumbnailURL,
		CreatedAt:    message.CreatedAt,
		UpdatedAt:    message.CreatedAt,
	}

	// Broadcast message through WebSocket
	if h.hub != nil {
		// Create WebSocket message
		wsMessage := Message{
			Type:      "chat_message",
			UserID:    userID,
			ChatID:    chatID,
			Content:   req.Content,
			Timestamp: time.Now().Unix(),
		}

		// Add message data
		msgData, _ := json.Marshal(messageResponse)
		wsMessage.Data = msgData

		// Broadcast to all connected clients
		broadcastData, _ := json.Marshal(wsMessage)
		h.hub.broadcast <- broadcastData
	}

	c.JSON(http.StatusCreated, messageResponse)
}

// MarkAsRead marks messages as read
func (h *ChatHandler) MarkAsRead(c *gin.Context) {
	chatID := c.Param("chatId")
	userID := c.GetString("userID")

	var messageIDs []string
	if err := c.ShouldBindJSON(&messageIDs); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// TODO: Update message status in database for chatID and userID
	// TODO: Send read receipts through WebSocket
	_ = chatID // Will be used to identify the chat
	_ = userID // Will be used to verify user permissions

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// DeleteMessage deletes a message
func (h *ChatHandler) DeleteMessage(c *gin.Context) {
	messageID := c.Param("messageId")
	userID := c.GetString("userID")

	// TODO: Verify user is the sender of the message using messageID and userID
	// TODO: Delete or mark message as deleted in database
	// TODO: Notify other participants through WebSocket
	_ = messageID // Will be used to identify the message to delete
	_ = userID    // Will be used to verify ownership

	c.Status(http.StatusNoContent)
}

func generateChatID() string {
	return "chat_" + uuid.New().String()
}

// GetUsers returns all users except the current user
func (h *ChatHandler) GetUsers(c *gin.Context) {
	currentUserID := c.GetString("userID")

	// Get all users from store
	users, err := services.Store.GetUsers()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	// Filter out current user and prepare response
	type UserResponse struct {
		ID       string `json:"id"`
		Username string `json:"username"`
		Name     string `json:"name"`
		Email    string `json:"email"`
	}

	userResponses := make([]UserResponse, 0)
	for _, user := range users {
		if user.ID != currentUserID {
			userResponses = append(userResponses, UserResponse{
				ID:       user.ID,
				Username: user.Username,
				Name:     user.Name,
				Email:    user.Email,
			})
		}
	}

	c.JSON(http.StatusOK, userResponses)
}

func generateMessageID() string {
	return uuid.New().String()
}