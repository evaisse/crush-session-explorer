# Crush Session Explorer

A Python CLI tool for exporting Crush chat sessions from SQLite databases to Markdown format with YAML frontmatter.

## Overview

This tool allows you to extract and export individual chat sessions from Crush's SQLite database, converting them into well-formatted Markdown files with structured metadata. Perfect for archiving, documentation, or further processing of conversation data.

## Features

- 📊 **SQLite Integration**: Direct access to Crush database using Python's built-in sqlite3
- 📝 **Markdown Export**: Clean conversion to Markdown with YAML frontmatter
- 🔍 **Interactive Session Selection**: Browse and select sessions interactively
- 📅 **Timestamp Formatting**: Automatic timestamp conversion to readable formats
- 🏷️ **Metadata Preservation**: Session metadata and message details preserved
- 🎯 **Type Safety**: Full type annotations with pyright compatibility
- ✅ **Well Tested**: Comprehensive test suite with pytest

## Requirements

- **Python 3.10+**
- **sqlite3** (included in Python standard library)

### Development Dependencies

```bash
pip install ruff pyright pytest
```

## Installation

Clone the repository and set up the environment:

```bash
git clone <repository-url>
cd crush-session-explorer
python -m venv .venv
source .venv/bin/activate  # On Windows: .venv\Scripts\activate
pip install -U pip
pip install -r requirements.txt
```

## Usage

### Export a Specific Session

```bash
python -m crush_md export --db ./.crush/crush.db --session <session-id> --out output.md
```

### Interactive Session Selection

```bash
python -m crush_md export --db ./.crush/crush.db
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
├── crush_md/
│   ├── __init__.py          # Package initialization
│   ├── cli.py               # Command-line interface and main logic
│   ├── db.py                # SQLite database operations and models
│   └── markdown.py          # Markdown rendering and formatting
├── tests/
│   ├── test_db.py           # Database functionality tests
│   └── test_markdown.py     # Markdown rendering tests
├── requirements.txt         # Development dependencies
├── README.md               # This file
└── Makefile                # Build and development commands
```

## Development

### Running Tests

```bash
pytest -q
```

### Code Quality

**Linting and Formatting:**
```bash
ruff check .
ruff format .
```

**Type Checking:**
```bash
pyright
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