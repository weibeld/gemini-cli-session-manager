package cache

import (
	"encoding/json"
	"fmt"
	"geminictl/internal/gemini"
	"os"
	"path/filepath"
)

// CalculateProjectID returns the SHA-256 hash of the absolute path.
func CalculateProjectID(path string) (string, error) {
	return gemini.HashProjectID(path)
}

// Cache stores the mapping of project IDs to directory paths.
type Cache struct {
	Data       map[string]string `json:"data"`
	configPath string
}

// NewCache creates a new Cache instance. If baseDir is provided, it uses it as the root
// and looks for cache.json directly in that directory.
func NewCache(baseDir string) (*Cache, error) {
	var configPath string
	if baseDir != "" {
		configPath = filepath.Join(baseDir, "cache.json")
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		configPath = filepath.Join(home, ".config", "geminictl", "cache.json")
	}

	return &Cache{
		Data:       make(map[string]string),
		configPath: configPath,
	}, nil
}

// Load reads the cache from the config file.
func (c *Cache) Load() error {
	data, err := os.ReadFile(c.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			c.Data = make(map[string]string)
			return nil
		}
		return err
	}

	return json.Unmarshal(data, &c.Data)
}

// Save writes the cache to the config file.
func (c *Cache) Save() error {
	dir := filepath.Dir(c.configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(c.Data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(c.configPath, data, 0644)
}

// Get returns the path for a given ID.
func (c *Cache) Get(id string) (string, bool) {
	path, ok := c.Data[id]
	return path, ok
}

// Set adds or updates a project in the cache.
func (c *Cache) Set(id, path string) {
	c.Data[id] = path
}

// VerifyAndSet checks if the path's hash matches the projectID and updates the cache.
func (c *Cache) VerifyAndSet(projectID, path string) error {
	hash, err := CalculateProjectID(path)
	if err != nil {
		return err
	}
	if hash != projectID {
		return fmt.Errorf("path hash mismatch: expected %s, got %s", projectID, hash)
	}
	c.Set(projectID, path)
	return c.Save()
}

// Delete removes a project from the cache.
func (c *Cache) Delete(id string) error {
	delete(c.Data, id)
	return c.Save()
}

// Clear removes all projects from the cache.
func (c *Cache) Clear() {
	c.Data = make(map[string]string)
}
