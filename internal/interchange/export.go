package interchange

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"crush-session-explorer/internal/db"
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
