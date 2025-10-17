# Crush Session Explorer

A fast, lightweight CLI tool written in Go for exporting chat sessions from multiple AI code tools (Crush, Claude Code, etc.) to Markdown or HTML format with YAML frontmatter.

## Overview

This tool allows you to extract and export individual chat sessions from various AI coding assistants' databases, converting them into well-formatted Markdown or HTML files with structured metadata. Perfect for archiving, documentation, or further processing of conversation data.

## Features

- 🔌 **Multi-Provider Support**: Auto-discover and export sessions from multiple AI code tools
  - Crush (`.crush/crush.db`)
  - Claude Code/Desktop
  - Extensible architecture for adding more providers
- 📊 **SQLite Integration**: Direct access to databases using Go's SQLite driver
- 📝 **Markdown Export**: Clean conversion to Markdown with YAML frontmatter
- 🌐 **HTML Export**: Interactive HTML with collapsible panels and syntax highlighting
- 🔍 **Interactive Session Selection**: Browse and select sessions from all available providers
- 📅 **Timestamp Formatting**: Automatic timestamp conversion to readable formats
- 🏷️ **Metadata Preservation**: Session metadata and message details preserved
- 🎯 **Type Safety**: Full compile-time type checking with Go
- ⚡ **Fast Performance**: Compiled binary with no runtime dependencies
- 🚀 **Cross-Platform**: Build for Linux, macOS, and Windows

## Requirements

- **Go 1.19+**
- **CGO enabled** (for SQLite driver)

## Installation

### Quick Setup

```bash
git clone <repository-url>
cd crush-session-explorer
go mod download
make build
```

### Install Globally

```bash
go install ./cmd/crush-md
```

### Cross-Platform Builds

```bash
make build-all  # Creates binaries for all platforms in bin/
```

## Usage

### Auto-Discovery Mode (Recommended)

The tool automatically discovers all available AI code tool sessions on your system:

```bash
./bin/crush-md export
```

This will scan for:
- Crush sessions in `.crush/crush.db`
- Claude Code sessions in the default location (OS-dependent)
- Any custom database paths you specify

Example output:
```
Available sessions:
 1. abc123 — 2024-01-15 14:30 — Project Discussion — 12 msg [crush]
 2. def456 — 2024-01-15 13:45 — Code Review — 8 msg [claude-code]
 3. ghi789 — 2024-01-15 12:20 — Planning Meeting — 15 msg [crush]
Select session number: 
```

### Export a Specific Session

```bash
./bin/crush-md export --session <session-id> --out output.md
```

### Specify Provider Explicitly

```bash
# Export from Crush only
./bin/crush-md export --provider crush --db ./.crush/crush.db

# Export from Claude Code only
./bin/crush-md export --provider claude-code --claude-db ~/Library/Application\ Support/Claude/state.db
```

### Custom Database Paths

```bash
# Custom Crush database
./bin/crush-md export --db /path/to/custom/crush.db

# Custom Claude database
./bin/crush-md export --claude-db /path/to/custom/claude.db

# Both providers with custom paths
./bin/crush-md export --db /path/to/crush.db --claude-db /path/to/claude.db
```

### Command Line Options

| Option | Description | Default |
|--------|-------------|---------|
| `--db` | Path to Crush SQLite database | `.crush/crush.db` |
| `--claude-db` | Path to Claude SQLite database | Auto-detected by OS |
| `--provider` | Specific provider to use | Auto-detect all |
| `--session` | Specific session ID to export | Interactive selection |
| `--out` | Output file path | Auto-generated |
| `--format` | Output format (markdown, html) | Interactive selection |

### Supported Providers

| Provider | Description | Default Database Path |
|----------|-------------|----------------------|
| `crush` | Crush AI code tool | `.crush/crush.db` |
| `claude-code` | Claude Desktop/Code | macOS: `~/Library/Application Support/Claude/state.db`<br>Linux: `~/.config/Claude/state.db`<br>Windows: `%APPDATA%/Claude/state.db` |

More providers can be easily added by implementing the `Provider` interface.

### Output Format

Exported Markdown files include:

- **YAML Frontmatter**: Session metadata (title, ID, timestamps, message count)
- **Message History**: Chronological conversation with role indicators
- **Timestamps**: Human-readable message timestamps
- **Model Information**: AI model and provider details where available

Example output structure:

```markdown
---
title: "Project Discussion"
session_id: abc123def456
created_at: 2024-01-15T14:30:00
message_count: 12
metadata:
  model: "claude-3"
  provider: "anthropic"
---

## user — 2024-01-15 14:30

<div>
Can you help me understand this code structure?
</div>

## assistant — 2024-01-15 14:31 (claude-3/anthropic)

<div>
I'd be happy to help explain the code structure...
</div>
```

## Project Structure

```
crush-session-explorer/
├── cmd/                          # CLI application entry point
│   └── crush-md/
│       └── main.go               # Main CLI application
├── internal/                     # Internal Go packages
│   ├── cli/
│   │   └── export.go             # Export command implementation
│   ├── db/
│   │   ├── connection.go         # Database connection
│   │   ├── models.go             # Data models
│   │   └── queries.go            # Database queries
│   ├── markdown/
│   │   ├── renderer.go           # Markdown rendering
│   │   ├── html_renderer.go     # HTML rendering
│   │   └── utils.go              # Utility functions
│   └── providers/                # AI code tool providers
│       ├── provider.go           # Provider interface
│       ├── crush.go              # Crush provider implementation
│       └── claude.go             # Claude Code provider implementation
├── bin/                          # Build output (created by make build)
├── go.mod                        # Go module definition
├── go.sum                        # Go dependencies
├── README.md                     # This file
└── Makefile                      # Build and development commands
```

## Development

### Building

```bash
make build          # Build for current platform
make build-all      # Build for all platforms
make dev            # Format, vet, test, and build
```

### Testing

```bash
make test           # Run tests
make test-coverage  # Run tests with coverage
```

### Code Quality

```bash
make fmt            # Format code
make vet            # Vet code
make lint           # Lint code (requires golangci-lint)
make check          # Run format, vet, and test
```

### Development Setup

```bash
make dev-setup      # Install development tools
```

### Making Changes

The codebase follows these principles:

- **Security First**: Uses parameterized queries to prevent SQL injection
- **Type Safety**: Full compile-time type checking with Go
- **Clean Code**: Formatted with gofmt, following Go best practices
- **Testable**: Modular design with comprehensive test coverage
- **Performance**: Efficient memory usage and fast execution

### Adding New Providers

The tool is designed to be easily extensible. To add support for a new AI code tool:

1. Create a new file in `internal/providers/` (e.g., `cursor.go`)
2. Implement the `Provider` interface:
   ```go
   type Provider interface {
       Name() string
       Discover() (bool, error)
       ListSessions(limit int) ([]db.Session, error)
       FetchSession(sessionID string) (*db.Session, error)
       ListMessages(sessionID string) ([]db.ParsedMessage, error)
   }
   ```
3. Add your provider to the `DiscoverAllProviders()` function in `provider.go`
4. Update the `GetProvider()` function to include your provider name

Example providers:
- **Crush**: SQLite-based sessions in `.crush/crush.db`
- **Claude Code**: SQLite-based conversations in Claude Desktop's database
- **Cursor**: (Future) Support for Cursor editor sessions
- **Copilot**: (Future) Support for GitHub Copilot sessions

## Security Notes

- All database queries use parameterized statements to prevent SQL injection
- No user credentials or sensitive data are logged or exposed
- SQLite database is accessed in read-only mode for exports
- Input validation is performed on all user-provided parameters

## Troubleshooting

### Common Issues

**No sessions found:**
```bash
# Check which providers are being detected
./bin/crush-md export --db .crush/crush.db --claude-db ~/Library/Application\ Support/Claude/state.db

# Verify database files exist
ls -la .crush/crush.db
ls -la ~/Library/Application\ Support/Claude/state.db  # macOS

# Verify database has expected tables
sqlite3 .crush/crush.db ".tables"
```

**Claude Code sessions not found:**
The Claude database location varies by operating system:
- **macOS**: `~/Library/Application Support/Claude/state.db`
- **Linux**: `~/.config/Claude/state.db`
- **Windows**: `%APPDATA%/Claude/state.db`

If Claude stores data in a different location, use the `--claude-db` flag:
```bash
./bin/crush-md export --claude-db /path/to/claude/database.db
```

**Database not found:**
```bash
# Check if the database exists
ls -la .crush/crush.db

# Verify database permissions
file .crush/crush.db
```

**Export fails:**
```bash
# Check output directory permissions
mkdir -p $(dirname your-output-file.md)
```

**CGO compilation issues:**
```bash
# Ensure CGO is enabled (required for SQLite)
export CGO_ENABLED=1
go build ./cmd/crush-md
```

**Provider-specific issues:**
```bash
# Test a specific provider
./bin/crush-md export --provider crush --db .crush/crush.db
./bin/crush-md export --provider claude-code --claude-db ~/path/to/claude.db

# Check provider detection
./bin/crush-md export  # Will show which providers were found
```

## Performance

The Go implementation offers significant performance advantages:

- **Startup Time**: Near-instantaneous startup (vs ~500ms for Python)
- **Memory Usage**: ~10MB RAM (vs ~50MB for Python with dependencies)
- **Binary Size**: ~15MB standalone executable
- **Export Speed**: 2-3x faster than Python implementation

## License

This project is provided as-is for educational and archival purposes.

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes following the existing code style
4. Add tests for new functionality
5. Ensure all checks pass (`make check`)
6. Build and test (`make dev`)
7. Commit your changes (`git commit -m 'feat: add amazing feature'`)
8. Push to the branch (`git push origin feature/amazing-feature`)
9. Open a Pull Request