package main

import (
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	ignoreDirs = map[string]bool{
		"node_modules": true,
		".git":         true,
		".npm":         true,
		".cache":       true,
		".gemini":      true,
		".vscode":      true,
		".idea":        true,
		"go":           true, // ~/go/pkg can be huge
		"Library":      true, // macOS specific system dir
	}

	// Common external dev paths (outside home)
	externalCommonDirs = []string{
		"/opt",
		"/var/www",
		"/usr/local/src",
		"/srv",
	}
)

type Session struct {
	SessionID string `json:"sessionId"`
}

type ProjectData struct {
	ID string
}

func main() {
	deepScan := flag.Bool("deep", false, "Enable deep scan of root filesystem (Tier 4)")
	flag.Parse()

	// 1. Get unknown project IDs
	unknownIDs, err := getUnknownProjectIDs()
	if err != nil {
		fmt.Printf("Error getting project IDs: %v\n", err)
		return
	}
	if len(unknownIDs) == 0 {
		fmt.Println("No projects found in ~/.gemini/tmp. Nothing to resolve.")
		return
	}
	fmt.Printf("Found %d projects to resolve: %v\n", len(unknownIDs), unknownIDs)

	resolved := make(map[string]string)
	visited := make(map[string]bool)
	home, _ := os.UserHomeDir()

	// Tier 1: Desktop Scan
	fmt.Println("\n--- Tier 1: Desktop Scan ---")
	t1Start := time.Now()
	desktopPath := filepath.Join(home, "Desktop")
	if _, err := os.Stat(desktopPath); err == nil {
		scan(desktopPath, unknownIDs, resolved, visited, true)
	}
	fmt.Printf("Tier 1 took: %v\n", time.Since(t1Start))
	report(resolved, unknownIDs)
	if len(resolved) == len(unknownIDs) {
		return
	}

	// Tier 2: Home Directory Scan
	fmt.Println("\n--- Tier 2: Home Directory Scan ---")
	t2Start := time.Now()
	scan(home, unknownIDs, resolved, visited, true)
	fmt.Printf("Tier 2 took: %v\n", time.Since(t2Start))
	report(resolved, unknownIDs)
	if len(resolved) == len(unknownIDs) {
		return
	}

	// Tier 3: External Common Paths
	fmt.Println("\n--- Tier 3: External Common Paths ---")
	t3Start := time.Now()
	for _, path := range externalCommonDirs {
		if _, err := os.Stat(path); err == nil {
			scan(path, unknownIDs, resolved, visited, true)
		}
	}
	fmt.Printf("Tier 3 took: %v\n", time.Since(t3Start))
	report(resolved, unknownIDs)
	if len(resolved) == len(unknownIDs) {
		return
	}

	// Tier 4: Root Filesystem
	if *deepScan {
		fmt.Println("\n--- Tier 4: Root Filesystem ---")
		t4Start := time.Now()
		scan("/", unknownIDs, resolved, visited, true)
		fmt.Printf("Tier 4 took: %v\n", time.Since(t4Start))
		report(resolved, unknownIDs)
	}
}

func getUnknownProjectIDs() ([]string, error) {
	home, _ := os.UserHomeDir()
	tmpDir := filepath.Join(home, ".gemini", "tmp")
	entries, err := os.ReadDir(tmpDir)
	if err != nil {
		return nil, err
	}

	var ids []string
	for _, entry := range entries {
		if entry.IsDir() && len(entry.Name()) == 64 { // SHA-256 length
			ids = append(ids, entry.Name())
		}
	}
	return ids, nil
}

func scan(root string, targets []string, resolved map[string]string, visited map[string]bool, recursive bool) {
	// Avoid re-scanning if we've already fully visited this root or its parent
	// Ideally we check visited inside walk, but here we just prevent top-level redundancy
	if visited[root] {
		return
	}

	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return filepath.SkipDir // Skip permission denied, etc.
		}

		// Check if we've already visited this path (e.g. if it was a Tier 1 dir inside Home)
		if visited[path] {
			return filepath.SkipDir
		}

		if d.IsDir() {
			name := d.Name()
			if ignoreDirs[name] || (strings.HasPrefix(name, ".") && name != ".") {
				visited[path] = true // Mark ignored as visited so we don't try again
				return filepath.SkipDir
			}

			// Check hash
			h := hashPath(path)
			for _, target := range targets {
				if h == target {
					fmt.Printf("MATCH FOUND: %s -> %s\n", target, path)
					resolved[target] = path
				}
			}

			// Optimization: If all resolved, stop walking?
			// Not implementing early exit from WalkDir easily without custom error,
			// but we check resolved count at Tier level.
		}

			visited[path] = true
			return nil
	})
}

func hashPath(path string) string {
	// Standardize path: absolute, no trailing slash
	abs, err := filepath.Abs(path)
	if err != nil {
		return ""
	}
	hash := sha256.Sum256([]byte(abs))
	return hex.EncodeToString(hash[:])
}

func report(resolved map[string]string, all []string) {
	fmt.Printf("Resolved %d/%d projects.\n", len(resolved), len(all))
}
