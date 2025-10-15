# Crush Session Explorer

A CLI tool for exporting Crush chat sessions from SQLite databases to Markdown format with YAML frontmatter.

Available in both **Python** and **Go** implementations with identical functionality.

## Overview

This tool allows you to extract and export individual chat sessions from Crush's SQLite database, converting them into well-formatted Markdown files with structured metadata. Perfect for archiving, documentation, or further processing of conversation data.

## Features

- 📊 **SQLite Integration**: Direct access to Crush database using built-in SQLite support
- 📝 **Markdown Export**: Clean conversion to Markdown with YAML frontmatter
- 🔍 **Interactive Session Selection**: Browse and select sessions interactively
- 📅 **Timestamp Formatting**: Automatic timestamp conversion to readable formats
- 🏷️ **Metadata Preservation**: Session metadata and message details preserved
- 🎯 **Type Safety**: Full type annotations (Python: pyright, Go: built-in)
- ✅ **Well Tested**: Comprehensive test coverage
- 🚀 **Dual Implementation**: Available in both Python and Go with identical CLI interface

## Requirements

### Python Implementation
- **Python 3.10+**
- **sqlite3** (included in Python standard library)

### Go Implementation  
- **Go 1.19+**
- **CGO enabled** (for SQLite driver)

## Installation

### Python Setup

```bash
git clone <repository-url>
cd crush-session-explorer
python -m venv .venv
source .venv/bin/activate  # On Windows: .venv\Scripts\activate
pip install -U pip
pip install -r requirements.txt
```

### Go Setup

```bash
git clone <repository-url>
cd crush-session-explorer
go mod download
make build
```

Or install directly:
```bash
go install ./cmd/crush-md
```

## Usage

The CLI interface is identical for both Python and Go implementations.

### Export a Specific Session

**Python:**
```bash
python -m crush_md export --db ./.crush/crush.db --session <session-id> --out output.md
```

**Go:**
```bash
./bin/crush-md export --db ./.crush/crush.db --session <session-id> --out output.md
# or if installed: crush-md export --db ./.crush/crush.db --session <session-id> --out output.md
```

### Interactive Session Selection

**Python:**
```bash
python -m crush_md export --db ./.crush/crush.db
```

**Go:**
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

| Option | Description | Default |
|--------|-------------|---------|
| `--db` | Path to the SQLite database | `.crush/crush.db` |
| `--session` | Specific session ID to export | Interactive selection |
| `--out` | Output Markdown file path | Auto-generated based on session |

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
├── crush_md/                     # Python implementation
│   ├── __init__.py               # Package initialization
│   ├── cli.py                    # Command-line interface and main logic
│   ├── db.py                     # SQLite database operations and models
│   └── markdown.py               # Markdown rendering and formatting
├── cmd/                          # Go implementation
│   └── crush-md/
│       └── main.go               # Go CLI entry point
├── internal/                     # Go internal packages
│   ├── cli/
│   │   └── export.go             # Export command implementation
│   ├── db/
│   │   ├── connection.go         # Database connection
│   │   ├── models.go             # Data models
│   │   └── queries.go            # Database queries
│   └── markdown/
│       ├── renderer.go           # Markdown rendering
│       └── utils.go              # Utility functions
├── tests/                        # Python tests
│   ├── test_db.py                # Database functionality tests
│   └── test_markdown.py          # Markdown rendering tests
├── bin/                          # Go build output (created by make build)
├── go.mod                        # Go module definition
├── go.sum                        # Go dependencies
├── requirements.txt              # Python dependencies
├── README.md                     # This file
└── Makefile                      # Build and development commands
```

## Development

### Python Development

**Running Tests:**
```bash
pytest -q
```

**Code Quality:**
```bash
ruff check .
ruff format .
pyright
```

### Go Development

**Building:**
```bash
make build          # Build for current platform
make build-all      # Build for all platforms
make dev            # Format, vet, test, and build
```

**Testing:**
```bash
make test           # Run tests
make test-coverage  # Run tests with coverage
```

**Code Quality:**
```bash
make fmt            # Format code
make vet            # Vet code
make lint           # Lint code (requires golangci-lint)
make check          # Run format, vet, and test
```

### Making Changes

The codebase follows these principles:

- **Security First**: Uses parameterized queries to prevent SQL injection
- **Type Safety**: Full type annotations for better IDE support and error catching
- **Clean Code**: Formatted with ruff, following Python best practices
- **Testable**: Modular design with comprehensive test coverage

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

## License

This project is provided as-is for educational and archival purposes.

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes following the existing code style
4. Add tests for new functionality
5. Ensure all tests pass (`pytest`)
6. Run linting (`ruff check . && ruff format .`)
7. Run type checking (`pyright`)
8. Commit your changes (`git commit -m 'feat: add amazing feature'`)
9. Push to the branch (`git push origin feature/amazing-feature`)
10. Open a Pull Request