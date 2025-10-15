package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

// ListSessions retrieves sessions from the database with a limit
func ListSessions(db *sql.DB, limit int) ([]Session, error) {
	query := `
		SELECT id, title, created_at, message_count
		FROM sessions
		ORDER BY created_at DESC
		LIMIT ?
	`

	rows, err := db.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query sessions: %w", err)
	}
	defer rows.Close()

	var sessions []Session
	for rows.Next() {
		var s Session
		err := rows.Scan(&s.ID, &s.Title, &s.CreatedAt, &s.MessageCount)
		if err != nil {
			return nil, fmt.Errorf("failed to scan session: %w", err)
		}
		sessions = append(sessions, s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating sessions: %w", err)
	}

	return sessions, nil
}

// FetchSession retrieves a specific session by ID
func FetchSession(db *sql.DB, sessionID string) (*Session, error) {
	query := `
		SELECT id, title, created_at, message_count
		FROM sessions
		WHERE id = ?
	`

	var s Session
	err := db.QueryRow(query, sessionID).Scan(&s.ID, &s.Title, &s.CreatedAt, &s.MessageCount)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("session not found: %s", sessionID)
		}
		return nil, fmt.Errorf("failed to fetch session: %w", err)
	}

	return &s, nil
}

// ListMessages retrieves all messages for a session
func ListMessages(db *sql.DB, sessionID string) ([]ParsedMessage, error) {
	query := `
		SELECT id, role, parts, model, provider, created_at
		FROM messages
		WHERE session_id = ?
		ORDER BY created_at ASC
	`

	rows, err := db.Query(query, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to query messages: %w", err)
	}
	defer rows.Close()

	var messages []ParsedMessage
	for rows.Next() {
		var m Message
		err := rows.Scan(&m.ID, &m.Role, &m.Parts, &m.Model, &m.Provider, &m.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}

		// Parse the JSON parts
		parsed := ParsedMessage{
			ID:        m.ID,
			Role:      m.Role,
			Model:     m.Model,
			Provider:  m.Provider,
			CreatedAt: m.CreatedAt,
		}

		// Parse parts JSON
		if m.Parts != nil {
			var rawParts []interface{}
			if err := json.Unmarshal(m.Parts, &rawParts); err == nil {
				for _, part := range rawParts {
					switch p := part.(type) {
					case string:
						parsed.Parts = append(parsed.Parts, p)
					case map[string]interface{}:
						// Handle object parts with text or data.text
						if textData, ok := p["text"]; ok {
							if text, ok := textData.(string); ok {
								parsed.Parts = append(parsed.Parts, text)
							}
						} else if data, ok := p["data"].(map[string]interface{}); ok {
							if text, ok := data["text"].(string); ok {
								parsed.Parts = append(parsed.Parts, text)
							}
						}
					}
				}
			}
		}

		messages = append(messages, parsed)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating messages: %w", err)
	}

	return messages, nil
}