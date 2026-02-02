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
		"Library":      true,
		"Applications": true,
		"Downloads":    true, // Often huge, better to skip for benchmark
		"Music":        true,
		"Movies":       true,
		"Pictures":     true,
		".npm":         true,
		".cache":       true,
		".gemini":      true, // Avoid circular ref
	}
)

func main() {
	root := flag.String("root", "", "Root directory to scan (default: user home)")
	flag.Parse()

	scanRoot := *root
	if scanRoot == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}
		scanRoot = home
	}

	fmt.Printf("Starting benchmark scan of: %s\n", scanRoot)
	start := time.Now()
	
dirsScanned := 0
hashesComputed := 0

	err := filepath.WalkDir(scanRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			// Permission denied or other error, skip
			return nil
		}

		if d.IsDir() {
			name := d.Name()
			if ignoreDirs[name] || strings.HasPrefix(name, ".") {
				return filepath.SkipDir
			}
			dirsScanned++
			
			// Simulate the cost of hashing (the core operation)
		hashPath(path)
		hashesComputed++
		}
		return nil
	})

	duration := time.Since(start)
	if err != nil {
		fmt.Printf("Error during scan: %v\n", err)
	}

	fmt.Printf("\n--- Benchmark Results ---\n")
	fmt.Printf("Time: %v\n", duration)
	fmt.Printf("Directories Scanned: %d\n", dirsScanned)
	fmt.Printf("Rate: %.2f dirs/sec\n", float64(dirsScanned)/duration.Seconds())
}

func hashPath(path string) string {
	hash := sha256.Sum256([]byte(path))
	return hex.EncodeToString(hash[:])
}
