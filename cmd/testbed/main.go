package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"geminictl/internal/gemini"
)

type Config struct {
	Projects []ProjectConfig `json:"projects"`
}

type ProjectConfig struct {
	Path     string   `json:"path"` // Relative/Absolute (created) or Empty (unlocated)
	Sessions []string `json:"sessions"`
}

func main() {
	configPath := flag.String("config", "", "MANDATORY: Path to test configuration JSON")
	testbedDir := flag.String("dir", "", "MANDATORY: Directory where the testbed will be generated")
	flag.Parse()

	if *configPath == "" || *testbedDir == "" {
		fmt.Println("Usage: testbed -config <file> -dir <dir>")
		flag.PrintDefaults()
		os.Exit(1)
	}

	fmt.Printf("Refreshing testbed in %s...\n", *testbedDir)
	
	// Clear and (re)create the output directory
	_ = os.RemoveAll(*testbedDir)
	if err := os.MkdirAll(*testbedDir, 0755); err != nil {
		fmt.Printf("Error creating output directory: %v\n", err)
		os.Exit(1)
	}

	geminiRoot := filepath.Join(*testbedDir, "gemini")
	workdirsRoot := filepath.Join(*testbedDir, "workdirs")
	_ = os.MkdirAll(geminiRoot, 0755)
	_ = os.MkdirAll(workdirsRoot, 0755)

	// 1. Load Config
	configData, err := os.ReadFile(*configPath)
	if err != nil {
		fmt.Printf("Error reading config %s: %v\n", *configPath, err)
		os.Exit(1)
	}
	var config Config
	if err := json.Unmarshal(configData, &config); err != nil {
		fmt.Printf("Error parsing config: %v\n", err)
		os.Exit(1)
	}

	// 2. Load Session Template
	sessionTemplateData, err := os.ReadFile("cmd/testbed/templates/session.json")
	if err != nil {
		fmt.Printf("Error reading template: %v\n", err)
		os.Exit(1)
	}

	now := time.Now().Format(time.RFC3339)

	// 3. Process Projects
	for i, p := range config.Projects {
		var finalPath string
		var projectHash string
		var isUnlocated bool

		if p.Path == "" {
			// Explicit Unlocated: use a unique virtual path and DO NOT create directory
			finalPath = fmt.Sprintf("/unlocated/project-%d", i)
			isUnlocated = true
		} else {
			// Valid Project: Create the directory (relative to workdirsRoot if relative)
			finalPath = filepath.Join(workdirsRoot, p.Path)
			_ = os.MkdirAll(finalPath, 0755)
			abs, _ := filepath.Abs(finalPath)
			finalPath = abs
			isUnlocated = false
		}

		// Compute hash based on the (virtual or real) absolute path
		projectHash, _ = gemini.HashProjectID(finalPath)

		if isUnlocated {
			fmt.Printf("Simulating [Unlocated] project: %s (%s)\n", finalPath, projectHash)
		} else {
			fmt.Printf("Creating [Valid] project: %s -> %s\n", p.Path, projectHash)
		}

		// 4. Create Sessions
		for _, sID := range p.Sessions {
			content := string(sessionTemplateData)
			content = strings.ReplaceAll(content, "{{PROJECT_HASH}}", projectHash)
			content = strings.ReplaceAll(content, "{{SESSION_ID}}", sID)

			var s gemini.Session
			if err := json.Unmarshal([]byte(content), &s); err != nil {
				fmt.Printf("Error unmarshalling template for %s: %v\n", sID, err)
				continue
			}
			
			s.StartTime = now
			s.LastUpdated = now

			if err := gemini.WriteSession(geminiRoot, projectHash, s); err != nil {
				fmt.Printf("Error writing session %s: %v\n", sID, err)
			}
		}
	}

	fmt.Println("Testbed refreshed successfully.")
}