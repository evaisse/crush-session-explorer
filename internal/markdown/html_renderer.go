package markdown

import (
	"encoding/json"
	"fmt"
	"html"
	"strconv"
	"strings"
	"time"

	"crush-session-explorer/internal/db"
)

// getRoleEmoji returns the appropriate emoji for a role
func getRoleEmoji(role string) string {
	switch strings.ToLower(role) {
	case "user":
		return "👤"
	case "assistant":
		return "🤖"
	case "tool":
		return "🔧"
	case "system":
		return "⚙️"
	default:
		return "💬"
	}
}

// formatTimeOnly formats a timestamp to show only time (HH:MM:SS)
func formatTimeOnly(ts *string) string {
	if ts == nil || *ts == "" {
		return "Unknown"
	}

	// Try parsing as Unix timestamp
	if timestamp, err := strconv.ParseInt(*ts, 10, 64); err == nil {
		return time.Unix(timestamp, 0).Format("15:04:05")
	}

	// Try parsing as ISO format
	if t, err := time.Parse(time.RFC3339, *ts); err == nil {
		return t.Local().Format("15:04:05")
	}

	// Return as-is if parsing fails
	return *ts
}

// formatDateOnly formats a timestamp to show only date (YYYY-MM-DD)
func formatDateOnly(ts *string) string {
	if ts == nil || *ts == "" {
		return ""
	}

	// Try parsing as Unix timestamp
	if timestamp, err := strconv.ParseInt(*ts, 10, 64); err == nil {
		return time.Unix(timestamp, 0).Format("2006-01-02")
	}

	// Try parsing as ISO format
	if t, err := time.Parse(time.RFC3339, *ts); err == nil {
		return t.Local().Format("2006-01-02")
	}

	return ""
}

// parseMessageTime parses a timestamp to time.Time
func parseMessageTime(ts *string) *time.Time {
	if ts == nil || *ts == "" {
		return nil
	}

	// Try parsing as Unix timestamp
	if timestamp, err := strconv.ParseInt(*ts, 10, 64); err == nil {
		t := time.Unix(timestamp, 0)
		return &t
	}

	// Try parsing as ISO format
	if t, err := time.Parse(time.RFC3339, *ts); err == nil {
		return &t
	}

	return nil
}

// RenderHTML converts a session and messages to HTML format with collapsible panels and timeline
func RenderHTML(session *db.Session, messages []db.ParsedMessage) string {
	var result strings.Builder

	// Generate title
	title := "Session " + session.ID
	if session.Title != nil && *session.Title != "" {
		title = *session.Title
	}

	// Start HTML document
	result.WriteString(generateHTMLHeader(title))

	// Add session metadata
	result.WriteString(generateSessionInfo(session))

	// Add conversation container
	result.WriteString("<div class=\"conversation\">\n")

	var lastDate string
	for i, msg := range messages {
		// Check if we need a date separator
		msgTime := parseMessageTime(msg.CreatedAt)
		if msgTime != nil {
			currentDate := msgTime.Format("2006-01-02")
			if currentDate != lastDate {
				if i > 0 { // Don't add separator before first message
					result.WriteString(generateDateSeparator(currentDate))
				}
				lastDate = currentDate
			}
		}
		
		result.WriteString(generateMessage(msg, i))
	}

	result.WriteString("</div>\n")

	// Close HTML document
	result.WriteString(generateHTMLFooter())

	return result.String()
}

// generateHTMLHeader creates the HTML header with embedded CSS and JavaScript
func generateHTMLHeader(title string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>%s</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
            line-height: 1.6;
            color: #333;
            background-color: #f5f5f5;
        }

        .container {
            max-width: 100rem;
            margin: 0 auto;
            padding: 20px;
        }

        .header {
            background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%);
            color: white;
            padding: 40px 20px;
            text-align: center;
            border-radius: 10px;
            margin-bottom: 30px;
            box-shadow: 0 4px 6px rgba(0,0,0,0.1);
        }

        .header h1 {
            font-size: 2.5em;
            margin-bottom: 10px;
            font-weight: 300;
        }

        .session-info {
            background: white;
            border-radius: 10px;
            padding: 20px;
            margin-bottom: 30px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }

        .session-info h2 {
            color: #667eea;
            margin-bottom: 15px;
            border-bottom: 2px solid #eee;
            padding-bottom: 10px;
        }

        .info-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 15px;
        }

        .info-item {
            padding: 10px;
            background: #f8f9fa;
            border-radius: 5px;
            border-left: 4px solid #667eea;
        }

        .info-label {
            font-weight: bold;
            color: #666;
            font-size: 0.9em;
            text-transform: uppercase;
            letter-spacing: 0.5px;
        }

        .info-value {
            margin-top: 5px;
            color: #333;
        }


        .conversation {
            background: white;
            border-radius: 10px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
            overflow: hidden;
        }

        .message {
            display: grid;
            grid-template-columns: 280px 1fr;
            min-height: 50px;
            border-bottom: 1px solid #f0f0f0;
        }

        .message:last-child {
            border-bottom: none;
        }

        .message-sidebar {
            padding: 12px 15px;
            border-right: 1px solid #f0f0f0;
            display: flex;
            flex-direction: row;
            align-items: center;
            gap: 12px;
        }

        .message-sidebar.user {
            background: linear-gradient(135deg, #e3f2fd 0%%, #bbdefb 100%%);
        }

        .message-sidebar.assistant {
            background: linear-gradient(135deg, #f3e5f5 0%%, #e1bee7 100%%);
        }

        .message-sidebar.system {
            background: linear-gradient(135deg, #e8f5e8 0%%, #c8e6c9 100%%);
        }

        .role-badge {
            font-weight: bold;
            font-size: 1.2em;
            color: #333;
            min-width: 30px;
            display: flex;
            align-items: center;
            justify-content: center;
        }

        .message-time {
            font-size: 0.8em;
            color: #666;
            line-height: 1.2;
            white-space: nowrap;
        }

        .message-model {
            font-size: 0.75em;
            color: #888;
            font-style: italic;
            white-space: nowrap;
        }

        .message-info {
            display: flex;
            align-items: center;
            gap: 8px;
            flex: 1;
        }

        .message-content {
            padding: 15px 20px;
            display: flex;
            flex-direction: column;
            gap: 10px;
        }

        .message-part {
            padding: 12px;
            background: #f8f9fa;
            border-radius: 6px;
            border-left: 3px solid #667eea;
            white-space: pre-wrap;
            word-wrap: break-word;
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            font-size: 0.9em;
            line-height: 1.5;
        }

        .message-part:only-child {
            background: transparent;
            border: none;
            padding: 0;
        }

        .message-part.tool {
            font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
            background: #f1f1f1;
            border-left-color: #666;
            whitespace: pre;
        }

        .date-separator {
            background: #f8f9fa;
            border-top: 1px solid #e9ecef;
            border-bottom: 1px solid #e9ecef;
            padding: 8px 20px;
            text-align: center;
            font-size: 0.85em;
            font-weight: 600;
            color: #6c757d;
            margin: 10px 0;
        }

        .message-time a {
            color: inherit;
            text-decoration: none;
            opacity: 0.8;
            transition: opacity 0.2s ease;
        }

        .message-time a:hover {
            opacity: 1;
        }

        .back-to-top {
            position: fixed;
            bottom: 30px;
            right: 30px;
            background: #667eea;
            color: white;
            border: none;
            border-radius: 50px;
            padding: 15px;
            cursor: pointer;
            box-shadow: 0 4px 12px rgba(0,0,0,0.3);
            transition: all 0.3s ease;
            font-weight: bold;
            z-index: 1000;
            width: 50px;
            height: 50px;
            display: flex;
            align-items: center;
            justify-content: center;
        }

        .back-to-top:hover {
            background: #5a67d8;
            transform: translateY(-2px);
            box-shadow: 0 6px 16px rgba(0,0,0,0.4);
        }

        @media (max-width: 768px) {
            .container {
                padding: 10px;
            }
            
            .header h1 {
                font-size: 2em;
            }
            
            .message {
                grid-template-columns: 1fr;
            }
            
            .message-sidebar {
                border-right: none;
                border-bottom: 1px solid #f0f0f0;
                flex-wrap: wrap;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>%s</h1>
        </div>
`, html.EscapeString(title), html.EscapeString(title))
}

// generateSessionInfo creates the session information section
func generateSessionInfo(session *db.Session) string {
	var result strings.Builder

	result.WriteString("<div class=\"session-info\">\n")
	result.WriteString("<h2>Session Information</h2>\n")
	result.WriteString("<div class=\"info-grid\">\n")

	// Session ID
	result.WriteString(fmt.Sprintf(`
        <div class="info-item">
            <div class="info-label">Session ID</div>
            <div class="info-value">%s</div>
        </div>
    `, html.EscapeString(session.ID)))

	// Created At
	if session.CreatedAt != nil {
		formattedTime := FormatTimestamp(session.CreatedAt)
		result.WriteString(fmt.Sprintf(`
        <div class="info-item">
            <div class="info-label">Created</div>
            <div class="info-value">%s</div>
        </div>
    `, html.EscapeString(formattedTime)))
	}

	// Message Count
	if session.MessageCount != nil {
		result.WriteString(fmt.Sprintf(`
        <div class="info-item">
            <div class="info-label">Messages</div>
            <div class="info-value">%d</div>
        </div>
    `, *session.MessageCount))
	}

	// Add metadata if present
	if session.Metadata != nil && *session.Metadata != "" {
		var metadata map[string]interface{}
		if err := json.Unmarshal([]byte(*session.Metadata), &metadata); err == nil {
			for k, v := range metadata {
				result.WriteString(fmt.Sprintf(`
        <div class="info-item">
            <div class="info-label">%s</div>
            <div class="info-value">%v</div>
        </div>
    `, html.EscapeString(k), html.EscapeString(fmt.Sprintf("%v", v))))
			}
		}
	}

	result.WriteString("</div>\n</div>\n")
	return result.String()
}

// generateDateSeparator creates a date separator line
func generateDateSeparator(date string) string {
	// Parse date to format it nicely
	if t, err := time.Parse("2006-01-02", date); err == nil {
		formattedDate := t.Format("Monday, January 2, 2006")
		return fmt.Sprintf(`
    <div class="date-separator">
        %s
    </div>
`, formattedDate)
	}
	
	return fmt.Sprintf(`
    <div class="date-separator">
        %s
    </div>
`, date)
}

// generateMessage creates a compact message layout
func generateMessage(msg db.ParsedMessage, index int) string {
	var result strings.Builder

	// Message metadata - use time only format
	timeOnly := "Unknown"
	if msg.CreatedAt != nil {
		timeOnly = formatTimeOnly(msg.CreatedAt)
	}

	var modelInfo []string
	if msg.Model != nil && *msg.Model != "" {
		modelInfo = append(modelInfo, *msg.Model)
	}
	if msg.Provider != nil && *msg.Provider != "" {
		modelInfo = append(modelInfo, *msg.Provider)
	}

	// Create anchor name
	anchorName := fmt.Sprintf("msg-%d", index+1)

	// Generate message
	result.WriteString(fmt.Sprintf(`
    <div class="message" id="%s">
        <div class="message-sidebar %s">
            <div class="role-badge" title="%s">%s</div>
            <div class="message-info">
                <div class="message-time"><a href="#%s">%s</a></div>
`, anchorName, html.EscapeString(msg.Role),
		html.EscapeString(strings.Title(msg.Role)), getRoleEmoji(msg.Role), anchorName, html.EscapeString(timeOnly)))

	// Add model info if available
	if len(modelInfo) > 0 {
		result.WriteString(fmt.Sprintf(`
                <div class="message-model">%s</div>
`, html.EscapeString(strings.Join(modelInfo, "/"))))
	}

	// Close message info and sidebar
	result.WriteString(`
            </div>
        </div>
        <div class="message-content">
`)

	// Add message parts
	for _, part := range msg.Parts {
		// Check if this is a tool message (starts with emoji indicators)
		isToolMessage := strings.HasPrefix(part, "🔧") || strings.HasPrefix(part, "📋")
		cssClass := "message-part"
		if isToolMessage {
			cssClass += " tool"
		}

		result.WriteString(fmt.Sprintf(`
            <div class="%s">%s</div>
`, cssClass, html.EscapeString(part)))
	}

	result.WriteString(`
        </div>
    </div>
`)

	return result.String()
}

// generateHTMLFooter creates the HTML footer with JavaScript
func generateHTMLFooter() string {
	return `
    </div>
    
    <button class="back-to-top" onclick="scrollToTop()" title="Back to top">
        ↑
    </button>

    <script>
        function scrollToTop() {
            window.scrollTo({ top: 0, behavior: 'smooth' });
        }

        // Add smooth scrolling for anchor links
        document.addEventListener('DOMContentLoaded', function() {
            // Handle anchor clicks for smooth scrolling
            document.querySelectorAll('a[href^="#"]').forEach(anchor => {
                anchor.addEventListener('click', function (e) {
                    e.preventDefault();
                    const target = document.querySelector(this.getAttribute('href'));
                    if (target) {
                        target.scrollIntoView({ 
                            behavior: 'smooth',
                            block: 'center'
                        });
                        
                        // Highlight the target message briefly
                        target.style.boxShadow = '0 0 20px rgba(102, 126, 234, 0.5)';
                        setTimeout(() => {
                            target.style.boxShadow = '';
                        }, 2000);
                    }
                });
            });
        });
    </script>
</body>
</html>`
}

// GenerateHTMLFilename generates a filename for the HTML export
func GenerateHTMLFilename(session *db.Session) string {
	// Generate base name from title or session ID
	base := slugify("session-" + session.ID[:8])
	if session.Title != nil && *session.Title != "" {
		base = slugify(*session.Title)
	}

	// Generate timestamp prefix
	prefix := time.Now().Format("2006-01-02_15-04")
	if session.CreatedAt != nil {
		if timestamp, err := time.Parse("1", *session.CreatedAt); err == nil {
			prefix = timestamp.Format("2006-01-02_15-04")
		}
	}

	return fmt.Sprintf("%s_%s.html", prefix, base)
}
