package markdown

import (
	"encoding/json"
	"fmt"
	"html"
	"strings"
	"time"

	"crush-session-explorer/internal/db"
)

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

	// Add timeline navigation
	result.WriteString(generateTimeline(messages))

	// Add messages container
	result.WriteString("<div class=\"messages-container\">\n")

	for i, msg := range messages {
		result.WriteString(generateMessagePanel(msg, i))
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
            max-width: 1200px;
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

        .timeline {
            background: white;
            border-radius: 10px;
            padding: 20px;
            margin-bottom: 30px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }

        .timeline h2 {
            color: #667eea;
            margin-bottom: 15px;
            border-bottom: 2px solid #eee;
            padding-bottom: 10px;
        }

        .timeline-nav {
            display: flex;
            flex-wrap: wrap;
            gap: 10px;
        }

        .timeline-item {
            padding: 8px 12px;
            background: #f8f9fa;
            border-radius: 20px;
            border: 1px solid #ddd;
            cursor: pointer;
            transition: all 0.3s ease;
            font-size: 0.9em;
        }

        .timeline-item:hover {
            background: #667eea;
            color: white;
            transform: translateY(-1px);
        }

        .timeline-item.user {
            background: #e3f2fd;
            border-color: #2196f3;
        }

        .timeline-item.assistant {
            background: #f3e5f5;
            border-color: #9c27b0;
        }

        .messages-container {
            display: flex;
            flex-direction: column;
            gap: 20px;
        }

        .message-panel {
            background: white;
            border-radius: 10px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
            overflow: hidden;
            transition: all 0.3s ease;
        }

        .message-panel:hover {
            box-shadow: 0 4px 8px rgba(0,0,0,0.15);
        }

        .message-header {
            padding: 15px 20px;
            cursor: pointer;
            display: flex;
            justify-content: space-between;
            align-items: center;
            transition: background-color 0.3s ease;
        }

        .message-header:hover {
            background-color: #f8f9fa;
        }

        .message-header.user {
            background: linear-gradient(135deg, #42a5f5 0%%, #1e88e5 100%%);
            color: white;
        }

        .message-header.assistant {
            background: linear-gradient(135deg, #ab47bc 0%%, #8e24aa 100%%);
            color: white;
        }

        .message-header.system {
            background: linear-gradient(135deg, #66bb6a 0%%, #43a047 100%%);
            color: white;
        }

        .message-info {
            display: flex;
            align-items: center;
            gap: 15px;
        }

        .role-badge {
            font-weight: bold;
            font-size: 1.1em;
            text-transform: capitalize;
        }

        .message-meta {
            font-size: 0.9em;
            opacity: 0.9;
        }

        .toggle-icon {
            font-size: 1.2em;
            transition: transform 0.3s ease;
        }

        .toggle-icon.expanded {
            transform: rotate(180deg);
        }

        .message-content {
            max-height: 0;
            overflow: hidden;
            transition: max-height 0.3s ease;
        }

        .message-content.expanded {
            max-height: 10000px;
        }

        .content-inner {
            padding: 20px;
            border-top: 1px solid #eee;
        }

        .message-part {
            margin-bottom: 15px;
            padding: 15px;
            background: #f8f9fa;
            border-radius: 8px;
            white-space: pre-wrap;
            word-wrap: break-word;
            font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
            font-size: 0.95em;
            line-height: 1.5;
        }

        .expand-all-btn {
            position: fixed;
            bottom: 30px;
            right: 30px;
            background: #667eea;
            color: white;
            border: none;
            border-radius: 50px;
            padding: 15px 20px;
            cursor: pointer;
            box-shadow: 0 4px 12px rgba(0,0,0,0.3);
            transition: all 0.3s ease;
            font-weight: bold;
            z-index: 1000;
        }

        .expand-all-btn:hover {
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
            
            .timeline-nav {
                justify-content: center;
            }
            
            .message-info {
                flex-direction: column;
                gap: 5px;
                align-items: flex-start;
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

// generateTimeline creates the timeline navigation
func generateTimeline(messages []db.ParsedMessage) string {
	var result strings.Builder

	result.WriteString("<div class=\"timeline\">\n")
	result.WriteString("<h2>Timeline</h2>\n")
	result.WriteString("<div class=\"timeline-nav\">\n")

	for i, msg := range messages {
		timestamp := "Unknown"
		if msg.CreatedAt != nil {
			timestamp = FormatTimestamp(msg.CreatedAt)
		}

		result.WriteString(fmt.Sprintf(`
        <div class="timeline-item %s" onclick="scrollToMessage(%d)" title="%s">
            %s: %s
        </div>
    `, html.EscapeString(msg.Role), i, html.EscapeString(timestamp), 
		html.EscapeString(strings.Title(msg.Role)), html.EscapeString(timestamp)))
	}

	result.WriteString("</div>\n</div>\n")
	return result.String()
}

// generateMessagePanel creates a collapsible message panel
func generateMessagePanel(msg db.ParsedMessage, index int) string {
	var result strings.Builder

	// Message metadata
	timestamp := "Unknown time"
	if msg.CreatedAt != nil {
		timestamp = FormatTimestamp(msg.CreatedAt)
	}

	var modelInfo []string
	if msg.Model != nil && *msg.Model != "" {
		modelInfo = append(modelInfo, *msg.Model)
	}
	if msg.Provider != nil && *msg.Provider != "" {
		modelInfo = append(modelInfo, *msg.Provider)
	}

	metaText := timestamp
	if len(modelInfo) > 0 {
		metaText += " • " + strings.Join(modelInfo, "/")
	}

	// Generate panel
	result.WriteString(fmt.Sprintf(`
    <div class="message-panel" id="message-%d">
        <div class="message-header %s" onclick="toggleMessage(%d)">
            <div class="message-info">
                <div class="role-badge">%s</div>
                <div class="message-meta">%s</div>
            </div>
            <div class="toggle-icon" id="toggle-%d">▼</div>
        </div>
        <div class="message-content" id="content-%d">
            <div class="content-inner">
`, index, html.EscapeString(msg.Role), index, 
	html.EscapeString(strings.Title(msg.Role)), html.EscapeString(metaText), index, index))

	// Add message parts
	for i, part := range msg.Parts {
		result.WriteString(fmt.Sprintf(`
                <div class="message-part">
                    <strong>Part %d:</strong><br>
                    %s
                </div>
`, i+1, html.EscapeString(part)))
	}

	result.WriteString(`
            </div>
        </div>
    </div>
`)

	return result.String()
}

// generateHTMLFooter creates the HTML footer with JavaScript
func generateHTMLFooter() string {
	return `
    </div>
    
    <button class="expand-all-btn" onclick="toggleAllMessages()" id="expandAllBtn">
        Expand All
    </button>

    <script>
        let allExpanded = false;

        function toggleMessage(index) {
            const content = document.getElementById('content-' + index);
            const toggle = document.getElementById('toggle-' + index);
            
            if (content.classList.contains('expanded')) {
                content.classList.remove('expanded');
                toggle.classList.remove('expanded');
            } else {
                content.classList.add('expanded');
                toggle.classList.add('expanded');
            }
        }

        function toggleAllMessages() {
            const contents = document.querySelectorAll('.message-content');
            const toggles = document.querySelectorAll('.toggle-icon');
            const btn = document.getElementById('expandAllBtn');
            
            if (allExpanded) {
                contents.forEach(content => content.classList.remove('expanded'));
                toggles.forEach(toggle => toggle.classList.remove('expanded'));
                btn.textContent = 'Expand All';
                allExpanded = false;
            } else {
                contents.forEach(content => content.classList.add('expanded'));
                toggles.forEach(toggle => toggle.classList.add('expanded'));
                btn.textContent = 'Collapse All';
                allExpanded = true;
            }
        }

        function scrollToMessage(index) {
            const element = document.getElementById('message-' + index);
            element.scrollIntoView({ behavior: 'smooth', block: 'center' });
            
            // Highlight the message briefly
            element.style.boxShadow = '0 0 20px rgba(102, 126, 234, 0.5)';
            setTimeout(() => {
                element.style.boxShadow = '0 2px 4px rgba(0,0,0,0.1)';
            }, 2000);
        }

        // Auto-expand first message on load
        document.addEventListener('DOMContentLoaded', function() {
            if (document.getElementById('content-0')) {
                toggleMessage(0);
            }
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