# Implementation Plan - Build `geminictl status`

## Phase 1: Scaffolding & Core Logic
- [~] Task: Initialize Go project
    - [ ] Run `go mod init`
    - [ ] Set up project structure (`cmd/`, `internal/`, `pkg/`)
    - [ ] Install Cobra: `go get -u github.com/spf13/cobra`
    - [ ] Install Bubbletea: `go get -u github.com/charmbracelet/bubbletea`
- [ ] Task: Implement Project Registry
    - [ ] Define `Registry` struct and JSON schema.
    - [ ] Implement `LoadRegistry` and `SaveRegistry` functions.
    - [ ] Implement logic to check if a project path still exists (orphan detection).
- [ ] Task: Implement Session Scanning
    - [ ] Define `Session` and `Project` domain models.
    - [ ] Implement `ScanGeminiTmp` function to walk `~/.gemini/tmp`.
    - [ ] Implement parsing of `session-*.json` to extract metadata (ID, count, time).
    - [ ] Correlate scanned hashes with the Registry to resolve names.

## Phase 2: CLI & TUI Implementation
- [ ] Task: Setup `geminictl status` command
    - [ ] Use Cobra to create the `status` command entry point.
- [ ] Task: Build TUI Model (Bubbletea)
    - [ ] Define the Bubbletea `Model` struct.
    - [ ] Implement `Init`, `Update`, and `View` methods.
    - [ ] Implement Project list view.
    - [ ] Implement Session list view (nested or side-by-side).
    - [ ] Add styling for "Orphaned" projects (red/dimmed).
- [ ] Task: Connect Data to TUI
    - [ ] Load Registry and Scan data on startup.
    - [ ] Pass data into the Bubbletea model.

## Phase 3: Polish & Verification
- [ ] Task: Handle Edge Cases
    - [ ] Verify behavior when `~/.gemini/tmp` is empty.
    - [ ] Verify behavior when `projects.json` is missing.
    - [ ] Verify behavior with corrupt session files.
- [ ] Task: Manual Verification
    - [ ] Run `geminictl status` and verify output matches `spec.md`.
    - [ ] Test navigation and quitting.
