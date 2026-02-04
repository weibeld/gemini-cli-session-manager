package gemini

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Constants for Gemini CLI storage structure.
const (
	SessionDir    = "chats"
	SessionPrefix = "session-"
	SessionSuffix = ".json"
)

// Thought represents a chain-of-thought step.
type Thought struct {
	Subject     string    `json:"subject"`
	Description string    `json:"description"`
	Timestamp   string    `json:"timestamp"`
}

// TokenStats represents usage statistics.
type TokenStats struct {
	Input    int `json:"input"`
	Output   int `json:"output"`
	Cached   int `json:"cached"`
	Thoughts int `json:"thoughts"`
	Tool     int `json:"tool"`
	Total    int `json:"total"`
}

// Message represents a single interaction in a session.
type Message struct {
	ID        string     `json:"id"`
	Timestamp string     `json:"timestamp"`
	Type      string     `json:"type"` // "user" or "gemini"
	Content   string     `json:"content"`
	Thoughts  []Thought  `json:"thoughts,omitempty"`
	Tokens    *TokenStats `json:"tokens,omitempty"`
	Model     string     `json:"model,omitempty"`
}

// Session represents the content of a Gemini CLI session file.
type Session struct {
	ID          string    `json:"sessionId"`
	ProjectHash string    `json:"projectHash"`
	StartTime   string    `json:"startTime"`
	LastUpdated string    `json:"lastUpdated"`
	Messages    []Message `json:"messages"`
	
	// Metadata not in JSON but extracted from file system
	FileLastUpdate time.Time `json:"-"`
	FilePath       string    `json:"-"`
}

// HashProjectID returns the SHA-256 hash of the standardized absolute path.
func HashProjectID(path string) (string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256([]byte(abs))
	return hex.EncodeToString(hash[:]), nil
}

// ListProjectIDs discovers all project hash directories in the root.
func ListProjectIDs(rootDir string) ([]string, error) {
	entries, err := os.ReadDir(rootDir)
	if err != nil {
		return nil, err
	}

	var ids []string
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		if len(name) == 64 { // Expecting full SHA-256
			ids = append(ids, name)
		}
	}
	return ids, nil
}

// ReadSessions parses all session files for a specific project.
func ReadSessions(rootDir, projectID string) ([]Session, error) {
	sessionPath := filepath.Join(rootDir, projectID, SessionDir)
	entries, err := os.ReadDir(sessionPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var sessions []Session
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasPrefix(entry.Name(), SessionPrefix) || !strings.HasSuffix(entry.Name(), SessionSuffix) {
			continue
		}

		path := filepath.Join(sessionPath, entry.Name())
		s, err := parseSessionFile(path)
		if err != nil {
			continue
		}
		s.FilePath = path
		sessions = append(sessions, s)
	}
	return sessions, nil
}

// WriteSession marshals and writes a session to the Gemini storage structure.
// It generates a filename in the format: session-YYYY-MM-DDTHH-mm-shortID.json
func WriteSession(rootDir, projectID string, s Session) error {
	sessionPath := filepath.Join(rootDir, projectID, SessionDir)
	if err := os.MkdirAll(sessionPath, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}

	// Format timestamp: 2026-02-02T12:55 -> 2026-02-02T12-55
	ts := strings.ReplaceAll(s.StartTime[:16], ":", "-")
	
	// Short ID: first 8 chars of sessionId
	shortID := s.ID
	if len(shortID) > 8 {
		shortID = shortID[:8]
	}

	filename := fmt.Sprintf("%s%s-%s%s", SessionPrefix, ts, shortID, SessionSuffix)
	return os.WriteFile(filepath.Join(sessionPath, filename), data, 0644)
}

func parseSessionFile(path string) (Session, error) {
	info, err := os.Stat(path)
	if err != nil {
		return Session{}, err
	}

	if info.Size() > 10*1024*1024 {
		return Session{}, fmt.Errorf("session file too large: %s", path)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return Session{}, err
	}

	var s Session
	if err := json.Unmarshal(data, &s); err != nil {
		return Session{}, err
	}
	s.FileLastUpdate = info.ModTime()
	return s, nil
}