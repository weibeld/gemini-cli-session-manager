# Specification: Optimize Cache Persistence & UI

## Overview
Enhance `geminictl` to intelligently manage the lifecycle of project paths. This includes persisting "unlocated" projects to prevent redundant scanning, validating existing cache entries to detect deleted/renamed directories, and ensuring the cache accurately reflects the state of `~/.gemini/tmp`.

## 1. Data Model (Persistence)
The system maintains a simple cache file (`cache.json`) mapping Project IDs to Directory Paths.
- **Structure:** `map[string]string`
    - Key: Project ID (Hash)
    - Value: Absolute Directory Path (or `""` if unlocated).
- **Semantics:** This file is strictly a cache. All lifecycle states are derived at runtime.

## 2. Runtime State Derivation & Integrity
On startup, the application performs the following logic to determine project states:

1.  **Load:** Load the `cache.json` cache.
2.  **Discover:** Scan `~/.gemini/tmp` to identify all currently *Active* Project IDs.
3.  **Garbage Collection (GC):**
    -   Iterate through the Cache.
    -   If a Cache ID is NOT in the Active ID list, **delete** it from the Cache (and save).
4.  **State Assignment:**
    -   For each Active ID:
        -   **NOT in Cache:** State = **New** (Trigger Background Scan).
        -   **In Cache, Value is Empty (`""`):** State = **Unlocated**.
        -   **In Cache, Value is Path:**
            -   `os.Stat(path)` exists? State = **Valid**.
            -   `os.Stat(path)` missing? State = **Orphaned**.

## 3. UI Elements
- **Projects Pane:**
    -   **Valid:** `(truncated hash) Path`
    -   **Orphaned:** `(truncated hash) Old Path [Orphaned]` (Path strikethrough).
    -   **Unlocated:** `(truncated hash) [Unlocated]` (Red).
    -   **Scanning:** `(truncated hash) ... [Scanning]`.
- **Sessions Pane:**
    -   Unchanged (ID, Count, Time).

## 4. Scan Behavior
- **New Projects:** Automatically queued for background scanning.
- **Unlocated Projects:** NOT re-scanned (cache hit prevents it).
- **Orphaned Projects:** NOT re-scanned (cache hit prevents it).


## 4. Future Operations (Context for Architecture)
*Note: These are not in scope for this track but inform the data structure.*
-   **Projects:** Change directory (manual resolve), Delete.
-   **Sessions:** Move to new project, Delete, Inspect (view messages), Open in Gemini CLI (`gemini --resume`).
-   **General:** Search across all projects/sessions.
