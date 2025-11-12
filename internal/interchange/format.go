package interchange

import (
	"time"
)

// AICS (AI Coding Session) Format
// A standard interchange format for AI coding sessions across different providers
// Inspired by HAR (HTTP Archive) format for standardization

const (
	// FormatVersion is the current version of the AICS format
	FormatVersion = "1.0"
	// FormatName is the standard name for this format
	FormatName = "AICS"
)

// Archive represents the root structure of an AICS file
type Archive struct {
	Version string   `json:"version"`           // Format version (e.g., "1.0")
	Creator Creator  `json:"creator"`           // Information about the tool that created this archive
	Browser *Browser `json:"browser,omitempty"` // Information about the AI tool/provider
	Log     Log      `json:"log"`               // Container for sessions
}

// Creator describes the tool that created the archive
type Creator struct {
	Name    string `json:"name"`              // Name of the export tool
	Version string `json:"version"`           // Version of the export tool
	Comment string `json:"comment,omitempty"` // Additional information
}

// Browser describes the AI coding tool/provider
type Browser struct {
	Name    string `json:"name"`              // Name of the AI tool (e.g., "Cursor", "Claude Code", "Crush")
	Version string `json:"version,omitempty"` // Version of the tool
	Comment string `json:"comment,omitempty"` // Additional information
}

// Log contains all sessions and metadata
type Log struct {
	Version  string    `json:"version"`           // Log version
	Creator  Creator   `json:"creator"`           // Creator information
	Browser  *Browser  `json:"browser,omitempty"` // Browser/tool information
	Sessions []Session `json:"sessions"`          // Array of sessions
	Comment  string    `json:"comment,omitempty"` // Additional information
}

// Session represents a single AI coding session
type Session struct {
	ID        string     `json:"id"`                  // Unique session identifier (UUID v7)
	ClientID  string     `json:"clientId,omitempty"`  // Client/machine identifier
	Title     string     `json:"title,omitempty"`     // Session title
	StartedAt *time.Time `json:"startedAt,omitempty"` // When the session started (ISO 8601)
	UpdatedAt *time.Time `json:"updatedAt,omitempty"` // When the session was last updated
	Messages  []Message  `json:"messages"`            // Array of messages in chronological order
	Metadata  Metadata   `json:"metadata,omitempty"`  // Additional session metadata
	Comment   string     `json:"comment,omitempty"`   // Additional information
}

// Message represents a single message in a session
type Message struct {
	ID        string     `json:"id"`                  // Unique message identifier
	Timestamp *time.Time `json:"timestamp,omitempty"` // When the message was created (ISO 8601)
	Role      string     `json:"role"`                // Message role: "user", "assistant", "system", "tool"
	Content   []Content  `json:"content"`             // Message content parts
	Model     string     `json:"model,omitempty"`     // AI model used (e.g., "claude-3-opus", "gpt-4")
	Provider  string     `json:"provider,omitempty"`  // Provider name (e.g., "anthropic", "openai")
	Metadata  Metadata   `json:"metadata,omitempty"`  // Additional message metadata
	Comment   string     `json:"comment,omitempty"`   // Additional information
}

// Content represents a content part within a message
type Content struct {
	Type     string   `json:"type"`               // Content type: "text", "tool_call", "tool_result", "code", "image"
	Text     string   `json:"text,omitempty"`     // Text content
	Data     Metadata `json:"data,omitempty"`     // Structured data for tool calls, results, etc.
	MimeType string   `json:"mimeType,omitempty"` // MIME type for binary/encoded content
	Encoding string   `json:"encoding,omitempty"` // Encoding for binary content (e.g., "base64")
	Comment  string   `json:"comment,omitempty"`  // Additional information
}

// Metadata represents flexible key-value metadata
type Metadata map[string]interface{}

// ProviderInfo contains information about the original provider
type ProviderInfo struct {
	Name         string   `json:"name"`                   // Provider name
	Version      string   `json:"version,omitempty"`      // Provider version
	DatabaseType string   `json:"databaseType,omitempty"` // Type of database used
	ExportedAt   string   `json:"exportedAt"`             // Export timestamp (ISO 8601)
	Metadata     Metadata `json:"metadata,omitempty"`     // Additional provider metadata
}
