# Implementation Plan - Build `geminictl status`

## Phase 1: Scaffolding & Core Logic
- [x] Task: Initialize Go project c86baab
    - [ ] Run `go mod init`
    - [ ] Set up project structure (`cmd/`, `internal/`, `pkg/`)
    - [ ] Install Cobra: `go get -u github.com/spf13/cobra`
    - [ ] Install Bubbletea: `go get -u github.com/charmbracelet/bubbletea`
- [x] Task: Implement Project Registry f448f40
    - [ ] Define `Registry` struct and JSON schema.
    - [ ] Implement `LoadRegistry` and `SaveRegistry` functions.
    - [ ] Implement logic to check if a project path still exists (orphan detection).
- [x] Task: Implement Session Scanning ef4d933
    - [ ] Define `Session` and `Project` domain models.
    - [ ] Implement `ScanGeminiTmp` function to walk `~/.gemini/tmp`.
    - [ ] Implement parsing of `session-*.json` to extract metadata (ID, count, time).
    - [ ] Correlate scanned hashes with the Registry to resolve names.

## Phase 2: CLI & TUI Implementation [checkpoint: 2ff10e6]
- [x] Task: Setup `geminictl status` command 2185954
    - [ ] Use Cobra to create the `status` command entry point.
- [x] Task: Build TUI Model (Bubbletea) f4e42d8
    - [ ] Define the Bubbletea `Model` struct.
    - [ ] Implement `Init`, `Update`, and `View` methods.
    - [ ] Implement Project list view.
    - [ ] Implement Session list view (nested or side-by-side).
    - [ ] Add styling for "Orphaned" projects (red/dimmed).
- [x] Task: Connect Data to TUI 2e062f4
    - [ ] Load Registry and Scan data on startup.
    - [ ] Pass data into the Bubbletea model.

## Phase 3: Polish & Verification
- [x] Task: Handle Edge Cases e6f8746
    - [ ] Verify behavior when `~/.gemini/tmp` is empty.
    - [ ] Verify behavior when `projects.json` is missing.
    - [ ] Verify behavior with corrupt session files.
- [~] Task: Manual Verification
    - [ ] Run `geminictl status` and verify output matches `spec.md`.
    - [ ] Test navigation and quitting.
