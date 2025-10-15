package markdown

// FormatTimestamp formats a timestamp for display (exported for CLI use)
func FormatTimestamp(ts *string) string {
	return formatTimestamp(ts)
}