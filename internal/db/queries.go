package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
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
		var id, role string
		var partsJSON *string
		var model, provider, createdAt *string

		err := rows.Scan(&id, &role, &partsJSON, &model, &provider, &createdAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}

		// Parse the JSON parts
		parsed := ParsedMessage{
			ID:        id,
			Role:      role,
			Model:     model,
			Provider:  provider,
			CreatedAt: createdAt,
		}

		// Parse parts JSON
		if partsJSON != nil && *partsJSON != "" {
			var rawParts []interface{}
			if err := json.Unmarshal([]byte(*partsJSON), &rawParts); err == nil {
				for _, part := range rawParts {
					switch p := part.(type) {
					case string:
						if strings.TrimSpace(p) != "" {
							parsed.Parts = append(parsed.Parts, p)
						}
					case map[string]interface{}:
						// Handle different message types
						if msgType, ok := p["type"].(string); ok {
							switch msgType {
							case "text":
								// Handle text messages
								if data, ok := p["data"].(map[string]interface{}); ok {
									if text, ok := data["text"].(string); ok && strings.TrimSpace(text) != "" {
										parsed.Parts = append(parsed.Parts, text)
									}
								}
							case "tool_call":
								// Handle tool calls - show what tool was called
								if data, ok := p["data"].(map[string]interface{}); ok {
									if name, ok := data["name"].(string); ok {
										toolInfo := fmt.Sprintf("🔧 Tool call: %s", name)
										if input, ok := data["input"].(string); ok && len(input) < 200 {
											toolInfo += fmt.Sprintf("\nInput: %s", input)
										}
										parsed.Parts = append(parsed.Parts, toolInfo)
									}
								}
							case "tool_result":
								// Handle tool results - show the result
								if data, ok := p["data"].(map[string]interface{}); ok {
									if content, ok := data["content"].(string); ok && strings.TrimSpace(content) != "" {
										result := fmt.Sprintf("📋 Tool result:\n%s", content)
										parsed.Parts = append(parsed.Parts, result)
									}
								}
							case "finish":
								// Skip finish messages as they don't contain user content
								continue
							}
						} else {
							// Fallback for old format
							if textData, ok := p["text"]; ok {
								if text, ok := textData.(string); ok && strings.TrimSpace(text) != "" {
									parsed.Parts = append(parsed.Parts, text)
								}
							} else if data, ok := p["data"].(map[string]interface{}); ok {
								if text, ok := data["text"].(string); ok && strings.TrimSpace(text) != "" {
									parsed.Parts = append(parsed.Parts, text)
								}
							}
						}
					}
				}
			}
		}

		// Only add message if it has actual content
		if len(parsed.Parts) > 0 {
			messages = append(messages, parsed)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating messages: %w", err)
	}

	return messages, nil
}
