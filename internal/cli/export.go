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
	"crush-session-explorer/internal/providers"

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
	var providerName string

	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export session to markdown or HTML",
		Long:  "Export a session from various AI code tools (Crush, Claude Code, etc.) to Markdown or HTML format",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Check if format was explicitly provided
			formatExplicit := cmd.Flags().Changed("format")

			var selectedProvider providers.Provider
			var allSessions []db.Session
			var providerMap map[string]providers.Provider // Maps session ID to provider

			// If provider is specified, use only that provider
			if providerName != "" {
				provider := providers.GetProvider(providerName)
				if provider == nil {
					return fmt.Errorf("unknown provider: %s", providerName)
				}

				// For Crush provider, allow custom db path
				if providerName == "crush" && dbPath != "" {
					crushProvider := provider.(*providers.CrushProvider)
					crushProvider.SetDBPath(dbPath)
				}

				found, err := provider.Discover()
				if err != nil || !found {
					return fmt.Errorf("provider '%s' data not found", providerName)
				}

				selectedProvider = provider
				sessions, err := provider.ListSessions(50)
				if err != nil {
					return fmt.Errorf("failed to list sessions from %s: %w", providerName, err)
				}
				allSessions = sessions
			} else {
				// Auto-discover all available providers
				availableProviders := providers.DiscoverAllProviders()

				// If dbPath is specified, also try Crush provider with that path
				if dbPath != "" && dbPath != ".crush/crush.db" {
					crushProvider := providers.NewCrushProviderWithPath(dbPath)
					if found, _ := crushProvider.Discover(); found {
						// Check if not already in list
						hasCustomCrush := false
						for _, p := range availableProviders {
							if cp, ok := p.(*providers.CrushProvider); ok && cp == crushProvider {
								hasCustomCrush = true
								break
							}
						}
						if !hasCustomCrush {
							availableProviders = append(availableProviders, crushProvider)
						}
					}
				}

				if len(availableProviders) == 0 {
					return fmt.Errorf("no AI code tool sessions found. Checked: Crush (.crush/crush.db), Claude Code")
				}

				// Collect sessions from all providers
				providerMap = make(map[string]providers.Provider)
				for _, provider := range availableProviders {
					sessions, err := provider.ListSessions(50)
					if err != nil {
						fmt.Fprintf(os.Stderr, "Warning: failed to list sessions from %s: %v\n", provider.Name(), err)
						continue
					}

					// Add provider name to session metadata for display
					for i := range sessions {
						// Store which provider owns this session
						providerMap[sessions[i].ID] = provider
					}

					allSessions = append(allSessions, sessions...)
				}

				if len(allSessions) == 0 {
					return fmt.Errorf("no sessions found in any available providers")
				}
			}

			// If no session ID provided, show interactive selection
			if sessionID == "" {
				if len(allSessions) == 0 {
					return fmt.Errorf("no sessions found")
				}

				// Display sessions with provider info
				fmt.Println("Available sessions:")
				for i, s := range allSessions {
					title := ""
					if s.Title != nil {
						title = *s.Title
					}
					messageCount := 0
					if s.MessageCount != nil {
						messageCount = *s.MessageCount
					}

					// Get provider name from metadata or map
					provider := ""
					if s.Metadata != nil && *s.Metadata != "" {
						provider = *s.Metadata
					} else if providerMap != nil {
						if p := providerMap[s.ID]; p != nil {
							provider = p.Name()
						}
					}

					providerTag := ""
					if provider != "" {
						providerTag = fmt.Sprintf(" [%s]", provider)
					}

					fmt.Printf("%2d. %s — %s — %s — %d msg%s\n",
						i+1, s.ID, formatTimestamp(s.CreatedAt), title, messageCount, providerTag)
				}

				// Get user selection
				fmt.Print("Select session number: ")
				reader := bufio.NewReader(os.Stdin)
				input, err := reader.ReadString('\n')
				if err != nil {
					return fmt.Errorf("failed to read input: %w", err)
				}

				selection, err := strconv.Atoi(strings.TrimSpace(input))
				if err != nil || selection < 1 || selection > len(allSessions) {
					return fmt.Errorf("invalid selection")
				}

				sessionID = allSessions[selection-1].ID

				// Set the provider for this session
				if providerMap != nil {
					selectedProvider = providerMap[sessionID]
				}
			}

			// If we still don't have a provider, determine it from the session ID
			if selectedProvider == nil {
				if providerMap != nil {
					selectedProvider = providerMap[sessionID]
				}
				if selectedProvider == nil {
					// Default to Crush provider
					selectedProvider = providers.NewCrushProviderWithPath(dbPath)
				}
			}

			// Fetch session using the selected provider
			session, err := selectedProvider.FetchSession(sessionID)
			if err != nil {
				return fmt.Errorf("failed to fetch session: %w", err)
			}

			// Fetch messages using the selected provider
			messages, err := selectedProvider.ListMessages(session.ID)
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

	cmd.Flags().StringVar(&dbPath, "db", ".crush/crush.db", "Path to sqlite database (for Crush provider)")
	cmd.Flags().StringVar(&sessionID, "session", "", "Session ID to export")
	cmd.Flags().StringVar(&outputPath, "out", "", "Output file path")
	cmd.Flags().StringVar(&format, "format", "markdown", "Output format: markdown, html, md (interactive selection if not specified)")
	cmd.Flags().StringVar(&providerName, "provider", "", "AI code tool provider: crush, claude-code (auto-detect if not specified)")

	return cmd
}
