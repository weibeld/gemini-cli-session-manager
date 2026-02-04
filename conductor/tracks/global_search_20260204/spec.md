# Specification: Global Search

## Overview
Implement a high-performance search engine within `geminictl` to locate specific projects and sessions based on keywords.

## Core Features

### 1. Unified Search Input
- **Trigger:** Interactive keybind (e.g., `/`).
- **UI:** A persistent or toggleable search bar at the bottom or top of the interface.

### 2. Search Scope
- **Project Search:** Filter the project list by directory name or hash.
- **Session Search:** Filter sessions within the selected project by session ID.
- **Full-Text Search (Deep Search):** Optionally search within the *content* of session JSON files for matching messages.

### 3. Result Highlighting
- Highlight matching substrings in the search results.

## Technical Constraints
- **Performance:** Full-text search must be non-blocking. Use worker pools or background goroutines to index/search session files.
- **Filtering:** Use standard Go `strings.Contains` or fuzzy matching for project/session names.
