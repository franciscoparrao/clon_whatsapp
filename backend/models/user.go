package models

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID           string    `json:"id" db:"id"`
	Username     string    `json:"username" db:"username"`
	Email        string    `json:"email" db:"email"`
	Name         string    `json:"name" db:"name"`
	PasswordHash string    `json:"-" db:"password_hash"` // Never send password hash in JSON
	Avatar       string    `json:"avatar,omitempty" db:"avatar"`
	Bio          string    `json:"bio,omitempty" db:"bio"`
	Phone        string    `json:"phone,omitempty" db:"phone"`
	Status       string    `json:"status" db:"status"` // online, offline, away, busy
	LastSeen     time.Time `json:"lastSeen" db:"last_seen"`
	IsOnline     bool      `json:"isOnline" db:"is_online"`
	Settings     UserSettings `json:"settings" db:"settings"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time `json:"updatedAt" db:"updated_at"`
	DeletedAt    *time.Time `json:"-" db:"deleted_at"` // Soft delete
}

// UserSettings represents user preferences and settings
type UserSettings struct {
	Theme               string `json:"theme" db:"theme"`                                 // light, dark, auto
	Language            string `json:"language" db:"language"`                           // en, es, fr, etc.
	NotificationsEnabled bool   `json:"notificationsEnabled" db:"notifications_enabled"`
	SoundEnabled        bool   `json:"soundEnabled" db:"sound_enabled"`
	ShowOnlineStatus    bool   `json:"showOnlineStatus" db:"show_online_status"`
	ShowLastSeen        bool   `json:"showLastSeen" db:"show_last_seen"`
	ShowProfilePhoto    string `json:"showProfilePhoto" db:"show_profile_photo"`        // everyone, contacts, nobody
	ShowStatus          string `json:"showStatus" db:"show_status"`                      // everyone, contacts, nobody
	ReadReceipts        bool   `json:"readReceipts" db:"read_receipts"`
	TwoFactorEnabled    bool   `json:"twoFactorEnabled" db:"two_factor_enabled"`
}

// UserProfile is a public representation of a user
type UserProfile struct {
	ID       string    `json:"id"`
	Username string    `json:"username"`
	Name     string    `json:"name"`
	Avatar   string    `json:"avatar,omitempty"`
	Bio      string    `json:"bio,omitempty"`
	Status   string    `json:"status,omitempty"`
	LastSeen time.Time `json:"lastSeen,omitempty"`
	IsOnline bool      `json:"isOnline"`
}

// Contact represents a user's contact
type Contact struct {
	ID         string    `json:"id" db:"id"`
	UserID     string    `json:"userId" db:"user_id"`
	ContactID  string    `json:"contactId" db:"contact_id"`
	Name       string    `json:"name,omitempty" db:"name"`        // Custom name for the contact
	IsBlocked  bool      `json:"isBlocked" db:"is_blocked"`
	IsFavorite bool      `json:"isFavorite" db:"is_favorite"`
	AddedAt    time.Time `json:"addedAt" db:"added_at"`
	UpdatedAt  time.Time `json:"updatedAt" db:"updated_at"`
}

// BlockedUser represents a blocked user
type BlockedUser struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"userId" db:"user_id"`
	BlockedID string    `json:"blockedId" db:"blocked_id"`
	BlockedAt time.Time `json:"blockedAt" db:"blocked_at"`
}

// ToProfile converts a User to UserProfile
func (u *User) ToProfile() *UserProfile {
	profile := &UserProfile{
		ID:       u.ID,
		Username: u.Username,
		Name:     u.Name,
		Avatar:   u.Avatar,
		Bio:      u.Bio,
		IsOnline: u.IsOnline,
	}

	// Only show status and last seen based on privacy settings
	if u.Settings.ShowOnlineStatus {
		profile.Status = u.Status
	}
	if u.Settings.ShowLastSeen {
		profile.LastSeen = u.LastSeen
	}

	return profile
}

// CanViewProfile checks if a user can view another user's profile details
func (u *User) CanViewProfile(viewerID string, isContact bool) bool {
	// User can always view their own profile
	if u.ID == viewerID {
		return true
	}

	// Check profile photo visibility settings
	switch u.Settings.ShowProfilePhoto {
	case "everyone":
		return true
	case "contacts":
		return isContact
	case "nobody":
		return false
	default:
		return true
	}
}

// CanViewStatus checks if a user can view another user's status
func (u *User) CanViewStatus(viewerID string, isContact bool) bool {
	// User can always view their own status
	if u.ID == viewerID {
		return true
	}

	// Check status visibility settings
	switch u.Settings.ShowStatus {
	case "everyone":
		return true
	case "contacts":
		return isContact
	case "nobody":
		return false
	default:
		return true
	}
}

// DefaultUserSettings returns default settings for a new user
func DefaultUserSettings() UserSettings {
	return UserSettings{
		Theme:               "light",
		Language:            "en",
		NotificationsEnabled: true,
		SoundEnabled:        true,
		ShowOnlineStatus:    true,
		ShowLastSeen:        true,
		ShowProfilePhoto:    "everyone",
		ShowStatus:          "everyone",
		ReadReceipts:        true,
		TwoFactorEnabled:    false,
	}
}