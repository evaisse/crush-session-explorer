# Crush Session Explorer

A fast, lightweight CLI tool written in Go for exporting Crush chat sessions from SQLite databases to Markdown format with YAML frontmatter.

## Overview

This tool allows you to extract and export individual chat sessions from Crush's SQLite database, converting them into well-formatted Markdown files with structured metadata. Perfect for archiving, documentation, or further processing of conversation data.

## Features

- 📊 **SQLite Integration**: Direct access to Crush database using Go's SQLite driver
- 📝 **Markdown Export**: Clean conversion to Markdown with YAML frontmatter
- 🔄 **AICS Format Support**: Export/import sessions using the standardized AI Coding Session interchange format
- 🔀 **Cross-Tool Migration**: Migrate sessions between different AI coding tools (Cursor, Claude Code, etc.)
- 🔍 **Interactive Session Selection**: Browse and select sessions interactively
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

### Export a Specific Session

```bash
./bin/crush-md export --db ./.crush/crush.db --session <session-id> --out output.md
# or if installed globally: crush-md export --db ./.crush/crush.db --session <session-id> --out output.md
```

### Export to AICS Interchange Format

Export sessions to the standardized AICS (AI Coding Session) format for migration to other tools:

```bash
./bin/crush-md export-aics --db ./.crush/crush.db --out sessions.aics.json
```

This creates a vendor-neutral JSON file that can be imported into other AI coding tools.

### Import from AICS Format

Import sessions from other AI coding tools that support AICS format:

```bash
./bin/crush-md import-aics --input sessions.aics.json --format markdown --out ./imported/
```

### Interactive Session Selection

```bash
./bin/crush-md export --db ./.crush/crush.db
```

This will display a list of recent sessions for you to choose from:

```
 1. abc123 — 2024-01-15 14:30 — Project Discussion — 12 msg
 2. def456 — 2024-01-15 13:45 — Code Review — 8 msg
 3. ghi789 — 2024-01-15 12:20 — Planning Meeting — 15 msg
Select session number: 
```

### Command Line Options

#### export command

| Option | Description | Default |
|--------|-------------|---------|
| `--db` | Path to the SQLite database | `.crush/crush.db` |
| `--session` | Specific session ID to export | Interactive selection |
| `--out` | Output Markdown file path | Auto-generated based on session |
| `--format` | Output format (markdown, html, md) | Interactive selection |

#### export-aics command

| Option | Description | Default |
|--------|-------------|---------|
| `--db` | Path to the SQLite database | `.crush/crush.db` |
| `--out` | Output AICS file path | `sessions.aics.json` |
| `--provider` | Name of the AI provider/tool | `Crush` |
| `--limit` | Maximum number of sessions to export | `50` |

#### import-aics command

| Option | Description | Default |
|--------|-------------|---------|
| `--input` | Input AICS file path | Required |
| `--out` | Output directory for exported sessions | `imported-sessions` |
| `--format` | Output format (markdown, html, md) | `markdown` |

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

## AICS Format (AI Coding Session Interchange Format)

The tool now supports the **AICS (AI Coding Session)** format, a standardized JSON-based interchange format for AI coding sessions. This format enables:

### Benefits

- **Tool Independence**: Switch between AI coding assistants (Cursor, Claude Code, GitHub Copilot, etc.) without losing history
- **Data Portability**: Export and import sessions in a vendor-neutral format
- **Archival**: Long-term preservation of AI conversations in a standardized format
- **Interoperability**: Share sessions with team members using different tools

### Usage Examples

#### Migrate from Crush to Another Tool

```bash
# Step 1: Export from Crush to AICS format
crush-md export-aics --db .crush/crush.db --out my-sessions.aics.json

# Step 2: Import the AICS file into your new tool (if supported)
# or convert to markdown for reference
crush-md import-aics --input my-sessions.aics.json --format markdown
```

#### Archive All Sessions

```bash
# Export all sessions to AICS format for long-term storage
crush-md export-aics --db .crush/crush.db --out archive-2024.aics.json --limit 1000
```

#### Share Sessions with Team

```bash
# Export specific sessions to AICS format
crush-md export-aics --db .crush/crush.db --out team-sessions.aics.json --limit 10

# Team member imports and converts to their preferred format
crush-md import-aics --input team-sessions.aics.json --format html
```

### Format Specification

For detailed information about the AICS format specification, see [AICS_FORMAT.md](AICS_FORMAT.md).

The AICS format is inspired by the HAR (HTTP Archive) format and provides:
- Standardized JSON structure
- ISO 8601 timestamps
- Support for multiple message types (text, tool calls, tool results)
- Flexible metadata system
- Version control for format evolution

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
│   └── markdown/
│       ├── renderer.go           # Markdown rendering
│       └── utils.go              # Utility functions
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

## Security Notes

- All database queries use parameterized statements to prevent SQL injection
- No user credentials or sensitive data are logged or exposed
- SQLite database is accessed in read-only mode for exports
- Input validation is performed on all user-provided parameters

## Troubleshooting

### Common Issues

**Database not found:**
```bash
# Check if the database exists
ls -la .crush/crush.db

# Verify database permissions
file .crush/crush.db
```

**No sessions found:**
```bash
# Verify database has sessions table
sqlite3 .crush/crush.db ".tables"
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