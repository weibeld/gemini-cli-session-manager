package scanner

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Session represents metadata for a single Gemini session.
type Session struct {
	ID           string    `json:"sessionId"`
	MessageCount int       `json:"messageCount"`
	LastUpdate   time.Time `json:"lastUpdate"`
}

// ProjectData aggregates sessions for a specific project hash.
type ProjectData struct {
	ID       string
	Sessions []Session
}

// Scanner handles discovery of Gemini sessions.
type Scanner struct {
	RootDir string
}

// NewScanner creates a scanner pointing to the default Gemini tmp directory.
func NewScanner() (*Scanner, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	return &Scanner{
		RootDir: filepath.Join(home, ".gemini", "tmp"),
	}, nil
}

// Scan discovery all projects and their sessions.
func (s *Scanner) Scan() ([]ProjectData, error) {
	entries, err := os.ReadDir(s.RootDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var projects []ProjectData

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		projectID := entry.Name()
		// Basic validation: Project IDs are hex hashes (typically SHA-256)
		if len(projectID) < 8 {
			continue
		}

		project := ProjectData{ID: projectID}
		sessions, err := s.scanSessions(filepath.Join(s.RootDir, projectID, "chats"))
		if err != nil {
			// Log error or handle as needed, for now we just skip or partial return
			continue
		}
		project.Sessions = sessions
		projects = append(projects, project)
	}

	return projects, nil
}

func (s *Scanner) scanSessions(chatsDir string) ([]Session, error) {
	entries, err := os.ReadDir(chatsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	sessionMap := make(map[string]*Session)

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasPrefix(entry.Name(), "session-") || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		path := filepath.Join(chatsDir, entry.Name())
		session, err := s.parseSessionFile(path)
		if err != nil {
			continue
		}

		// Since a session can have multiple files (sub-agents), we merge them.
		if existing, ok := sessionMap[session.ID]; ok {
			existing.MessageCount += session.MessageCount
			if session.LastUpdate.After(existing.LastUpdate) {
				existing.LastUpdate = session.LastUpdate
			}
		} else {
			sessionMap[session.ID] = &session
		}
	}

	var sessions []Session
	for _, sess := range sessionMap {
		sessions = append(sessions, *sess)
	}

	// Sort sessions by last update descending
	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].LastUpdate.After(sessions[j].LastUpdate)
	})

	return sessions, nil
}

type rawSession struct {
	SessionID string `json:"sessionId"`
	Messages  []any  `json:"messages"`
}

func (s *Scanner) parseSessionFile(path string) (Session, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Session{}, err
	}

	var raw rawSession
	if err := json.Unmarshal(data, &raw); err != nil {
		return Session{}, err
	}

	info, err := os.Stat(path)
	if err != nil {
		return Session{}, err
	}

	return Session{
		ID:           raw.SessionID,
		MessageCount: len(raw.Messages),
		LastUpdate:   info.ModTime(),
	}, nil
}
