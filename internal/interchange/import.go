package interchange

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"crush-session-explorer/internal/db"
)

// ImportFromAICS imports sessions from AICS format to internal database format
func ImportFromAICS(data []byte) (*Archive, error) {
	var archive Archive
	if err := json.Unmarshal(data, &archive); err != nil {
		return nil, fmt.Errorf("failed to parse AICS file: %w", err)
	}

	// Validate format version
	if archive.Version != FormatVersion {
		return nil, fmt.Errorf("unsupported AICS version: %s (expected: %s)", archive.Version, FormatVersion)
	}

	return &archive, nil
}

// ImportFromFile imports sessions from an AICS file
func ImportFromFile(filePath string) (*Archive, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return ImportFromAICS(data)
}

// ConvertToDBFormat converts AICS sessions to database format
func (a *Archive) ConvertToDBFormat() ([]db.Session, map[string][]db.ParsedMessage, error) {
	sessions := make([]db.Session, 0, len(a.Log.Sessions))
	messagesMap := make(map[string][]db.ParsedMessage)

	for _, aicsSession := range a.Log.Sessions {
		dbSession, dbMessages, err := convertAICSSession(aicsSession)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to convert session %s: %w", aicsSession.ID, err)
		}

		sessions = append(sessions, *dbSession)
		messagesMap[aicsSession.ID] = dbMessages
	}

	return sessions, messagesMap, nil
}

// convertAICSSession converts an AICS session to database format
func convertAICSSession(aicsSession Session) (*db.Session, []db.ParsedMessage, error) {
	dbSession := &db.Session{
		ID: aicsSession.ID,
	}

	// Set title
	if aicsSession.Title != "" {
		dbSession.Title = &aicsSession.Title
	}

	// Set created_at
	if aicsSession.StartedAt != nil {
		createdAtStr := aicsSession.StartedAt.Format("2006-01-02T15:04:05Z")
		dbSession.CreatedAt = &createdAtStr
	}

	// Set message count
	messageCount := len(aicsSession.Messages)
	dbSession.MessageCount = &messageCount

	// Convert messages
	dbMessages := make([]db.ParsedMessage, 0, len(aicsSession.Messages))
	for _, aicsMsg := range aicsSession.Messages {
		dbMsg, err := convertAICSMessage(aicsMsg)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to convert message %s: %w", aicsMsg.ID, err)
		}
		dbMessages = append(dbMessages, *dbMsg)
	}

	return dbSession, dbMessages, nil
}

// convertAICSMessage converts an AICS message to database format
func convertAICSMessage(aicsMsg Message) (*db.ParsedMessage, error) {
	dbMsg := &db.ParsedMessage{
		ID:   aicsMsg.ID,
		Role: aicsMsg.Role,
	}

	// Set model and provider
	if aicsMsg.Model != "" {
		dbMsg.Model = &aicsMsg.Model
	}
	if aicsMsg.Provider != "" {
		dbMsg.Provider = &aicsMsg.Provider
	}

	// Set timestamp
	if aicsMsg.Timestamp != nil {
		createdAtStr := aicsMsg.Timestamp.Format("2006-01-02T15:04:05Z")
		dbMsg.CreatedAt = &createdAtStr
	}

	// Convert content to parts
	dbMsg.Parts = make([]string, 0, len(aicsMsg.Content))
	for _, content := range aicsMsg.Content {
		if content.Text != "" {
			// Add emoji prefix based on content type
			text := content.Text
			switch content.Type {
			case "tool_call":
				if !strings.HasPrefix(text, "🔧") {
					text = "🔧 " + text
				}
			case "tool_result":
				if !strings.HasPrefix(text, "📋") {
					text = "📋 " + text
				}
			}
			dbMsg.Parts = append(dbMsg.Parts, text)
		}
	}

	return dbMsg, nil
}

// ValidateArchive performs basic validation on an AICS archive
func ValidateArchive(archive *Archive) error {
	if archive.Version == "" {
		return fmt.Errorf("missing version field")
	}

	if archive.Version != FormatVersion {
		return fmt.Errorf("unsupported version: %s (expected: %s)", archive.Version, FormatVersion)
	}

	if archive.Creator.Name == "" {
		return fmt.Errorf("missing creator name")
	}

	if len(archive.Log.Sessions) == 0 {
		return fmt.Errorf("archive contains no sessions")
	}

	// Validate each session
	for i, session := range archive.Log.Sessions {
		if session.ID == "" {
			return fmt.Errorf("session %d: missing ID", i)
		}

		if len(session.Messages) == 0 {
			return fmt.Errorf("session %s: no messages", session.ID)
		}

		// Validate each message
		for j, msg := range session.Messages {
			if msg.ID == "" {
				return fmt.Errorf("session %s, message %d: missing ID", session.ID, j)
			}
			if msg.Role == "" {
				return fmt.Errorf("session %s, message %s: missing role", session.ID, msg.ID)
			}
			if len(msg.Content) == 0 {
				return fmt.Errorf("session %s, message %s: no content", session.ID, msg.ID)
			}
		}
	}

	return nil
}
