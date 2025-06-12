package models

import (
	"time"
)

// Chat represents a conversation between users
type Chat struct {
	ID             string      `json:"id" db:"id"`
	Type           ChatType    `json:"type" db:"type"`
	Name           string      `json:"name,omitempty" db:"name"` // For group chats
	Description    string      `json:"description,omitempty" db:"description"`
	Avatar         string      `json:"avatar,omitempty" db:"avatar"`
	CreatedBy      string      `json:"createdBy" db:"created_by"`
	LastMessageID  *string     `json:"lastMessageId,omitempty" db:"last_message_id"`
	LastMessage    *Message    `json:"lastMessage,omitempty"` // Populated when fetching
	LastActivity   time.Time   `json:"lastActivity" db:"last_activity"`
	CreatedAt      time.Time   `json:"createdAt" db:"created_at"`
	UpdatedAt      time.Time   `json:"updatedAt" db:"updated_at"`
	DeletedAt      *time.Time  `json:"-" db:"deleted_at"` // Soft delete
	
	// Group-specific fields
	Settings       GroupSettings `json:"settings,omitempty" db:"settings"`
	
	// Populated fields
	Participants   []ChatParticipant `json:"participants,omitempty"`
	UnreadCount    int              `json:"unreadCount,omitempty"` // Per-user count
	IsMuted        bool             `json:"isMuted,omitempty"`     // Per-user setting
	IsPinned       bool             `json:"isPinned,omitempty"`    // Per-user setting
	IsArchived     bool             `json:"isArchived,omitempty"`  // Per-user setting
}

// ChatType represents the type of chat
type ChatType string

const (
	ChatTypePrivate   ChatType = "private"   // One-on-one chat
	ChatTypeGroup     ChatType = "group"     // Group chat
	ChatTypeBroadcast ChatType = "broadcast" // Broadcast list
	ChatTypeChannel   ChatType = "channel"   // One-way channel
)

// ChatParticipant represents a user's participation in a chat
type ChatParticipant struct {
	ID               string    `json:"id" db:"id"`
	ChatID           string    `json:"chatId" db:"chat_id"`
	UserID           string    `json:"userId" db:"user_id"`
	User             *User     `json:"user,omitempty"` // Populated when fetching
	Role             ParticipantRole `json:"role" db:"role"`
	JoinedAt         time.Time `json:"joinedAt" db:"joined_at"`
	LastReadMessageID *string  `json:"lastReadMessageId,omitempty" db:"last_read_message_id"`
	LastReadAt       *time.Time `json:"lastReadAt,omitempty" db:"last_read_at"`
	IsMuted          bool      `json:"isMuted" db:"is_muted"`
	MutedUntil       *time.Time `json:"mutedUntil,omitempty" db:"muted_until"`
	IsPinned         bool      `json:"isPinned" db:"is_pinned"`
	IsArchived       bool      `json:"isArchived" db:"is_archived"`
	NotificationSound string   `json:"notificationSound,omitempty" db:"notification_sound"`
	LeftAt           *time.Time `json:"leftAt,omitempty" db:"left_at"`
}

// ParticipantRole represents a participant's role in a chat
type ParticipantRole string

const (
	RoleMember    ParticipantRole = "member"
	RoleAdmin     ParticipantRole = "admin"
	RoleOwner     ParticipantRole = "owner"
	RoleModerator ParticipantRole = "moderator"
)

// GroupSettings represents settings for a group chat
type GroupSettings struct {
	OnlyAdminsCanSend        bool   `json:"onlyAdminsCanSend" db:"only_admins_can_send"`
	OnlyAdminsCanEditInfo    bool   `json:"onlyAdminsCanEditInfo" db:"only_admins_can_edit_info"`
	OnlyAdminsCanAddMembers  bool   `json:"onlyAdminsCanAddMembers" db:"only_admins_can_add_members"`
	AutoDeleteMessages       bool   `json:"autoDeleteMessages" db:"auto_delete_messages"`
	AutoDeleteDuration       int    `json:"autoDeleteDuration" db:"auto_delete_duration"` // In hours
	JoinByLink               bool   `json:"joinByLink" db:"join_by_link"`
	JoinLink                 string `json:"joinLink,omitempty" db:"join_link"`
	RequireAdminApproval     bool   `json:"requireAdminApproval" db:"require_admin_approval"`
	MaxParticipants          int    `json:"maxParticipants" db:"max_participants"`
}

// ChatInvite represents an invitation to join a chat
type ChatInvite struct {
	ID         string    `json:"id" db:"id"`
	ChatID     string    `json:"chatId" db:"chat_id"`
	InviteCode string    `json:"inviteCode" db:"invite_code"`
	CreatedBy  string    `json:"createdBy" db:"created_by"`
	ExpiresAt  *time.Time `json:"expiresAt,omitempty" db:"expires_at"`
	MaxUses    int       `json:"maxUses" db:"max_uses"`
	Uses       int       `json:"uses" db:"uses"`
	CreatedAt  time.Time `json:"createdAt" db:"created_at"`
}

// ChatMedia represents media shared in a chat
type ChatMedia struct {
	ID        string      `json:"id" db:"id"`
	ChatID    string      `json:"chatId" db:"chat_id"`
	MessageID string      `json:"messageId" db:"message_id"`
	Type      MessageType `json:"type" db:"type"`
	URL       string      `json:"url" db:"url"`
	ThumbnailURL string   `json:"thumbnailUrl,omitempty" db:"thumbnail_url"`
	Size      int64       `json:"size" db:"size"`
	CreatedAt time.Time   `json:"createdAt" db:"created_at"`
}

// PinnedMessage represents a pinned message in a chat
type PinnedMessage struct {
	ID        string    `json:"id" db:"id"`
	ChatID    string    `json:"chatId" db:"chat_id"`
	MessageID string    `json:"messageId" db:"message_id"`
	Message   *Message  `json:"message,omitempty"` // Populated when fetching
	PinnedBy  string    `json:"pinnedBy" db:"pinned_by"`
	PinnedAt  time.Time `json:"pinnedAt" db:"pinned_at"`
}

// GetOtherParticipant returns the other participant in a private chat
func (c *Chat) GetOtherParticipant(currentUserID string) *ChatParticipant {
	if c.Type != ChatTypePrivate {
		return nil
	}
	
	for _, participant := range c.Participants {
		if participant.UserID != currentUserID {
			return &participant
		}
	}
	
	return nil
}

// IsUserParticipant checks if a user is a participant in the chat
func (c *Chat) IsUserParticipant(userID string) bool {
	for _, participant := range c.Participants {
		if participant.UserID == userID && participant.LeftAt == nil {
			return true
		}
	}
	return false
}

// GetParticipant returns a specific participant
func (c *Chat) GetParticipant(userID string) *ChatParticipant {
	for _, participant := range c.Participants {
		if participant.UserID == userID {
			return &participant
		}
	}
	return nil
}

// CanUserSendMessage checks if a user can send messages to this chat
func (c *Chat) CanUserSendMessage(userID string) bool {
	participant := c.GetParticipant(userID)
	if participant == nil || participant.LeftAt != nil {
		return false
	}
	
	// Check group settings
	if c.Type == ChatTypeGroup && c.Settings.OnlyAdminsCanSend {
		return participant.Role == RoleAdmin || participant.Role == RoleOwner
	}
	
	return true
}

// CanUserEditInfo checks if a user can edit chat info
func (c *Chat) CanUserEditInfo(userID string) bool {
	if c.Type == ChatTypePrivate {
		return false // Private chats don't have editable info
	}
	
	participant := c.GetParticipant(userID)
	if participant == nil || participant.LeftAt != nil {
		return false
	}
	
	if c.Settings.OnlyAdminsCanEditInfo {
		return participant.Role == RoleAdmin || participant.Role == RoleOwner
	}
	
	return true
}

// CanUserAddMembers checks if a user can add members to the chat
func (c *Chat) CanUserAddMembers(userID string) bool {
	if c.Type != ChatTypeGroup {
		return false
	}
	
	participant := c.GetParticipant(userID)
	if participant == nil || participant.LeftAt != nil {
		return false
	}
	
	if c.Settings.OnlyAdminsCanAddMembers {
		return participant.Role == RoleAdmin || participant.Role == RoleOwner
	}
	
	return true
}

// DefaultGroupSettings returns default settings for a new group
func DefaultGroupSettings() GroupSettings {
	return GroupSettings{
		OnlyAdminsCanSend:       false,
		OnlyAdminsCanEditInfo:   false,
		OnlyAdminsCanAddMembers: false,
		AutoDeleteMessages:      false,
		AutoDeleteDuration:      0,
		JoinByLink:             false,
		RequireAdminApproval:   false,
		MaxParticipants:        256,
	}
}