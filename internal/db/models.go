package db

import (
	"encoding/json"
	"time"
)

// Session represents a chat session from the database
type Session struct {
	ID           string    `json:"id"`
	Title        *string   `json:"title"`
	CreatedAt    *string   `json:"created_at"`
	Metadata     *string   `json:"metadata"`
	Content      *string   `json:"content"`
	MessageCount *int      `json:"message_count"`
}

// Message represents a message within a session
type Message struct {
	ID        string          `json:"id"`
	Role      string          `json:"role"`
	Parts     json.RawMessage `json:"parts"`
	Model     *string         `json:"model"`
	Provider  *string         `json:"provider"`
	CreatedAt *string         `json:"created_at"`
}

// ParsedMessage represents a message with parsed parts
type ParsedMessage struct {
	ID        string   `json:"id"`
	Role      string   `json:"role"`
	Parts     []string `json:"parts"`
	Model     *string  `json:"model"`
	Provider  *string  `json:"provider"`
	CreatedAt *string  `json:"created_at"`
}

// ParsedCreatedAt returns the created_at timestamp as a time.Time
func (s *Session) ParsedCreatedAt() *time.Time {
	if s.CreatedAt == nil {
		return nil
	}
	
	// Try parsing as Unix timestamp first
	if t, err := time.Parse("1", *s.CreatedAt); err == nil {
		return &t
	}
	
	// Try parsing as ISO format
	if t, err := time.Parse(time.RFC3339, *s.CreatedAt); err == nil {
		return &t
	}
	
	return nil
}

// ParsedCreatedAt returns the created_at timestamp as a time.Time
func (m *Message) ParsedCreatedAt() *time.Time {
	if m.CreatedAt == nil {
		return nil
	}
	
	// Try parsing as Unix timestamp first
	if t, err := time.Parse("1", *m.CreatedAt); err == nil {
		return &t
	}
	
	// Try parsing as ISO format
	if t, err := time.Parse(time.RFC3339, *m.CreatedAt); err == nil {
		return &t
	}
	
	return nil
}