package markdown

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"crush-session-explorer/internal/db"
)

// yamlEscape escapes strings for YAML frontmatter
func yamlEscape(s string) string {
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\"", "'")
	return s
}

// FormatTimestamp formats a timestamp for display
func FormatTimestamp(ts *string) string {
	if ts == nil || *ts == "" {
		return ""
	}

	// Try parsing as Unix timestamp
	if timestamp, err := strconv.ParseInt(*ts, 10, 64); err == nil {
		return time.Unix(timestamp, 0).Format("2006-01-02 15:04")
	}

	// Try parsing as ISO format
	if t, err := time.Parse(time.RFC3339, *ts); err == nil {
		return t.Local().Format("2006-01-02 15:04")
	}

	// Return as-is if parsing fails
	return *ts
}

// formatTimestampISO formats a timestamp as ISO string for frontmatter
func formatTimestampISO(ts *string) string {
	if ts == nil || *ts == "" {
		return ""
	}

	// Try parsing as Unix timestamp
	if timestamp, err := strconv.ParseInt(*ts, 10, 64); err == nil {
		return time.Unix(timestamp, 0).Format(time.RFC3339)
	}

	// Try parsing as ISO format and return normalized
	if t, err := time.Parse(time.RFC3339, *ts); err == nil {
		return t.Format(time.RFC3339)
	}

	// Return as-is if parsing fails
	return *ts
}

// slugify converts text to a URL-friendly slug
func slugify(text string) string {
	text = strings.ToLower(strings.TrimSpace(text))

	// Remove non-alphanumeric characters except hyphens, spaces, and underscores
	reg := regexp.MustCompile(`[^a-z0-9\-\s_]+`)
	text = reg.ReplaceAllString(text, "")

	// Replace spaces and underscores with hyphens
	reg = regexp.MustCompile(`[\s_]+`)
	text = reg.ReplaceAllString(text, "-")

	if text == "" {
		return "untitled"
	}

	return text
}

// RenderMarkdown converts a session and messages to markdown format
func RenderMarkdown(session *db.Session, messages []db.ParsedMessage) string {
	var result strings.Builder

	// Generate title
	title := "Session " + session.ID
	if session.Title != nil && *session.Title != "" {
		title = *session.Title
	}

	// Generate frontmatter
	result.WriteString("---\n")
	result.WriteString(fmt.Sprintf("title: \"%s\"\n", yamlEscape(title)))
	result.WriteString(fmt.Sprintf("session_id: %s\n", session.ID))

	if session.CreatedAt != nil {
		if iso := formatTimestampISO(session.CreatedAt); iso != "" {
			result.WriteString(fmt.Sprintf("created_at: %s\n", iso))
		}
	}

	if session.MessageCount != nil {
		result.WriteString(fmt.Sprintf("message_count: %d\n", *session.MessageCount))
	}

	// Add metadata if present
	if session.Metadata != nil && *session.Metadata != "" {
		var metadata map[string]interface{}
		if err := json.Unmarshal([]byte(*session.Metadata), &metadata); err == nil {
			result.WriteString("metadata:\n")
			for k, v := range metadata {
				jsonValue, _ := json.Marshal(v)
				result.WriteString(fmt.Sprintf("  %s: %s\n", k, string(jsonValue)))
			}
		}
	}

	result.WriteString("---\n\n")

	// Generate message content
	for _, msg := range messages {
		// Generate header
		role := msg.Role
		ts := FormatTimestamp(msg.CreatedAt)
		header := fmt.Sprintf("## %s — %s", role, ts)

		// Add model/provider info if available
		var modelInfo []string
		if msg.Model != nil && *msg.Model != "" {
			modelInfo = append(modelInfo, *msg.Model)
		}
		if msg.Provider != nil && *msg.Provider != "" {
			modelInfo = append(modelInfo, *msg.Provider)
		}
		if len(modelInfo) > 0 {
			header += fmt.Sprintf(" (%s)", strings.Join(modelInfo, "/"))
		}

		result.WriteString(header + "\n\n")

		// Add message content
		result.WriteString("<div>\n")
		for _, part := range msg.Parts {
			result.WriteString(part + "\n")
		}
		result.WriteString("</div>\n\n")
	}

	return result.String()
}

// GenerateFilename generates a filename for the session
func GenerateFilename(session *db.Session) string {
	// Generate base name from title or session ID
	base := slugify("session-" + session.ID[:8])
	if session.Title != nil && *session.Title != "" {
		base = slugify(*session.Title)
	}

	// Generate timestamp prefix
	prefix := time.Now().Format("2006-01-02_15-04")
	if session.CreatedAt != nil {
		if timestamp, err := strconv.ParseInt(*session.CreatedAt, 10, 64); err == nil {
			prefix = time.Unix(timestamp, 0).Format("2006-01-02_15-04")
		}
	}

	return fmt.Sprintf("%s_%s.md", prefix, base)
}
