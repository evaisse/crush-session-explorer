package cli

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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

// openInBrowser opens a file in the default browser
func openInBrowser(filePath string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", filePath)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", filePath)
	case "linux":
		cmd = exec.Command("xdg-open", filePath)
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	return cmd.Start()
}

// promptOpenInBrowser asks user if they want to open the HTML file in browser
func promptOpenInBrowser(filePath string) {
	if !strings.HasSuffix(strings.ToLower(filePath), ".html") {
		return // Only prompt for HTML files
	}

	fmt.Printf("\n🌐 Open %s in browser? [y/N]: ", filepath.Base(filePath))
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return
	}

	choice := strings.TrimSpace(strings.ToLower(input))
	if choice == "y" || choice == "yes" {
		// Convert to absolute path for better browser compatibility
		absPath, err := filepath.Abs(filePath)
		if err != nil {
			fmt.Printf("❌ Failed to get absolute path: %v\n", err)
			return
		}

		if err := openInBrowser(absPath); err != nil {
			fmt.Printf("❌ Failed to open in browser: %v\n", err)
			fmt.Printf("💡 You can manually open: file://%s\n", absPath)
		} else {
			fmt.Printf("✅ Opening in browser...\n")
		}
	}
}

// ExportCmd creates the export command
func ExportCmd() *cobra.Command {
	var dbPath string
	var sessionID string
	var outputPath string
	var format string

	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export session to markdown or HTML",
		Long:  "Export a Crush session from SQLite database to Markdown or HTML format",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Check if format was explicitly provided
			formatExplicit := cmd.Flags().Changed("format")
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

			// Interactive format selection if not explicitly provided
			if !formatExplicit {
				fmt.Println("Choose export format:")
				fmt.Println("1. Markdown (.md)")
				fmt.Println("2. HTML with interactive panels (.html)")

				fmt.Print("Select format [1-2]: ")
				reader := bufio.NewReader(os.Stdin)
				input, err := reader.ReadString('\n')
				if err != nil {
					return fmt.Errorf("failed to read input: %w", err)
				}

				choice := strings.TrimSpace(input)
				switch choice {
				case "1", "":
					format = "markdown"
				case "2":
					format = "html"
				default:
					return fmt.Errorf("invalid choice: %s (choose 1 or 2)", choice)
				}
			} else {
				// Validate format when explicitly provided
				if format != "markdown" && format != "html" && format != "md" {
					return fmt.Errorf("invalid format: %s (supported: markdown, html, md)", format)
				}

				// Normalize format
				if format == "md" {
					format = "markdown"
				}
			}

			// Generate output path if not provided
			if outputPath == "" {
				var filename string
				if format == "html" {
					filename = markdown.GenerateHTMLFilename(session)
				} else {
					filename = markdown.GenerateFilename(session)
				}
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

			// Render content based on format
			var content string
			if format == "html" {
				content = markdown.RenderHTML(session, messages)
			} else {
				content = markdown.RenderMarkdown(session, messages)
			}

			// Ensure output directory exists
			if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
				return fmt.Errorf("failed to create output directory: %w", err)
			}

			// Write file
			if err := os.WriteFile(outputPath, []byte(content), 0644); err != nil {
				return fmt.Errorf("failed to write output file: %w", err)
			}

			fmt.Println(outputPath)

			// Prompt to open in browser if HTML format
			promptOpenInBrowser(outputPath)

			return nil
		},
	}

	cmd.Flags().StringVar(&dbPath, "db", ".crush/crush.db", "Path to sqlite database")
	cmd.Flags().StringVar(&sessionID, "session", "", "Session ID to export")
	cmd.Flags().StringVar(&outputPath, "out", "", "Output file path")
	cmd.Flags().StringVar(&format, "format", "markdown", "Output format: markdown, html, md (interactive selection if not specified)")

	return cmd
}
