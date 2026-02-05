# Specification: Project Operations

## Overview
Implement deep management capabilities for projects within `geminictl`. Unlike simple cache updates, these operations modify the underlying Gemini CLI storage structure to permanently delete data or migrate projects to new directory contexts.

## Core Features

### 1. Change Working Directory (Move Project)
- **Goal:** Re-assign an existing project history to a completely new file system location.
- **Trigger:** Interactive keybind (`c`).
- **Workflow:**
    1.  User selects a project.
    2.  User inputs a new absolute path.
    3.  **Calculation:** Compute `NewID` = SHA-256(NewPath).
    4.  **Validation:** Ensure `NewID` doesn't already exist in storage (collision check).
    5.  **Migration:**
        -   Rename directory: `~/.gemini/tmp/<OldID>` -> `~/.gemini/tmp/<NewID>`.
        -   Update Metadata: Iterate all `session-*.json` files and replace the `projectHash` field with `NewID`.
    6.  **Persistence:** Remove `OldID` from cache; Add `NewID` -> `NewPath` to cache.

### 2. Delete Project
- **Goal:** Permanently destroy a project and all associated sessions.
- **Trigger:** Interactive keybind (`d` or `x`).
- **Safety:**
    -   Calculate total sessions and messages to be deleted.
    -   Prompt: "Permanently delete project and <N> sessions? (y/n)"
- **Action:** `rm -rf ~/.gemini/tmp/<ProjectID>`.
- **Persistence:** Remove from cache.

## Technical Constraints
- **Atomic-ish Operations:** Moving a project involves multiple file IO steps. If renaming fails, abort. If JSON updating fails, warn the user (state might be inconsistent).
- **Session Parsing:** Must parse every session file to update the hash.
