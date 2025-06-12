package models

import (
	"time"
)

// Message represents a chat message
type Message struct {
	ID           string       `json:"id" db:"id"`
	ChatID       string       `json:"chatId" db:"chat_id"`
	SenderID     string       `json:"senderId" db:"sender_id"`
	Content      string       `json:"content" db:"content"`
	Type         MessageType  `json:"type" db:"type"`
	Status       MessageStatus `json:"status" db:"status"`
	ReplyToID    *string      `json:"replyToId,omitempty" db:"reply_to_id"`
	ReplyTo      *Message     `json:"replyTo,omitempty"` // Populated when fetching
	ForwardedFrom *string     `json:"forwardedFrom,omitempty" db:"forwarded_from"`
	EditedAt     *time.Time   `json:"editedAt,omitempty" db:"edited_at"`
	DeletedAt    *time.Time   `json:"-" db:"deleted_at"` // Soft delete
	CreatedAt    time.Time    `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time    `json:"updatedAt" db:"updated_at"`
	
	// Media-specific fields
	MediaURL     string      `json:"mediaUrl,omitempty" db:"media_url"`
	ThumbnailURL string      `json:"thumbnailUrl,omitempty" db:"thumbnail_url"`
	MediaSize    int64       `json:"mediaSize,omitempty" db:"media_size"`
	MediaDuration int        `json:"mediaDuration,omitempty" db:"media_duration"` // For audio/video in seconds
	FileName     string      `json:"fileName,omitempty" db:"file_name"`
	MimeType     string      `json:"mimeType,omitempty" db:"mime_type"`
	
	// Location-specific fields
	Latitude     float64     `json:"latitude,omitempty" db:"latitude"`
	Longitude    float64     `json:"longitude,omitempty" db:"longitude"`
	LocationName string      `json:"locationName,omitempty" db:"location_name"`
	
	// Additional metadata
	Metadata     map[string]interface{} `json:"metadata,omitempty" db:"metadata"`
	Reactions    []Reaction             `json:"reactions,omitempty"`
	ReadBy       []ReadReceipt          `json:"readBy,omitempty"`
}

// MessageType represents the type of message
type MessageType string

const (
	MessageTypeText     MessageType = "text"
	MessageTypeImage    MessageType = "image"
	MessageTypeVideo    MessageType = "video"
	MessageTypeAudio    MessageType = "audio"
	MessageTypeDocument MessageType = "document"
	MessageTypeLocation MessageType = "location"
	MessageTypeContact  MessageType = "contact"
	MessageTypeSticker  MessageType = "sticker"
	MessageTypeGIF      MessageType = "gif"
	MessageTypeSystem   MessageType = "system" // For system messages like "user joined"
)

// MessageStatus represents the delivery status of a message
type MessageStatus string

const (
	MessageStatusPending   MessageStatus = "pending"   // Message is being sent
	MessageStatusSent      MessageStatus = "sent"      // Message sent to server
	MessageStatusDelivered MessageStatus = "delivered" // Message delivered to recipient
	MessageStatusRead      MessageStatus = "read"      // Message read by recipient
	MessageStatusFailed    MessageStatus = "failed"    // Message failed to send
)

// Reaction represents a reaction to a message
type Reaction struct {
	ID        string    `json:"id" db:"id"`
	MessageID string    `json:"messageId" db:"message_id"`
	UserID    string    `json:"userId" db:"user_id"`
	Emoji     string    `json:"emoji" db:"emoji"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

// ReadReceipt represents when a message was read by a user
type ReadReceipt struct {
	ID        string    `json:"id" db:"id"`
	MessageID string    `json:"messageId" db:"message_id"`
	UserID    string    `json:"userId" db:"user_id"`
	ReadAt    time.Time `json:"readAt" db:"read_at"`
}

// MessageDelivery tracks message delivery to users
type MessageDelivery struct {
	ID          string    `json:"id" db:"id"`
	MessageID   string    `json:"messageId" db:"message_id"`
	UserID      string    `json:"userId" db:"user_id"`
	DeliveredAt time.Time `json:"deliveredAt" db:"delivered_at"`
}

// DeletedMessage tracks deleted messages for sync
type DeletedMessage struct {
	ID             string    `json:"id" db:"id"`
	MessageID      string    `json:"messageId" db:"message_id"`
	DeletedBy      string    `json:"deletedBy" db:"deleted_by"`
	DeletedAt      time.Time `json:"deletedAt" db:"deleted_at"`
	DeleteForEveryone bool   `json:"deleteForEveryone" db:"delete_for_everyone"`
}

// MessageEdit tracks message edit history
type MessageEdit struct {
	ID              string    `json:"id" db:"id"`
	MessageID       string    `json:"messageId" db:"message_id"`
	PreviousContent string    `json:"previousContent" db:"previous_content"`
	EditedBy        string    `json:"editedBy" db:"edited_by"`
	EditedAt        time.Time `json:"editedAt" db:"edited_at"`
}

// Validate checks if the message is valid
func (m *Message) Validate() error {
	// TODO: Implement validation logic
	// - Check required fields
	// - Validate message type
	// - Check content length limits
	// - Validate media URLs if present
	return nil
}

// CanEdit checks if a user can edit this message
func (m *Message) CanEdit(userID string) bool {
	// Only the sender can edit their message
	if m.SenderID != userID {
		return false
	}
	
	// Can't edit if message is deleted
	if m.DeletedAt != nil {
		return false
	}
	
	// Can't edit system messages
	if m.Type == MessageTypeSystem {
		return false
	}
	
	// Can only edit within 15 minutes of creation
	editDeadline := m.CreatedAt.Add(15 * time.Minute)
	return time.Now().Before(editDeadline)
}

// CanDelete checks if a user can delete this message
func (m *Message) CanDelete(userID string) bool {
	// User can always delete their own messages
	if m.SenderID == userID {
		return true
	}
	
	// TODO: Add admin/group admin logic
	return false
}

// IsMedia checks if the message contains media
func (m *Message) IsMedia() bool {
	switch m.Type {
	case MessageTypeImage, MessageTypeVideo, MessageTypeAudio, MessageTypeDocument, MessageTypeSticker, MessageTypeGIF:
		return true
	default:
		return false
	}
}