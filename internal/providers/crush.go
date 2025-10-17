package providers

import (
	"crush-session-explorer/internal/db"
	"database/sql"
	"os"
	"path/filepath"
)

// CrushProvider implements the Provider interface for Crush sessions
type CrushProvider struct {
	dbPath string
	conn   *sql.DB
}

// NewCrushProvider creates a new Crush provider instance
func NewCrushProvider() *CrushProvider {
	return &CrushProvider{
		dbPath: ".crush/crush.db",
	}
}

// NewCrushProviderWithPath creates a new Crush provider with custom db path
func NewCrushProviderWithPath(dbPath string) *CrushProvider {
	return &CrushProvider{
		dbPath: dbPath,
	}
}

// Name returns the provider name
func (p *CrushProvider) Name() string {
	return "crush"
}

// Discover checks if Crush database exists
func (p *CrushProvider) Discover() (bool, error) {
	// Check if database file exists
	if _, err := os.Stat(p.dbPath); os.IsNotExist(err) {
		return false, nil
	}

	// Try to connect to verify it's a valid database
	conn, err := db.Connect(p.dbPath)
	if err != nil {
		return false, nil
	}
	defer conn.Close()

	return true, nil
}

// getConnection returns or creates a database connection
func (p *CrushProvider) getConnection() (*sql.DB, error) {
	if p.conn == nil {
		conn, err := db.Connect(p.dbPath)
		if err != nil {
			return nil, err
		}
		p.conn = conn
	}
	return p.conn, nil
}

// ListSessions retrieves sessions from Crush database
func (p *CrushProvider) ListSessions(limit int) ([]db.Session, error) {
	conn, err := p.getConnection()
	if err != nil {
		return nil, err
	}

	sessions, err := db.ListSessions(conn, limit)
	if err != nil {
		return nil, err
	}

	// Add provider metadata to each session
	for i := range sessions {
		provider := "crush"
		sessions[i].Metadata = &provider
	}

	return sessions, nil
}

// FetchSession retrieves a specific session
func (p *CrushProvider) FetchSession(sessionID string) (*db.Session, error) {
	conn, err := p.getConnection()
	if err != nil {
		return nil, err
	}

	session, err := db.FetchSession(conn, sessionID)
	if err != nil {
		return nil, err
	}

	// Add provider metadata
	provider := "crush"
	session.Metadata = &provider

	return session, nil
}

// ListMessages retrieves messages for a session
func (p *CrushProvider) ListMessages(sessionID string) ([]db.ParsedMessage, error) {
	conn, err := p.getConnection()
	if err != nil {
		return nil, err
	}

	return db.ListMessages(conn, sessionID)
}

// SetDBPath allows setting a custom database path
func (p *CrushProvider) SetDBPath(path string) {
	// Expand home directory if needed
	if len(path) > 0 && path[0] == '~' {
		home, err := os.UserHomeDir()
		if err == nil {
			path = filepath.Join(home, path[1:])
		}
	}
	p.dbPath = path
	// Close existing connection since path changed
	if p.conn != nil {
		p.conn.Close()
		p.conn = nil
	}
}

// Close closes the database connection
func (p *CrushProvider) Close() error {
	if p.conn != nil {
		err := p.conn.Close()
		p.conn = nil
		return err
	}
	return nil
}
