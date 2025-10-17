# AICS Format Specification

## AI Coding Session Interchange Format (AICS)

Version: 1.0

### Overview

The **AICS (AI Coding Session)** format is a standardized JSON-based interchange format for AI coding sessions across different providers and tools. It is designed to facilitate:

- Migration between AI coding assistants (Cursor, Claude Code, GitHub Copilot, etc.)
- Archival and backup of AI conversation history
- Data portability and vendor independence
- Analysis and auditing of AI interactions

This format is inspired by the [HAR (HTTP Archive)](https://en.wikipedia.org/wiki/HAR_(file_format)) format, which provides a standard way to export HTTP transaction data.

### Design Goals

1. **Vendor Neutral**: Not tied to any specific AI provider or tool
2. **Extensible**: Can accommodate new features through metadata
3. **Human Readable**: JSON format that can be easily inspected
4. **Complete**: Preserves all essential information about sessions
5. **Verifiable**: Can be validated against a schema

### File Structure

An AICS file is a JSON document with the following structure:

```json
{
  "version": "1.0",
  "creator": {
    "name": "crush-session-explorer",
    "version": "v1.0.1",
    "comment": "Exported from Crush database"
  },
  "browser": {
    "name": "Crush",
    "version": "1.0",
    "comment": "Original AI coding tool"
  },
  "log": {
    "version": "1.0",
    "creator": { ... },
    "browser": { ... },
    "sessions": [ ... ]
  }
}
```

### Root Object

The root object contains:

- `version` (string, required): Format version (e.g., "1.0")
- `creator` (object, required): Information about the tool that created this archive
- `browser` (object, optional): Information about the original AI tool/provider
- `log` (object, required): Container for sessions and metadata

### Creator Object

```json
{
  "name": "tool-name",
  "version": "1.0.0",
  "comment": "Optional additional information"
}
```

### Browser Object

Represents the AI coding tool or provider:

```json
{
  "name": "Cursor",
  "version": "0.42.0",
  "comment": "Optional additional information"
}
```

### Log Object

Contains all sessions:

```json
{
  "version": "1.0",
  "creator": { ... },
  "browser": { ... },
  "sessions": [ ... ],
  "comment": "Optional additional information"
}
```

### Session Object

Represents a single AI coding session:

```json
{
  "id": "unique-session-id",
  "title": "Session Title",
  "startedAt": "2024-01-15T14:30:00Z",
  "updatedAt": "2024-01-15T15:45:00Z",
  "messages": [ ... ],
  "metadata": {
    "message_count": 12,
    "project": "my-project",
    "branch": "feature/new-feature"
  },
  "comment": "Optional additional information"
}
```

Fields:
- `id` (string, required): Unique identifier for the session
- `title` (string, optional): Human-readable session title
- `startedAt` (ISO 8601 timestamp, optional): When the session started
- `updatedAt` (ISO 8601 timestamp, optional): Last update time
- `messages` (array, required): Array of message objects
- `metadata` (object, optional): Flexible key-value pairs for additional data
- `comment` (string, optional): Additional information

### Message Object

Represents a single message in a session:

```json
{
  "id": "unique-message-id",
  "timestamp": "2024-01-15T14:30:00Z",
  "role": "user",
  "content": [
    {
      "type": "text",
      "text": "Can you help me refactor this function?"
    }
  ],
  "model": "claude-3-opus",
  "provider": "anthropic",
  "metadata": {},
  "comment": "Optional additional information"
}
```

Fields:
- `id` (string, required): Unique identifier for the message
- `timestamp` (ISO 8601 timestamp, optional): When the message was created
- `role` (string, required): Message role - one of:
  - `"user"`: User prompt
  - `"assistant"`: AI response
  - `"system"`: System message
  - `"tool"`: Tool execution message
- `content` (array, required): Array of content parts
- `model` (string, optional): AI model identifier (e.g., "gpt-4", "claude-3-opus")
- `provider` (string, optional): Provider name (e.g., "openai", "anthropic")
- `metadata` (object, optional): Additional message-specific data
- `comment` (string, optional): Additional information

### Content Object

Represents a content part within a message:

```json
{
  "type": "text",
  "text": "Content text here",
  "data": {},
  "mimeType": "text/plain",
  "encoding": "utf-8",
  "comment": "Optional additional information"
}
```

Fields:
- `type` (string, required): Content type - one of:
  - `"text"`: Plain text content
  - `"tool_call"`: Tool invocation
  - `"tool_result"`: Result from tool execution
  - `"code"`: Code snippet
  - `"image"`: Image content
- `text` (string, optional): Text content
- `data` (object, optional): Structured data for tool calls, results, etc.
- `mimeType` (string, optional): MIME type for content
- `encoding` (string, optional): Encoding for binary content (e.g., "base64")
- `comment` (string, optional): Additional information

### Metadata Object

A flexible key-value object for storing additional information:

```json
{
  "key1": "value1",
  "key2": 42,
  "key3": ["array", "of", "values"],
  "key4": {
    "nested": "object"
  }
}
```

### Timestamps

All timestamps in the AICS format use ISO 8601 format:

```
2024-01-15T14:30:00Z        // UTC
2024-01-15T14:30:00+00:00   // UTC with timezone
2024-01-15T09:30:00-05:00   // EST timezone
```

### Example AICS File

```json
{
  "version": "1.0",
  "creator": {
    "name": "crush-session-explorer",
    "version": "v1.0.1"
  },
  "browser": {
    "name": "Crush"
  },
  "log": {
    "version": "1.0",
    "creator": {
      "name": "crush-session-explorer",
      "version": "v1.0.1"
    },
    "sessions": [
      {
        "id": "abc123",
        "title": "Refactoring user authentication",
        "startedAt": "2024-01-15T14:30:00Z",
        "updatedAt": "2024-01-15T15:45:00Z",
        "messages": [
          {
            "id": "msg1",
            "timestamp": "2024-01-15T14:30:00Z",
            "role": "user",
            "content": [
              {
                "type": "text",
                "text": "Can you help me refactor the authentication code?"
              }
            ]
          },
          {
            "id": "msg2",
            "timestamp": "2024-01-15T14:31:00Z",
            "role": "assistant",
            "content": [
              {
                "type": "text",
                "text": "I'd be happy to help! Let me analyze the code..."
              }
            ],
            "model": "claude-3-opus",
            "provider": "anthropic"
          }
        ],
        "metadata": {
          "message_count": 2
        }
      }
    ]
  }
}
```

### Validation

An AICS file is valid if:

1. It is well-formed JSON
2. The `version` field is present and matches a supported version
3. The `creator` object has a `name` field
4. The `log` object is present
5. Each session has an `id` and at least one message
6. Each message has an `id`, `role`, and at least one content part
7. All timestamps are valid ISO 8601 format

### Usage Examples

#### Exporting from a tool

```bash
crush-md export-aics --db .crush/crush.db --out sessions.aics.json
```

#### Importing to another tool

```bash
crush-md import-aics --input sessions.aics.json --format markdown --out ./imported/
```

#### Converting between formats

```bash
# Export from Crush to AICS
crush-md export-aics --db .crush/crush.db --out crush-sessions.aics.json

# Import AICS to Markdown
crush-md import-aics --input crush-sessions.aics.json --format markdown
```

### Extensions

Tools may add custom fields to the `metadata` objects at any level. Custom fields should use a vendor-specific prefix to avoid conflicts:

```json
{
  "metadata": {
    "cursor_project_id": "proj123",
    "claude_code_workspace": "/path/to/workspace",
    "custom_field": "value"
  }
}
```

### Version History

#### Version 1.0 (Initial Release)

- Basic structure for sessions and messages
- Support for text, tool calls, and tool results
- Flexible metadata system
- ISO 8601 timestamps

### References

- [HAR (HTTP Archive) Format](https://en.wikipedia.org/wiki/HAR_(file_format))
- [JSON Schema](https://json-schema.org/)
- [ISO 8601 Date/Time Format](https://en.wikipedia.org/wiki/ISO_8601)

### Contributing

To propose extensions or changes to the AICS format, please:

1. Open an issue describing the use case
2. Provide example data structures
3. Ensure backward compatibility
4. Update this specification document

### License

The AICS format specification is released under CC0 1.0 Universal (public domain).
