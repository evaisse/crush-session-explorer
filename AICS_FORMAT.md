# AICS Format Specification

## AI Coding Session Interchange Format (AICS)

Version: 1.0

### Overview

The **AICS (AI Coding Session)** format is a standardized JSON-based interchange format for AI coding sessions across different providers and tools. It is designed to facilitate:

- Migration between AI coding assistants (Cursor, Claude Code, GitHub Copilot, etc.)
- Archival and backup of AI conversation history
- Data portability and vendor independence
- Analysis and auditing of AI interactions

This format is inspired by the [HAR (HTTP Archive)](<https://en.wikipedia.org/wiki/HAR_(file_format)>) format, which provides a standard way to export HTTP transaction data.

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
  "id": "01234567-89ab-7def-0123-456789abcdef",
  "clientId": "fedcba98-7654-3210-fedc-ba9876543210",
  "title": "Session Title",
  "startedAt": "2024-01-15T14:30:00Z",
  "updatedAt": "2024-01-15T15:45:00Z",
  "messages": [ ... ],
  "gitRefs": {
    "branches": ["main", "feature/new-api"],
    "issues": ["#123", "org/repo#456"],
    "commits": ["abc1234"],
    "tags": ["v1.0.0"],
    "repos": ["owner/repository"]
  },
  "metadata": {
    "message_count": 12,
    "project": "my-project"
  },
  "comment": "Optional additional information"
}
```

Fields:
- `id` (string, required): Unique identifier for the session (UUID v7 recommended)
- `clientId` (string, optional): Client/machine identifier for session grouping
- `title` (string, optional): Human-readable session title
- `startedAt` (ISO 8601 timestamp, optional): When the session started
- `updatedAt` (ISO 8601 timestamp, optional): Last update time
- `messages` (array, required): Array of message objects
- `gitRefs` (object, optional): Git references mentioned during the session
- `metadata` (object, optional): Flexible key-value pairs for additional data
- `comment` (string, optional): Additional information

**Note on Session IDs**: It is recommended to use UUID v7 for session IDs. UUID v7 provides:
- Time-ordered identifiers (sortable by creation time)
- Globally unique across all systems
- Better database indexing performance
- Embedded timestamp information

### GitRefs Object

Tracks git-related references mentioned during a session:

```json
{
  "branches": ["main", "feature/new-api", "develop"],
  "issues": ["#123", "owner/repo#456", "GH-789"],
  "commits": ["abc1234", "def5678901234567890123456789012345678901"],
  "tags": ["v1.0.0", "release-2024-01"],
  "repos": ["owner/repository", "org/project"]
}
```

Fields:
- `branches` (array of strings, optional): Git branch names mentioned (e.g., "main", "feature/auth")
- `issues` (array of strings, optional): Issue or PR references (e.g., "#123", "owner/repo#456")
- `commits` (array of strings, optional): Commit SHAs (short or full)
- `tags` (array of strings, optional): Git tag names (e.g., "v1.0.0")
- `repos` (array of strings, optional): Repository identifiers (e.g., "owner/repo")

This field helps track the development context and makes it easier to:
- Link sessions to specific branches or features
- Associate conversations with issues or pull requests
- Track which commits were discussed
- Maintain context when switching between projects

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
  "mcp": {
    "version": "1.0",
    "tools": [
      {
        "name": "read_file",
        "input": {"path": "/src/main.go"},
        "output": "file content..."
      }
    ],
    "resources": [
      {
        "uri": "file:///workspace/src/main.go",
        "name": "main.go",
        "mimeType": "text/x-go"
      }
    ]
  },
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
- `mcp` (object, optional): Model Context Protocol information
- `metadata` (object, optional): Additional message-specific data
- `comment` (string, optional): Additional information

### Model Context Protocol (MCP)

The `mcp` field contains information about Model Context Protocol usage in AI coding tools. MCP is a protocol that allows AI assistants to interact with external tools, resources, and prompts in a standardized way.

**MCP Object Structure:**

```json
{
  "version": "1.0",
  "tools": [
    {
      "name": "read_file",
      "description": "Read contents of a file",
      "input": {
        "path": "/src/utils.js"
      },
      "output": "// File contents here..."
    }
  ],
  "resources": [
    {
      "uri": "file:///workspace/project/src/main.go",
      "name": "main.go",
      "description": "Main application file",
      "mimeType": "text/x-go",
      "metadata": {
        "size": 1024,
        "lastModified": "2024-01-15T14:30:00Z"
      }
    }
  ],
  "prompts": [
    {
      "name": "code_review",
      "description": "Review code for issues",
      "arguments": {
        "language": "go",
        "focus": "performance"
      }
    }
  ],
  "metadata": {
    "serverName": "my-mcp-server",
    "capabilities": ["tools", "resources", "prompts"]
  }
}
```

**MCP Fields:**
- `version` (string, optional): MCP protocol version (e.g., "1.0")
- `tools` (array, optional): Array of MCP tool invocations
- `resources` (array, optional): Array of MCP resources accessed
- `prompts` (array, optional): Array of MCP prompts used
- `metadata` (object, optional): Additional MCP-specific metadata

**MCP Tool Object:**
- `name` (string, required): Tool name
- `description` (string, optional): Tool description
- `input` (object, optional): Tool input parameters
- `output` (any, optional): Tool output/result

**MCP Resource Object:**
- `uri` (string, required): Resource URI (e.g., "file:///path", "git://repo")
- `name` (string, optional): Human-readable resource name
- `description` (string, optional): Resource description
- `mimeType` (string, optional): Resource MIME type
- `metadata` (object, optional): Additional resource metadata

**MCP Prompt Object:**
- `name` (string, required): Prompt name
- `description` (string, optional): Prompt description
- `arguments` (object, optional): Prompt arguments/parameters

**MCP Use Cases:**

The MCP field enables tracking of:
- **Tool Usage**: Which tools were called (file operations, git commands, search, etc.)
- **Resource Access**: What files, databases, or APIs were accessed
- **Prompt Chains**: Which prompts were used to guide the AI
- **Context Tracking**: Full context of how the AI assistant interacted with the environment

This is particularly useful for:
- Debugging AI assistant behavior
- Reproducing AI sessions in different environments
- Understanding what data the AI accessed
- Auditing AI tool usage
- Migrating sessions between different AI coding tools that support MCP

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
        "id": "01234567-89ab-7def-0123-456789abcdef",
        "clientId": "fedcba98-7654-3210-fedc-ba9876543210",
        "title": "Refactoring user authentication",
        "startedAt": "2024-01-15T14:30:00Z",
        "updatedAt": "2024-01-15T15:45:00Z",
        "gitRefs": {
          "branches": ["feature/auth-refactor"],
          "issues": ["#234"],
          "commits": ["abc1234"],
          "repos": ["myorg/auth-service"]
        },
        "messages": [
          {
            "id": "msg1",
            "timestamp": "2024-01-15T14:30:00Z",
            "role": "user",
            "content": [
              {
                "type": "text",
                "text": "Can you help me refactor the authentication code in #234?"
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
            "provider": "anthropic",
            "mcp": {
              "version": "1.0",
              "tools": [
                {
                  "name": "read_file",
                  "input": {"path": "src/auth.go"},
                  "output": "package auth..."
                }
              ],
              "resources": [
                {
                  "uri": "file:///workspace/src/auth.go",
                  "name": "auth.go",
                  "mimeType": "text/x-go"
                }
              ]
            }
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

- [HAR (HTTP Archive) Format](<https://en.wikipedia.org/wiki/HAR_(file_format)>)
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
