package registry

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

// CalculateProjectID returns the SHA-256 hash of the absolute path.
func CalculateProjectID(path string) (string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256([]byte(abs))
	return hex.EncodeToString(hash[:]), nil
}

// Project represents a Gemini CLI project mapping.
type Project struct {
	ID   string `json:"id"`   // SHA-256 hash
	Path string `json:"path"` // Absolute directory path
}

// Registry stores the mapping of project IDs to directory paths.
type Registry struct {
	Projects []Project `json:"projects"`
	configPath string
}

// NewRegistry creates a new Registry instance with the default config path.
func NewRegistry() (*Registry, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	configPath := filepath.Join(home, ".config", "geminictl", "projects.json")
	return &Registry{
		Projects:   []Project{},
		configPath: configPath,
	}, nil
}

// Load reads the registry from the config file.
func (r *Registry) Load() error {
	data, err := os.ReadFile(r.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			r.Projects = []Project{}
			return nil
		}
		return err
	}

	return json.Unmarshal(data, r)
}

// Save writes the registry to the config file.
func (r *Registry) Save() error {
	dir := filepath.Dir(r.configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(r.configPath, data, 0644)
}

// GetProjectPath returns the path for a given ID and whether the project is an orphan.
func (r *Registry) GetProjectPath(id string) (string, bool, error) {
	for _, p := range r.Projects {
		if p.ID == id {
			_, err := os.Stat(p.Path)
			isOrphan := os.IsNotExist(err)
			return p.Path, isOrphan, nil
		}
	}
	return "", false, errors.New("project not found in registry")
}

// AddProject adds or updates a project in the registry.
func (r *Registry) AddProject(id, path string) {
	for i, p := range r.Projects {
		if p.ID == id {
			r.Projects[i].Path = path
			return
		}
	}
	r.Projects = append(r.Projects, Project{ID: id, Path: path})
}
