# Specification: Build `geminictl status` Command

## Goal
Implement the `geminictl status` command to provide a comprehensive, interactive view of Gemini CLI projects and sessions. This includes setting up the underlying data registry for project path resolution and the TUI for display.

## Core Features

### 1. Project Registry
- **Mechanism:** Maintain a JSON file at `~/.config/geminictl/projects.json`.
- **Data Structure:** Map SHA-256 Project IDs to absolute directory paths.
- **Functionality:**
    - Load registry on startup.
    - Ability to add/update entries (auto-registration logic to be defined in implementation).
    - Handle missing/renamed directories (orphan detection).

### 2. Session Scanning
- **Source:** Read `~/.gemini/tmp/`.
- **Logic:**
    - Iterate through project subdirectories in `tmp`.
    - For each project, parse `chats/session-*.json` files to extract:
        - Session ID.
        - Message count.
        - Last modified timestamp.
    - Parse `logs.json` (optional for MVP, focus on session files first for speed).

### 3. TUI Display (Bubbletea)
- **Layout:**
    - **Projects List:** Sidebar or main section listing resolved project paths.
    - **Session Details:** When a project is selected, show its sessions.
- **Visuals:**
    - **Active Projects:** Standard text.
    - **Orphaned Projects:** Red/Dimmed text with "[Orphan]" label.
    - **Sessions:** Show ID (short), Message Count, Relative Time (e.g., "2h ago").
- **Interactivity:**
    - Navigate list with Up/Down arrows.
    - Quit with `q` or `Esc`.

## Technical Constraints
- **OS:** macOS/Linux.
- **Language:** Go.
- **Libraries:** Bubbletea (TUI), Cobra (CLI).
