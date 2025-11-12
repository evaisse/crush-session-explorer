package interchange

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"crush-session-explorer/internal/db"

	"github.com/google/uuid"
)

const toolVersion = "v1.0.1"

// ExportToAICS exports sessions to the AICS (AI Coding Session) format
func ExportToAICS(sessions []db.Session, messages map[string][]db.ParsedMessage, providerName string) (*Archive, error) {
	archive := &Archive{
		Version: FormatVersion,
		Creator: Creator{
			Name:    "crush-session-explorer",
			Version: toolVersion,
			Comment: "Exported from Crush database",
		},
		Browser: &Browser{
			Name:    providerName,
			Comment: "Original AI coding tool",
		},
		Log: Log{
			Version: FormatVersion,
			Creator: Creator{
				Name:    "crush-session-explorer",
				Version: toolVersion,
			},
			Browser: &Browser{
				Name: providerName,
			},
			Sessions: make([]Session, 0, len(sessions)),
		},
	}

	// Convert each session
	for _, dbSession := range sessions {
		session, err := convertSession(dbSession, messages[dbSession.ID])
		if err != nil {
			return nil, fmt.Errorf("failed to convert session %s: %w", dbSession.ID, err)
		}
		archive.Log.Sessions = append(archive.Log.Sessions, *session)
	}

	return archive, nil
}

// convertSession converts a database session to AICS format
func convertSession(dbSession db.Session, dbMessages []db.ParsedMessage) (*Session, error) {
	session := &Session{
		ID:       dbSession.ID,
		Messages: make([]Message, 0, len(dbMessages)),
		Metadata: make(Metadata),
	}

	// Set title
	if dbSession.Title != nil {
		session.Title = *dbSession.Title
	}

	// Parse created_at timestamp
	if dbSession.CreatedAt != nil {
		createdAt := parseTimestamp(*dbSession.CreatedAt)
		session.StartedAt = createdAt
		session.UpdatedAt = createdAt
	}

	// Add message count to metadata
	if dbSession.MessageCount != nil {
		session.Metadata["message_count"] = *dbSession.MessageCount
	}

	// Convert messages
	for _, dbMsg := range dbMessages {
		msg, err := convertMessage(dbMsg)
		if err != nil {
			return nil, fmt.Errorf("failed to convert message %s: %w", dbMsg.ID, err)
		}
		session.Messages = append(session.Messages, *msg)

		// Update session updated_at to latest message timestamp
		if msg.Timestamp != nil && (session.UpdatedAt == nil || msg.Timestamp.After(*session.UpdatedAt)) {
			session.UpdatedAt = msg.Timestamp
		}
	}

	return session, nil
}

// convertMessage converts a database message to AICS format
func convertMessage(dbMsg db.ParsedMessage) (*Message, error) {
	msg := &Message{
		ID:       dbMsg.ID,
		Role:     dbMsg.Role,
		Content:  make([]Content, 0, len(dbMsg.Parts)),
		Metadata: make(Metadata),
	}

	// Set model and provider
	if dbMsg.Model != nil {
		msg.Model = *dbMsg.Model
	}
	if dbMsg.Provider != nil {
		msg.Provider = *dbMsg.Provider
	}

	// Parse timestamp
	if dbMsg.CreatedAt != nil {
		msg.Timestamp = parseTimestamp(*dbMsg.CreatedAt)
	}

	// Convert message parts to content
	for _, part := range dbMsg.Parts {
		content := Content{
			Type: "text",
			Text: part,
		}

		// Detect tool calls and results based on content patterns
		if len(part) > 0 {
			if strings.HasPrefix(part, "🔧") {
				content.Type = "tool_call"
			} else if strings.HasPrefix(part, "📋") {
				content.Type = "tool_result"
			}
		}

		msg.Content = append(msg.Content, content)
	}

	return msg, nil
}

// parseTimestamp attempts to parse various timestamp formats
func parseTimestamp(ts string) *time.Time {
	if ts == "" {
		return nil
	}

	// Try RFC3339 format first (ISO 8601)
	if t, err := time.Parse(time.RFC3339, ts); err == nil {
		return &t
	}

	// Try RFC3339Nano format
	if t, err := time.Parse(time.RFC3339Nano, ts); err == nil {
		return &t
	}

	// Try common date-time formats
	formats := []string{
		"2006-01-02T15:04:05Z",
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05",
		time.RFC1123,
		time.RFC822,
	}

	for _, format := range formats {
		if t, err := time.Parse(format, ts); err == nil {
			return &t
		}
	}

	return nil
}

// ToJSON converts the archive to JSON format
func (a *Archive) ToJSON() ([]byte, error) {
	return json.MarshalIndent(a, "", "  ")
}

// ToJSONCompact converts the archive to compact JSON format
func (a *Archive) ToJSONCompact() ([]byte, error) {
	return json.Marshal(a)
}

// GenerateSessionID generates a UUID v7 for a session
func GenerateSessionID() string {
	id, err := uuid.NewV7()
	if err != nil {
		// Fallback to UUID v4 if v7 generation fails
		return uuid.NewString()
	}
	return id.String()
}

// GetClientID retrieves or generates a persistent client ID
func GetClientID() (string, error) {
	// Try to get from environment first
	if clientID := os.Getenv("CRUSH_CLIENT_ID"); clientID != "" {
		return clientID, nil
	}

	// Try to read from config file
	configDir, err := os.UserConfigDir()
	if err != nil {
		// Fall back to generating a new one
		return uuid.NewString(), nil
	}

	clientIDPath := filepath.Join(configDir, "crush-session-explorer", "client-id")

	// Try to read existing client ID
	if data, err := os.ReadFile(clientIDPath); err == nil {
		return strings.TrimSpace(string(data)), nil
	}

	// Generate new client ID (using v4 for client ID is fine)
	clientID := uuid.NewString()

	// Try to save it for future use
	if err := os.MkdirAll(filepath.Dir(clientIDPath), 0755); err == nil {
		_ = os.WriteFile(clientIDPath, []byte(clientID), 0644)
	}

	return clientID, nil
}

// ExportSessionToFile exports a single session to a file in a date-based folder structure
func ExportSessionToFile(session *Session, baseDir string, providerName string) (string, error) {
	if session.StartedAt == nil {
		return "", fmt.Errorf("session has no start time")
	}

	// Create folder structure: baseDir/YYYY/MM/DD/
	year := session.StartedAt.Format("2006")
	month := session.StartedAt.Format("01")
	day := session.StartedAt.Format("02")

	sessionDir := filepath.Join(baseDir, year, month, day)
	if err := os.MkdirAll(sessionDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create session directory: %w", err)
	}

	// Create a single-session archive
	archive := &Archive{
		Version: FormatVersion,
		Creator: Creator{
			Name:    "crush-session-explorer",
			Version: toolVersion,
			Comment: "Exported from Crush database",
		},
		Browser: &Browser{
			Name:    providerName,
			Comment: "Original AI coding tool",
		},
		Log: Log{
			Version: FormatVersion,
			Creator: Creator{
				Name:    "crush-session-explorer",
				Version: toolVersion,
			},
			Browser: &Browser{
				Name: providerName,
			},
			Sessions: []Session{*session},
		},
	}

	// Convert to JSON
	jsonData, err := archive.ToJSON()
	if err != nil {
		return "", fmt.Errorf("failed to convert to JSON: %w", err)
	}

	// Create filename based on session ID
	filename := fmt.Sprintf("%s.aics.json", session.ID)
	filePath := filepath.Join(sessionDir, filename)

	// Write file
	if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return filePath, nil
}

// ExportSessionsIndividually exports each session to its own file in a date-based folder structure
func ExportSessionsIndividually(sessions []db.Session, messages map[string][]db.ParsedMessage, baseDir string, providerName string, clientID string) ([]string, error) {
	var exportedFiles []string

	for _, dbSession := range sessions {
		// Convert session
		session, err := convertSession(dbSession, messages[dbSession.ID])
		if err != nil {
			return exportedFiles, fmt.Errorf("failed to convert session %s: %w", dbSession.ID, err)
		}

		// Generate new UUID v7 for the session
		session.ID = GenerateSessionID()

		// Set client ID
		if clientID != "" {
			session.ClientID = clientID
		}

		// Export to file
		filePath, err := ExportSessionToFile(session, baseDir, providerName)
		if err != nil {
			return exportedFiles, fmt.Errorf("failed to export session %s: %w", session.ID, err)
		}

		exportedFiles = append(exportedFiles, filePath)
	}

	return exportedFiles, nil
}
