# Specification: Implement Background Project Scanning & Resolution

## Overview
Implement the validated 4-tier project resolution strategy within the `geminictl status` command. The goal is to automatically resolve project hashes to directory paths in the background while keeping the TUI responsive.

## Core Features

### 1. Automatic Background Scanning
- **Trigger:** Scan starts automatically when `geminictl status` is launched.
- **Strategy:** Execute the 4-tier strategy sequentially:
    1.  **Tier 1:** Desktop Scan
    2.  **Tier 2:** Home Directory Scan
    3.  **Tier 3:** External Common Paths
    4.  **Tier 4:** Root Filesystem
- **Responsiveness:** Scanning must occur in a separate Goroutine to ensure the TUI remains interactive (scrolling, selecting) at all times.

### 2. Real-Time UI Updates
- **Initial State:** Unresolved projects display their hash ID.
- **Resolution Event:** As soon as a hash is resolved to a path, the UI must update that specific entry from Hash -> Path.
- **Indicators:**
    - Display a small spinner/indicator next to *unresolved* projects while the scan is active.
    - Remove the indicator once the project is resolved or the scan completes.

### 3. Persistence
- **Auto-Save:** Resolved paths must be immediately saved to the `projects.json` registry so they persist across runs.

### 4. Orphan Handling
- **Completion:** If the 4-tier scan completes and a project remains unresolved:
    - Stop the spinner.
    - Mark the project visually (e.g., dim text, "Orphaned" label) to indicate it could not be found.

## Technical Constraints
- **Concurrency:** Use Go channels/Bubbletea commands to stream results from the scanner to the UI model.
- **Performance:** UI framerate must not stutter during the Tier 4 deep scan.
