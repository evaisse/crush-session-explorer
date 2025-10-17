package cli

import (
	"fmt"
	"path/filepath"

	"crush-session-explorer/internal/interchange"
	"crush-session-explorer/internal/markdown"

	"github.com/spf13/cobra"
)

// ImportAICSCmd creates the import-aics command
func ImportAICSCmd() *cobra.Command {
	var inputPath string
	var outputDir string
	var format string

	cmd := &cobra.Command{
		Use:   "import-aics",
		Short: "Import sessions from AICS (AI Coding Session) interchange format",
		Long: `Import sessions from the AICS standard interchange format.

The AICS format is a standardized JSON format for AI coding sessions.
This command imports AICS files and converts them to markdown or HTML format.

Use this to:
- Import sessions from other AI coding tools
- Migrate from another tool to your current workflow
- Convert archived sessions to readable format`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if inputPath == "" {
				return fmt.Errorf("input file path is required (use --input)")
			}

			// Import from AICS file
			fmt.Printf("📥 Importing from %s...\n", inputPath)
			archive, err := interchange.ImportFromFile(inputPath)
			if err != nil {
				return fmt.Errorf("failed to import AICS file: %w", err)
			}

			// Validate the archive
			if err := interchange.ValidateArchive(archive); err != nil {
				return fmt.Errorf("invalid AICS file: %w", err)
			}

			fmt.Printf("✅ Successfully imported AICS archive\n")
			fmt.Printf("📊 Format version: %s\n", archive.Version)
			fmt.Printf("🔧 Created by: %s v%s\n", archive.Creator.Name, archive.Creator.Version)
			if archive.Browser != nil {
				fmt.Printf("🌐 Original tool: %s\n", archive.Browser.Name)
			}
			fmt.Printf("📝 Sessions: %d\n", len(archive.Log.Sessions))

			// Convert to database format
			sessions, messagesMap, err := archive.ConvertToDBFormat()
			if err != nil {
				return fmt.Errorf("failed to convert to database format: %w", err)
			}

			// Validate format
			if format != "markdown" && format != "html" && format != "md" {
				return fmt.Errorf("invalid format: %s (supported: markdown, html, md)", format)
			}

			// Normalize format
			if format == "md" {
				format = "markdown"
			}

			// Export each session to the specified format
			fmt.Printf("\n📤 Exporting sessions to %s format...\n", format)

			successCount := 0
			for _, session := range sessions {
				messages := messagesMap[session.ID]

				// Generate filename
				var filename string
				var content string

				if format == "html" {
					filename = markdown.GenerateHTMLFilename(&session)
					content = markdown.RenderHTML(&session, messages)
				} else {
					filename = markdown.GenerateFilename(&session)
					content = markdown.RenderMarkdown(&session, messages)
				}

				// Create output path
				outputPath := filepath.Join(outputDir, filename)

				// Write file
				if err := markdown.WriteFile(outputPath, content); err != nil {
					fmt.Printf("❌ Failed to export session %s: %v\n", session.ID, err)
					continue
				}

				successCount++
				fmt.Printf("  ✓ %s\n", filename)
			}

			fmt.Printf("\n✅ Successfully exported %d/%d sessions to %s\n",
				successCount, len(sessions), outputDir)

			return nil
		},
	}

	cmd.Flags().StringVar(&inputPath, "input", "", "Input AICS file path (required)")
	cmd.Flags().StringVar(&outputDir, "out", "imported-sessions", "Output directory for exported sessions")
	cmd.Flags().StringVar(&format, "format", "markdown", "Output format: markdown, html, md")

	cmd.MarkFlagRequired("input")

	return cmd
}
