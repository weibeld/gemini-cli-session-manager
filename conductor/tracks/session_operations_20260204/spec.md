# Specification: Session Operations

## Overview
Implement the core lifecycle and navigation capabilities for sessions. This allows users to view conversation history and interact with the Gemini CLI directly from `geminictl`.

## Core Features

### 1. Inspect Session
- **Goal:** View the full transcript of a selected session.
- **Trigger:** Interactive keybind (e.g., `Enter` when focusing the Sessions pane).
- **Interface:** A full-screen or popup viewport showing messages, roles (User/Assistant), and content.

### 2. Resume in Gemini CLI
- **Goal:** Launch the official Gemini CLI to continue a session.
- **Mechanism:** `chdir` to project path -> `tea.ExecProcess(gemini --resume <session-id>)`.
- **UX:** `geminictl` pauses while Gemini CLI runs and resumes immediately upon exit.

### 3. Delete Session
- **Goal:** Permanently remove a session file from the disk.
- **Scope:** Deletes `session-*.json` from `~/.gemini/tmp/<ProjectID>/chats/`.
- **Safety:** Must prompt for confirmation.

### 4. Move Session
- **Goal:** Re-assign a session to a different Project ID.
- **Mechanism:** Move the physical file to a different subdirectory in `~/.gemini/tmp/`.
- **Requirement:** User selects the target project from a list.

## Technical Constraints
- **Session Resolution:** Must handle multi-file sessions correctly when deleting or moving (move/delete all related files).
- **Process Management:** Ensure TUI state is preserved when executing Gemini CLI.
