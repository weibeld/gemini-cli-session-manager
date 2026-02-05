package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"geminictl/internal/gemini"
	"github.com/spf13/cobra"
)

type Config struct {
	Projects []ProjectConfig `json:"projects"`
}

type ProjectConfig struct {
	Path     string   `json:"path"`
	Sessions []string `json:"sessions"`
}

var (
	configPath string
	testbedDir string
	verbose    bool
)

var rootCmd = &cobra.Command{
	Use:   "testbedgen",
	Short: "Gemini CLI Session Manager Testbed Generator",
	Long:  `A tool to generate realistic, isolated Gemini CLI data for development and testing.`,
	Run:   runGenerator,
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "", "Path to test configuration JSON (MANDATORY)")
	rootCmd.PersistentFlags().StringVarP(&testbedDir, "dir", "d", "", "Directory where the testbed will be generated (MANDATORY)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Print detailed generation logs")
	_ = rootCmd.MarkPersistentFlagRequired("config")
	_ = rootCmd.MarkPersistentFlagRequired("dir")
}

func runGenerator(cmd *cobra.Command, args []string) {
	// 1. Initialize
	_ = os.RemoveAll(testbedDir)
	if err := os.MkdirAll(testbedDir, 0755); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	geminiRoot := filepath.Join(testbedDir, "gemini")
	workdirsRoot := filepath.Join(testbedDir, "workdirs")
	geminictlConfigRoot := filepath.Join(testbedDir, "geminictl")
	
	_ = os.MkdirAll(geminiRoot, 0755)
	_ = os.MkdirAll(workdirsRoot, 0755)
	_ = os.MkdirAll(geminictlConfigRoot, 0755)

	// 2. Load Config
	configData, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Printf("Error reading config %s: %v\n", configPath, err)
		os.Exit(1)
	}
	var config Config
	if err := json.Unmarshal(configData, &config); err != nil {
		fmt.Printf("Error parsing config: %v\n", err)
		os.Exit(1)
	}

	// 3. Load Session Template
	templatePath := filepath.Join("cmd", "testbedgen", "templates", "session.json")
	sessionTemplateData, err := os.ReadFile(templatePath)
	if err != nil {
		fmt.Printf("Error reading template: %v\n", err)
		os.Exit(1)
	}

	now := time.Now().Format(time.RFC3339)
	cacheData := make(map[string]string)

	// 4. Process Projects
	for i, p := range config.Projects {
		var finalPath string
		var projectHash string
		var isUnlocated bool

		if p.Path == "" {
			finalPath = fmt.Sprintf("/unlocated/project-%d", i)
			isUnlocated = true
		} else {
			finalPath = filepath.Join(workdirsRoot, p.Path)
			_ = os.MkdirAll(finalPath, 0755)
			abs, _ := filepath.Abs(finalPath)
			finalPath = abs
			isUnlocated = false
		}

		projectHash, _ = gemini.HashProjectID(finalPath)

		if !isUnlocated {
			cacheData[projectHash] = finalPath
		} else {
			cacheData[projectHash] = ""
		}

		if verbose {
			status := "Valid"
			if isUnlocated {
				status = "Unlocated"
			}
			fmt.Printf("[%s] %s -> %s\n", status, p.Path, projectHash)
		}

		// 5. Create Sessions
		for _, sID := range p.Sessions {
			content := string(sessionTemplateData)
			content = strings.ReplaceAll(content, "{{PROJECT_HASH}}", projectHash)
			content = strings.ReplaceAll(content, "{{SESSION_ID}}", sID)

			var s gemini.Session
			if err := json.Unmarshal([]byte(content), &s); err != nil {
				continue
			}
			s.StartTime = now
			s.LastUpdated = now

			_ = gemini.WriteSession(geminiRoot, projectHash, s)
		}
	}

	// 6. Write Cache
	cacheJSON, _ := json.MarshalIndent(cacheData, "", "  ")
	_ = os.WriteFile(filepath.Join(geminictlConfigRoot, "cache.json"), cacheJSON, 0644)

	fmt.Printf("Created testbed in %s\n", testbedDir)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}