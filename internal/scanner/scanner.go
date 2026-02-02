package scanner

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/fs"
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

// Resolution represents a found project mapping.
type Resolution struct {
	Hash string
	Path string
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
	info, err := os.Stat(s.RootDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", s.RootDir)
	}

	entries, err := os.ReadDir(s.RootDir)
	if err != nil {
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
			continue
		}
		project.Sessions = sessions
		projects = append(projects, project)
	}

	return projects, nil
}

// ResolveBackground starts a 4-tier scan to resolve project hashes to paths.
// It returns a channel that emits resolutions as they are found.
func (s *Scanner) ResolveBackground(unknownIDs []string) <-chan Resolution {
	out := make(chan Resolution)
	go func() {
		defer close(out)

		targets := make(map[string]bool)
		for _, id := range unknownIDs {
			targets[id] = true
		}

		if len(targets) == 0 {
			return
		}

		visited := make(map[string]bool)
		home, _ := os.UserHomeDir()

		// Tier 1: Desktop
		s.scanTier(filepath.Join(home, "Desktop"), targets, visited, out)
		if len(targets) == 0 {
			return
		}

		// Tier 2: Home
		s.scanTier(home, targets, visited, out)
		if len(targets) == 0 {
			return
		}

		// Tier 3: External Common
		externalCommon := []string{"/opt", "/var/www", "/usr/local/src", "/srv"}
		for _, path := range externalCommon {
			s.scanTier(path, targets, visited, out)
			if len(targets) == 0 {
				return
			}
		}

		// Tier 4: Root
		s.scanTier("/", targets, visited, out)
	}()
	return out
}

var ignoreDirs = map[string]bool{
	"node_modules": true,
	".git":         true,
	".npm":         true,
	".cache":       true,
	".gemini":      true,
	".vscode":      true,
	".idea":        true,
	"go":           true,
	"Library":      true,
}

func (s *Scanner) scanTier(root string, targets map[string]bool, visited map[string]bool, out chan<- Resolution) {
	if visited[root] {
		return
	}

	_ = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return filepath.SkipDir
		}

		if visited[path] {
			return filepath.SkipDir
		}

		if d.IsDir() {
			name := d.Name()
			if ignoreDirs[name] || (strings.HasPrefix(name, ".") && name != ".") {
				visited[path] = true
				return filepath.SkipDir
			}

			// Compute and check hash
			h := s.hashPath(path)
			if targets[h] {
				out <- Resolution{Hash: h, Path: path}
				delete(targets, h)
				if len(targets) == 0 {
					return fmt.Errorf("all resolved") // Shortcut to stop walking
				}
			}
		}

		visited[path] = true
		return nil
	})
}

func (s *Scanner) hashPath(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		return ""
	}
	hash := sha256.Sum256([]byte(abs))
	return hex.EncodeToString(hash[:])
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
	info, err := os.Stat(path)
	if err != nil {
		return Session{}, err
	}

	if info.Size() > 10*1024*1024 {
		return Session{}, fmt.Errorf("file too large: %s", path)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return Session{}, err
	}

	var raw rawSession
	if err := json.Unmarshal(data, &raw); err != nil {
		return Session{}, err
	}

	return Session{
		ID:           raw.SessionID,
		MessageCount: len(raw.Messages),
		LastUpdate:   info.ModTime(),
	}, nil
}