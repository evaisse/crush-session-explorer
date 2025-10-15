package cli

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"crush-session-explorer/internal/db"
	"crush-session-explorer/internal/markdown"

	"github.com/spf13/cobra"
)

// formatTimestamp formats a timestamp for display in session list
func formatTimestamp(ts *string) string {
	if ts == nil || *ts == "" {
		return ""
	}
	return markdown.FormatTimestamp(ts)
}

// ExportCmd creates the export command
func ExportCmd() *cobra.Command {
	var dbPath string
	var sessionID string
	var outputPath string

	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export session to markdown",
		Long:  "Export a Crush session from SQLite database to Markdown format with YAML frontmatter",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Connect to database
			database, err := db.Connect(dbPath)
			if err != nil {
				return fmt.Errorf("failed to connect to database: %w", err)
			}
			defer database.Close()

			// If no session ID provided, show interactive selection
			if sessionID == "" {
				sessions, err := db.ListSessions(database, 50)
				if err != nil {
					return fmt.Errorf("failed to list sessions: %w", err)
				}

				if len(sessions) == 0 {
					return fmt.Errorf("no sessions found in database")
				}

				// Display sessions
				fmt.Println("Available sessions:")
				for i, s := range sessions {
					title := ""
					if s.Title != nil {
						title = *s.Title
					}
					messageCount := 0
					if s.MessageCount != nil {
						messageCount = *s.MessageCount
					}
					fmt.Printf("%2d. %s — %s — %s — %d msg\n", 
						i+1, s.ID, formatTimestamp(s.CreatedAt), title, messageCount)
				}

				// Get user selection
				fmt.Print("Select session number: ")
				reader := bufio.NewReader(os.Stdin)
				input, err := reader.ReadString('\n')
				if err != nil {
					return fmt.Errorf("failed to read input: %w", err)
				}

				selection, err := strconv.Atoi(strings.TrimSpace(input))
				if err != nil || selection < 1 || selection > len(sessions) {
					return fmt.Errorf("invalid selection")
				}

				sessionID = sessions[selection-1].ID
			}

			// Fetch session
			session, err := db.FetchSession(database, sessionID)
			if err != nil {
				return fmt.Errorf("failed to fetch session: %w", err)
			}

			// Fetch messages
			messages, err := db.ListMessages(database, session.ID)
			if err != nil {
				return fmt.Errorf("failed to fetch messages: %w", err)
			}

			// Set session content as JSON for compatibility
			if len(messages) > 0 {
				contentBytes, _ := json.Marshal(messages)
				contentStr := string(contentBytes)
				session.Content = &contentStr
			}

			// Generate output path if not provided
			if outputPath == "" {
				filename := markdown.GenerateFilename(session)
				defaultDir := ".crush/sessions"
				defaultPath := filepath.Join(defaultDir, filename)

				fmt.Printf("Output file path [%s]: ", defaultPath)
				reader := bufio.NewReader(os.Stdin)
				input, err := reader.ReadString('\n')
				if err != nil {
					return fmt.Errorf("failed to read input: %w", err)
				}

				input = strings.TrimSpace(input)
				if input == "" {
					outputPath = defaultPath
				} else {
					outputPath = input
				}
			}

			// Render markdown
			markdownContent := markdown.RenderMarkdown(session, messages)

			// Ensure output directory exists
			if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
				return fmt.Errorf("failed to create output directory: %w", err)
			}

			// Write file
			if err := os.WriteFile(outputPath, []byte(markdownContent), 0644); err != nil {
				return fmt.Errorf("failed to write output file: %w", err)
			}

			fmt.Println(outputPath)
			return nil
		},
	}

	cmd.Flags().StringVar(&dbPath, "db", ".crush/crush.db", "Path to sqlite database")
	cmd.Flags().StringVar(&sessionID, "session", "", "Session ID to export")
	cmd.Flags().StringVar(&outputPath, "out", "", "Output markdown file path")

	return cmd
}