package providers

import (
	"crush-session-explorer/internal/db"
)

// Provider defines the interface for session providers from different AI code tools
type Provider interface {
	// Name returns the provider name (e.g., "crush", "claude-code", "cursor")
	Name() string

	// Discover checks if this provider's data exists on the system
	Discover() (bool, error)

	// ListSessions retrieves available sessions from this provider
	ListSessions(limit int) ([]db.Session, error)

	// FetchSession retrieves a specific session by ID
	FetchSession(sessionID string) (*db.Session, error)

	// ListMessages retrieves messages for a session
	ListMessages(sessionID string) ([]db.ParsedMessage, error)
}

// DiscoverAllProviders finds all available providers on the system
func DiscoverAllProviders() []Provider {
	allProviders := []Provider{
		NewCrushProvider(),
		NewClaudeProvider(),
	}

	var available []Provider
	for _, provider := range allProviders {
		if found, err := provider.Discover(); err == nil && found {
			available = append(available, provider)
		}
	}

	return available
}

// GetProvider returns a specific provider by name
func GetProvider(name string) Provider {
	providers := map[string]Provider{
		"crush":       NewCrushProvider(),
		"claude-code": NewClaudeProvider(),
		"claude":      NewClaudeProvider(),
	}

	return providers[name]
}
