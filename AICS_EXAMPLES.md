# AICS Format Examples

This document provides practical examples of using the AICS (AI Coding Session) interchange format for migrating sessions between different AI coding tools.

## Example 1: Migrating from Crush to Another Tool

### Scenario
You've been using Crush for a year and have accumulated many valuable coding sessions. Now you want to switch to Cursor or Claude Code but want to preserve your session history.

### Solution

```bash
# Step 1: Export all sessions from Crush to AICS format
crush-md export-aics \
  --db ~/.crush/crush.db \
  --out my-crush-sessions.aics.json \
  --provider "Crush" \
  --limit 500

# Output:
# Found 150 sessions to export
# ✅ Exported 150 sessions to my-crush-sessions.aics.json
# 📊 Format: AICS v1.0 (AI Coding Session Interchange Format)
# 💡 This file can be imported into other AI coding tools that support AICS

# Step 2: Convert to markdown for reference or import into new tool
crush-md import-aics \
  --input my-crush-sessions.aics.json \
  --format markdown \
  --out ~/Documents/ai-sessions-archive/

# Output:
# 📥 Importing from my-crush-sessions.aics.json...
# ✅ Successfully imported AICS archive
# 📊 Format version: 1.0
# 🔧 Created by: crush-session-explorer v1.0.1
# 🌐 Original tool: Crush
# 📝 Sessions: 150
# 
# 📤 Exporting sessions to markdown format...
#   ✓ 2024-01-15_14-30_refactoring-auth.md
#   ✓ 2024-01-16_10-00_database-optimization.md
#   ... (148 more)
# 
# ✅ Successfully exported 150/150 sessions to ~/Documents/ai-sessions-archive/
```

## Example 2: Sharing Sessions with Team Members

### Scenario
Your team uses different AI coding tools. You want to share some helpful sessions with your teammates.

### Solution

```bash
# Export specific recent sessions
crush-md export-aics \
  --db ~/.crush/crush.db \
  --out team-sessions.aics.json \
  --limit 10

# Share the AICS file with team members
# They can import it regardless of their AI tool:

# Team member using Cursor (converts to markdown)
crush-md import-aics \
  --input team-sessions.aics.json \
  --format markdown \
  --out ./shared-sessions/

# Team member who prefers HTML format
crush-md import-aics \
  --input team-sessions.aics.json \
  --format html \
  --out ./shared-sessions/
```

## Example 3: Long-term Archival

### Scenario
You want to preserve your AI coding sessions for future reference, compliance, or training purposes.

### Solution

```bash
# Create monthly archives
crush-md export-aics \
  --db ~/.crush/crush.db \
  --out archives/sessions-2024-01.aics.json \
  --limit 1000

# Store the AICS file in version control or backup system
git add archives/sessions-2024-01.aics.json
git commit -m "Archive AI coding sessions for January 2024"

# Later, when needed, convert to readable format
crush-md import-aics \
  --input archives/sessions-2024-01.aics.json \
  --format html \
  --out ./review/january-2024/
```

## Example 4: Testing a New AI Tool

### Scenario
You want to test a new AI coding tool but don't want to lose access to your existing sessions.

### Solution

```bash
# Export your current sessions to AICS
crush-md export-aics \
  --db ~/.crush/crush.db \
  --out backup-sessions.aics.json

# Try the new tool for a few weeks
# ... (use new tool) ...

# If you want to go back, your sessions are preserved in AICS format
# You can always convert them to any format supported by your tools

# Generate HTML reports of your old sessions for reference
crush-md import-aics \
  --input backup-sessions.aics.json \
  --format html \
  --out ./session-reference/
```

## Example 5: Cross-Tool Workflow

### Scenario
You use different AI tools for different projects and want to consolidate all your session history.

### Solution

```bash
# Export from Crush (Tool A)
crush-md export-aics \
  --db ~/.crush/crush.db \
  --out crush-sessions.aics.json \
  --provider "Crush"

# Export from Cursor (Tool B) - if supported
# cursor-export --format aics --out cursor-sessions.aics.json

# Merge AICS files (manual or with script)
# Then import the consolidated sessions

crush-md import-aics \
  --input consolidated-sessions.aics.json \
  --format markdown \
  --out ./all-sessions/
```

## AICS File Structure Example

Here's what a typical AICS file looks like:

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
    "comment": "Original AI coding tool"
  },
  "log": {
    "version": "1.0",
    "creator": {
      "name": "crush-session-explorer",
      "version": "v1.0.1"
    },
    "browser": {
      "name": "Crush"
    },
    "sessions": [
      {
        "id": "session-abc123",
        "title": "Refactoring Authentication Module",
        "startedAt": "2024-01-15T14:30:00Z",
        "updatedAt": "2024-01-15T16:45:00Z",
        "messages": [
          {
            "id": "msg-001",
            "timestamp": "2024-01-15T14:30:00Z",
            "role": "user",
            "content": [
              {
                "type": "text",
                "text": "How can I refactor this authentication code to be more secure?"
              }
            ]
          },
          {
            "id": "msg-002",
            "timestamp": "2024-01-15T14:32:00Z",
            "role": "assistant",
            "content": [
              {
                "type": "text",
                "text": "I'll help you improve the security of your authentication code. Here are several recommendations..."
              }
            ],
            "model": "claude-3-opus",
            "provider": "anthropic"
          }
        ],
        "metadata": {
          "message_count": 2,
          "project": "auth-service",
          "language": "python"
        }
      }
    ]
  }
}
```

## Integration with Other Tools

### For Tool Developers

If you're developing an AI coding tool and want to support AICS format:

1. **Read the specification**: See [AICS_FORMAT.md](AICS_FORMAT.md)
2. **Implement export**: Convert your internal format to AICS JSON
3. **Implement import**: Parse AICS JSON and convert to your internal format
4. **Test with examples**: Use the sample files provided

### Sample Import Code (Conceptual)

```go
// Import AICS file
archive, err := interchange.ImportFromFile("sessions.aics.json")
if err != nil {
    log.Fatal(err)
}

// Validate the archive
if err := interchange.ValidateArchive(archive); err != nil {
    log.Fatal(err)
}

// Convert to your tool's format
for _, session := range archive.Log.Sessions {
    // Process each session
    for _, message := range session.Messages {
        // Import message into your database
        importMessage(message)
    }
}
```

## Benefits Summary

### For Users
- ✅ **Freedom to switch tools** without losing history
- ✅ **Portable data** in vendor-neutral format
- ✅ **Long-term preservation** of valuable conversations
- ✅ **Easy sharing** with team members

### For Developers
- ✅ **Standard format** reduces custom integration work
- ✅ **Interoperability** with other AI coding tools
- ✅ **Clear specification** makes implementation straightforward
- ✅ **Community support** for format evolution

## Next Steps

1. Try exporting your sessions: `crush-md export-aics --help`
2. Read the format specification: [AICS_FORMAT.md](AICS_FORMAT.md)
3. Import sessions from other tools: `crush-md import-aics --help`
4. Share your feedback and contribute to the format evolution

## Resources

- [AICS Format Specification](AICS_FORMAT.md)
- [HAR Format Specification](https://en.wikipedia.org/wiki/HAR_(file_format)) (inspiration)
- [ISO 8601 Date/Time Format](https://en.wikipedia.org/wiki/ISO_8601)
