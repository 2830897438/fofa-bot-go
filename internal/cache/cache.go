package cache

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

const (
	CacheDir        = "fofa_cache"
	HistoryFile     = "history.json"
	MaxHistorySize  = 50
	CacheExpiration = 24 * time.Hour
)

// QueryCache represents a cached query
type QueryCache struct {
	QueryText string    `json:"query_text"`
	Timestamp time.Time `json:"timestamp"`
	FilePath  string    `json:"file_path"`
	Count     int       `json:"count"`
}

// History represents query history
type History struct {
	Queries []QueryCache `json:"queries"`
}

// Manager manages cache operations
type Manager struct {
	historyFile string
	cacheDir    string
}

// NewManager creates a new cache manager
func NewManager() *Manager {
	return &Manager{
		historyFile: HistoryFile,
		cacheDir:    CacheDir,
	}
}

// Init initializes cache directories
func (m *Manager) Init() error {
	return os.MkdirAll(m.cacheDir, 0755)
}

// LoadHistory loads query history
func (m *Manager) LoadHistory() (*History, error) {
	data, err := os.ReadFile(m.historyFile)
	if err != nil {
		if os.IsNotExist(err) {
			return &History{Queries: []QueryCache{}}, nil
		}
		return nil, err
	}

	var history History
	if err := json.Unmarshal(data, &history); err != nil {
		return nil, err
	}

	return &history, nil
}

// SaveHistory saves query history
func (m *Manager) SaveHistory(history *History) error {
	data, err := json.MarshalIndent(history, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(m.historyFile, data, 0644)
}

// AddQuery adds a query to history
func (m *Manager) AddQuery(query string, filePath string, count int) error {
	history, err := m.LoadHistory()
	if err != nil {
		return err
	}

	// Remove existing entry if present
	for i, q := range history.Queries {
		if q.QueryText == query {
			history.Queries = append(history.Queries[:i], history.Queries[i+1:]...)
			break
		}
	}

	// Add new entry at the beginning
	cache := QueryCache{
		QueryText: query,
		Timestamp: time.Now(),
		FilePath:  filePath,
		Count:     count,
	}
	history.Queries = append([]QueryCache{cache}, history.Queries...)

	// Limit history size
	if len(history.Queries) > MaxHistorySize {
		history.Queries = history.Queries[:MaxHistorySize]
	}

	return m.SaveHistory(history)
}

// FindCache finds a cached query
func (m *Manager) FindCache(query string) *QueryCache {
	history, err := m.LoadHistory()
	if err != nil {
		return nil
	}

	for _, q := range history.Queries {
		if q.QueryText == query {
			// Check if file exists
			if _, err := os.Stat(q.FilePath); err == nil {
				return &q
			}
		}
	}

	return nil
}

// GetCachePath returns the full path for a cache file
func (m *Manager) GetCachePath(filename string) string {
	return filepath.Join(m.cacheDir, filename)
}
