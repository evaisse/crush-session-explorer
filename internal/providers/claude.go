package providers

import (
	"crush-session-explorer/internal/db"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// ClaudeProvider implements the Provider interface for Claude Code/Desktop sessions
type ClaudeProvider struct {
	dbPath string
	conn   *sql.DB
}

// NewClaudeProvider creates a new Claude provider instance
func NewClaudeProvider() *ClaudeProvider {
	provider := &ClaudeProvider{}
	provider.dbPath = provider.getDefaultDBPath()
	return provider
}

// NewClaudeProviderWithPath creates a new Claude provider with custom db path
func NewClaudeProviderWithPath(dbPath string) *ClaudeProvider {
	return &ClaudeProvider{
		dbPath: dbPath,
	}
}

// Name returns the provider name
func (p *ClaudeProvider) Name() string {
	return "claude-code"
}

// GetDBPath returns the database path
func (p *ClaudeProvider) GetDBPath() string {
	return p.dbPath
}

// getDefaultDBPath returns the default Claude database path based on OS
func (p *ClaudeProvider) getDefaultDBPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	switch runtime.GOOS {
	case "darwin":
		// macOS: ~/Library/Application Support/Claude/
		return filepath.Join(home, "Library", "Application Support", "Claude", "state.db")
	case "windows":
		// Windows: %APPDATA%/Claude/
		appData := os.Getenv("APPDATA")
		if appData == "" {
			appData = filepath.Join(home, "AppData", "Roaming")
		}
		return filepath.Join(appData, "Claude", "state.db")
	case "linux":
		// Linux: ~/.config/Claude/
		return filepath.Join(home, ".config", "Claude", "state.db")
	default:
		return ""
	}
}

// Discover checks if Claude database exists
func (p *ClaudeProvider) Discover() (bool, error) {
	if p.dbPath == "" {
		return false, nil
	}

	// Check if database file exists
	if _, err := os.Stat(p.dbPath); os.IsNotExist(err) {
		return false, nil
	}

	// Try to connect and check for expected tables
	conn, err := sql.Open("sqlite3", p.dbPath)
	if err != nil {
		return false, nil
	}
	defer conn.Close()

	// Check if conversations table exists
	var tableName string
	err = conn.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='conversations'").Scan(&tableName)
	if err != nil {
		return false, nil
	}

	return true, nil
}

// getConnection returns or creates a database connection
func (p *ClaudeProvider) getConnection() (*sql.DB, error) {
	if p.conn == nil {
		conn, err := sql.Open("sqlite3", p.dbPath)
		if err != nil {
			return nil, err
		}
		if err := conn.Ping(); err != nil {
			conn.Close()
			return nil, err
		}
		p.conn = conn
	}
	return p.conn, nil
}

// ListSessions retrieves sessions from Claude database
func (p *ClaudeProvider) ListSessions(limit int) ([]db.Session, error) {
	conn, err := p.getConnection()
	if err != nil {
		return nil, err
	}

	// Claude Desktop schema typically has: conversations table
	// Try to query with common column names
	query := `
		SELECT 
			uuid,
			COALESCE(name, ''),
			created_at,
			updated_at
		FROM conversations
		ORDER BY updated_at DESC
		LIMIT ?
	`

	rows, err := conn.Query(query, limit)
	if err != nil {
		// If query fails, it might be a different schema - return empty
		return nil, fmt.Errorf("failed to query Claude conversations: %w", err)
	}
	defer rows.Close()

	var sessions []db.Session
	for rows.Next() {
		var id, name, createdAt, updatedAt string

		err := rows.Scan(&id, &name, &createdAt, &updatedAt)
		if err != nil {
			continue
		}

		// Format the session to match our standard Session model
		title := name
		if title == "" {
			title = "Untitled Conversation"
		}

		// Use updated_at as the primary timestamp
		timestamp := updatedAt
		if timestamp == "" {
			timestamp = createdAt
		}

		// Count messages for this conversation
		var msgCount int
		msgQuery := "SELECT COUNT(*) FROM chat_messages WHERE conversation_uuid = ?"
		if err := conn.QueryRow(msgQuery, id).Scan(&msgCount); err != nil {
			msgCount = 0
		}

		provider := "claude-code"
		sessions = append(sessions, db.Session{
			ID:           id,
			Title:        &title,
			CreatedAt:    &timestamp,
			Metadata:     &provider,
			MessageCount: &msgCount,
		})
	}

	return sessions, nil
}

// FetchSession retrieves a specific session
func (p *ClaudeProvider) FetchSession(sessionID string) (*db.Session, error) {
	conn, err := p.getConnection()
	if err != nil {
		return nil, err
	}

	query := `
		SELECT 
			uuid,
			COALESCE(name, ''),
			created_at,
			updated_at
		FROM conversations
		WHERE uuid = ?
	`

	var id, name, createdAt, updatedAt string
	err = conn.QueryRow(query, sessionID).Scan(&id, &name, &createdAt, &updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("session not found: %s", sessionID)
		}
		return nil, fmt.Errorf("failed to fetch Claude session: %w", err)
	}

	title := name
	if title == "" {
		title = "Untitled Conversation"
	}

	timestamp := updatedAt
	if timestamp == "" {
		timestamp = createdAt
	}

	// Count messages
	var msgCount int
	msgQuery := "SELECT COUNT(*) FROM chat_messages WHERE conversation_uuid = ?"
	if err := conn.QueryRow(msgQuery, id).Scan(&msgCount); err != nil {
		msgCount = 0
	}

	provider := "claude-code"
	return &db.Session{
		ID:           id,
		Title:        &title,
		CreatedAt:    &timestamp,
		Metadata:     &provider,
		MessageCount: &msgCount,
	}, nil
}

// ListMessages retrieves messages for a session
func (p *ClaudeProvider) ListMessages(sessionID string) ([]db.ParsedMessage, error) {
	conn, err := p.getConnection()
	if err != nil {
		return nil, err
	}

	query := `
		SELECT 
			uuid,
			COALESCE(sender, ''),
			COALESCE(text, ''),
			created_at
		FROM chat_messages
		WHERE conversation_uuid = ?
		ORDER BY created_at ASC
	`

	rows, err := conn.Query(query, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to query Claude messages: %w", err)
	}
	defer rows.Close()

	var messages []db.ParsedMessage
	for rows.Next() {
		var id, sender, text, createdAt string

		err := rows.Scan(&id, &sender, &text, &createdAt)
		if err != nil {
			continue
		}

		// Map Claude sender types to standard roles
		role := "user"
		if strings.ToLower(sender) == "assistant" || strings.ToLower(sender) == "claude" {
			role = "assistant"
		}

		// Parse timestamp to RFC3339 format if needed
		timestamp := p.normalizeTimestamp(createdAt)

		model := "claude"
		provider := "anthropic"

		messages = append(messages, db.ParsedMessage{
			ID:        id,
			Role:      role,
			Parts:     []string{text},
			Model:     &model,
			Provider:  &provider,
			CreatedAt: &timestamp,
		})
	}

	return messages, nil
}

// normalizeTimestamp converts various timestamp formats to RFC3339
func (p *ClaudeProvider) normalizeTimestamp(ts string) string {
	if ts == "" {
		return time.Now().Format(time.RFC3339)
	}

	// Try parsing as RFC3339 first
	if t, err := time.Parse(time.RFC3339, ts); err == nil {
		return t.Format(time.RFC3339)
	}

	// Try parsing as Unix timestamp (milliseconds)
	if t, err := time.Parse("2006-01-02T15:04:05.999Z", ts); err == nil {
		return t.Format(time.RFC3339)
	}

	// Try other common formats
	formats := []string{
		time.RFC3339Nano,
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05",
		"2006-01-02",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, ts); err == nil {
			return t.Format(time.RFC3339)
		}
	}

	// If all else fails, return original
	return ts
}

// SetDBPath allows setting a custom database path
func (p *ClaudeProvider) SetDBPath(path string) {
	// Expand home directory if needed
	if len(path) > 0 && path[0] == '~' {
		home, err := os.UserHomeDir()
		if err == nil {
			path = filepath.Join(home, path[1:])
		}
	}
	p.dbPath = path
	// Close existing connection since path changed
	if p.conn != nil {
		p.conn.Close()
		p.conn = nil
	}
}

// Close closes the database connection
func (p *ClaudeProvider) Close() error {
	if p.conn != nil {
		err := p.conn.Close()
		p.conn = nil
		return err
	}
	return nil
}

// Helper function to parse JSON content if needed
func parseJSONContent(content string) []string {
	var parts []string

	// Try to parse as JSON array
	var jsonArray []interface{}
	if err := json.Unmarshal([]byte(content), &jsonArray); err == nil {
		for _, item := range jsonArray {
			if str, ok := item.(string); ok {
				parts = append(parts, str)
			} else if obj, ok := item.(map[string]interface{}); ok {
				if text, ok := obj["text"].(string); ok {
					parts = append(parts, text)
				}
			}
		}
		return parts
	}

	// If not JSON, return as single part
	return []string{content}
}
