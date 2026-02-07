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

// GetLastUpdate returns the parsed LastUpdated timestamp from the JSON, 
// falling back to the file system modification time if parsing fails.
func (s Session) GetLastUpdate() time.Time {
	t, err := time.Parse(time.RFC3339, s.LastUpdated)
	if err != nil {
		return s.FileLastUpdate
	}
	return t
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

// GetSession aggregates all messages for a specific session ID in a project.
// Gemini CLI may split a single logical session across multiple files if it spans long periods.
func GetSession(rootDir, projectID, sessionID string) (Session, error) {
	allSessions, err := ReadSessions(rootDir, projectID)
	if err != nil {
		return Session{}, err
	}

	var result Session
	found := false

	for _, s := range allSessions {
		if s.ID == sessionID {
			if !found {
				result = s
				found = true
			} else {
				// Aggregate messages
				result.Messages = append(result.Messages, s.Messages...)
				// Update timestamps if necessary (assuming they are already sorted or we sort later)
			}
		}
	}

	if !found {
		return Session{}, fmt.Errorf("session %s not found in project %s", sessionID, projectID)
	}

	// Ensure messages are sorted by timestamp if they came from multiple files
	// (Gemini CLI session files are usually chronological, but aggregation might need care)
	return result, nil
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

// DeleteProject removes the entire project directory from Gemini storage.
func DeleteProject(rootDir, projectID string) error {
	path := filepath.Join(rootDir, projectID)
	return os.RemoveAll(path)
}

// DeleteSession removes all files associated with a specific session ID.
func DeleteSession(rootDir, projectID, sessionID string) error {
	allSessions, err := ReadSessions(rootDir, projectID)
	if err != nil {
		return err
	}

	for _, s := range allSessions {
		if s.ID == sessionID {
			if err := os.Remove(s.FilePath); err != nil {
				return err
			}
		}
	}
	return nil
}

// MoveSession relocates a session to a different project.
func MoveSession(rootDir, oldProjectID, newProjectID, sessionID string) error {
	if oldProjectID == newProjectID {
		return nil
	}

	allSessions, err := ReadSessions(rootDir, oldProjectID)
	if err != nil {
		return err
	}

	for _, s := range allSessions {
		if s.ID == sessionID {
			oldPath := s.FilePath
			s.ProjectHash = newProjectID
			
			if err := WriteSession(rootDir, newProjectID, s); err != nil {
				return err
			}
			
			if err := os.Remove(oldPath); err != nil {
				// We could try to roll back but it's complex. 
				// At least the new file is written.
				return fmt.Errorf("failed to remove old session file: %w", err)
			}
		}
	}
	return nil
}

// MoveProject migrates a project to a new directory path.
// It renames the storage directory and updates the projectHash in all sessions.
func MoveProject(rootDir, oldID, newPath string) (string, error) {
	newID, err := HashProjectID(newPath)
	if err != nil {
		return "", err
	}

	if oldID == newID {
		return newID, nil // No change needed
	}

	oldPath := filepath.Join(rootDir, oldID)
	newStoragePath := filepath.Join(rootDir, newID)

	// 1. Rename the project directory
	if err := os.Rename(oldPath, newStoragePath); err != nil {
		return "", err
	}

	// 2. Update all sessions in the new directory
	sessions, err := ReadSessions(rootDir, newID)
	if err != nil {
		return newID, fmt.Errorf("failed to read sessions after move: %w", err)
	}

	for _, s := range sessions {
		// Store the old file path before we potentially rename/overwrite
		oldFilePath := s.FilePath

		s.ProjectHash = newID
		if err := WriteSession(rootDir, newID, s); err != nil {
			return newID, fmt.Errorf("failed to update session %s: %w", s.ID, err)
		}

		// If the filename changed (because the short ID changed), remove the old one
		newShortID := newID[:8]
		if !strings.Contains(filepath.Base(oldFilePath), newShortID) {
			_ = os.Remove(oldFilePath)
		}
	}

	return newID, nil
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