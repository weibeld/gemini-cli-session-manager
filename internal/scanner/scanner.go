package scanner

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"geminictl/internal/gemini"
)

// Session metadata for TUI display.
type Session struct {
	ID           string
	MessageCount int
	LastUpdate   time.Time
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

// NewScanner creates a scanner. If baseDir is provided, it uses it as the root
// and looks for storage in a 'gemini' subdirectory of that path.
func NewScanner(baseDir string) (*Scanner, error) {
	var root string
	if baseDir != "" {
		root = filepath.Join(baseDir, "gemini")
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		root = filepath.Join(home, ".gemini", "tmp")
	}

	return &Scanner{
		RootDir: root,
	}, nil
}

// Scan discovery all projects and their sessions using the gemini abstraction.
func (s *Scanner) Scan() ([]ProjectData, error) {
	ids, err := gemini.ListProjectIDs(s.RootDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var projects []ProjectData
	for _, id := range ids {
		project := ProjectData{ID: id}
		
		sessions, err := gemini.ReadSessions(s.RootDir, id)
		if err != nil {
			continue
		}

		// Aggregate multi-file sessions
		sessionMap := make(map[string]*Session)
		for _, sess := range sessions {
			lastUpdate := sess.GetLastUpdate()
			if existing, ok := sessionMap[sess.ID]; ok {
				existing.MessageCount += len(sess.Messages)
				if lastUpdate.After(existing.LastUpdate) {
					existing.LastUpdate = lastUpdate
				}
			} else {
				sessionMap[sess.ID] = &Session{
					ID:           sess.ID,
					MessageCount: len(sess.Messages),
					LastUpdate:   lastUpdate,
				}
			}
		}

		var projectSessions []Session
		for _, sess := range sessionMap {
			projectSessions = append(projectSessions, *sess)
		}

		// Sort sessions by last update descending
		sort.Slice(projectSessions, func(i, j int) bool {
			return projectSessions[i].LastUpdate.After(projectSessions[j].LastUpdate)
		})

		project.Sessions = projectSessions
		projects = append(projects, project)
	}

	return projects, nil
}

// ResolveBackground starts a 4-tier scan to resolve project hashes to paths.
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

			h, err := gemini.HashProjectID(path)
			if err != nil {
				return nil
			}
			
			if targets[h] {
				out <- Resolution{Hash: h, Path: path}
				delete(targets, h)
				if len(targets) == 0 {
					return fmt.Errorf("all resolved")
				}
			}
		}

		visited[path] = true
		return nil
	})
}
