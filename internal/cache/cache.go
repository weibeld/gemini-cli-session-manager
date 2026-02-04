package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
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

// Cache stores the mapping of project IDs to directory paths.
type Cache struct {
	Data       map[string]string `json:"data"`
	configPath string
}

// NewCache creates a new Cache instance with the default config path.
func NewCache() (*Cache, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	configPath := filepath.Join(home, ".config", "geminictl", "cache.json")
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

// Delete removes a project from the cache.
func (c *Cache) Delete(id string) {
	delete(c.Data, id)
}

// Clear removes all projects from the cache.
func (c *Cache) Clear() {
	c.Data = make(map[string]string)
}