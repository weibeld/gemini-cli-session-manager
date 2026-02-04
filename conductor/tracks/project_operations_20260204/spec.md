# Specification: Project Operations

## Overview
Implement essential management capabilities for projects within `geminictl`. This allows users to rectify incorrect path mappings and clean up their project landscape.

## Core Features

### 1. Change Working Directory
- **Goal:** Allow users to manually assign or update the directory path for a specific Project ID.
- **Trigger:** An interactive keybind (e.g., `c`) in the Projects pane.
- **Workflow:**
    1.  User selects a project (Valid, Unlocated, or Orphaned).
    2.  User triggers "Change Directory".
    3.  User provides a new absolute path (input prompt or fuzzy search).
    4.  The system validates the path and computes its SHA-256 hash.
    5.  **Validation:** If the hash of the new path does NOT match the Project ID, the system warns the user (or rejects the change).
    6.  **Persistence:** Upon success, the new path is saved to `cache.json`.

### 2. Delete Project
- **Goal:** Remove a project entry from the local `cache.json`.
- **Trigger:** An interactive keybind (e.g., `d` or `x`).
- **Safety:** Must prompt for confirmation.
- **Scope:** This only deletes the *cache mapping*. It does NOT delete the Gemini CLI data in `~/.gemini/tmp` or the actual project directory.

## Technical Constraints
- **State Transition:** Changing a directory for an `Orphaned` or `Unlocated` project should move it to the `Valid` state if the path exists.
- **UI Responsiveness:** Prompts must be handled gracefully within the Bubbletea TUI.
