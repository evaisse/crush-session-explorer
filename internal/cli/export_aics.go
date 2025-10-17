package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"crush-session-explorer/internal/db"
	"crush-session-explorer/internal/interchange"

	"github.com/spf13/cobra"
)

// ExportAICSCmd creates the export-aics command
func ExportAICSCmd() *cobra.Command {
	var dbPath string
	var outputPath string
	var providerName string
	var limit int

	cmd := &cobra.Command{
		Use:   "export-aics",
		Short: "Export sessions to AICS (AI Coding Session) interchange format",
		Long: `Export sessions to the AICS standard interchange format.

The AICS format is a standardized JSON format for AI coding sessions,
inspired by HAR (HTTP Archive) format. It allows migration and synchronization
between different AI coding tools like Cursor, Claude Code, and others.

Benefits:
- Switch between AI tools while preserving session history
- Share session data in a vendor-neutral format
- Archive conversations for future reference
- Migrate from one tool to another seamlessly`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Connect to database
			database, err := db.Connect(dbPath)
			if err != nil {
				return fmt.Errorf("failed to connect to database: %w", err)
			}
			defer database.Close()

			// Fetch sessions
			sessions, err := db.ListSessions(database, limit)
			if err != nil {
				return fmt.Errorf("failed to list sessions: %w", err)
			}

			if len(sessions) == 0 {
				return fmt.Errorf("no sessions found in database")
			}

			fmt.Printf("Found %d sessions to export\n", len(sessions))

			// Fetch messages for each session
			messagesMap := make(map[string][]db.ParsedMessage)
			for _, session := range sessions {
				messages, err := db.ListMessages(database, session.ID)
				if err != nil {
					return fmt.Errorf("failed to fetch messages for session %s: %w", session.ID, err)
				}
				messagesMap[session.ID] = messages
			}

			// Export to AICS format
			archive, err := interchange.ExportToAICS(sessions, messagesMap, providerName)
			if err != nil {
				return fmt.Errorf("failed to export to AICS: %w", err)
			}

			// Convert to JSON
			jsonData, err := archive.ToJSON()
			if err != nil {
				return fmt.Errorf("failed to convert to JSON: %w", err)
			}

			// Generate output path if not provided
			if outputPath == "" {
				outputPath = "sessions.aics.json"
			}

			// Ensure output directory exists
			if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
				return fmt.Errorf("failed to create output directory: %w", err)
			}

			// Write file
			if err := os.WriteFile(outputPath, jsonData, 0644); err != nil {
				return fmt.Errorf("failed to write output file: %w", err)
			}

			fmt.Printf("✅ Exported %d sessions to %s\n", len(sessions), outputPath)
			fmt.Printf("📊 Format: AICS v%s (AI Coding Session Interchange Format)\n", interchange.FormatVersion)
			fmt.Printf("💡 This file can be imported into other AI coding tools that support AICS\n")

			return nil
		},
	}

	cmd.Flags().StringVar(&dbPath, "db", ".crush/crush.db", "Path to sqlite database")
	cmd.Flags().StringVar(&outputPath, "out", "", "Output AICS file path (default: sessions.aics.json)")
	cmd.Flags().StringVar(&providerName, "provider", "Crush", "Name of the AI provider/tool")
	cmd.Flags().IntVar(&limit, "limit", 50, "Maximum number of sessions to export")

	return cmd
}
