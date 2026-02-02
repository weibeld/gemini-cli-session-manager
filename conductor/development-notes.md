# Gemini CLI Internals & Session Management

*Notes provided by user on 2026-02-02.*

## Storage & Structure
- **Root Storage:** Gemini CLI maintains projects and sessions in `~/.gemini/tmp`.
- **Project mapping:** A project corresponds to a directory on the local system where the CLI is used.
- **Project ID:** The SHA-256 hash of the absolute path of the project directory.
    - Example: `echo -n "/path/to/project" | sha256` -> `6a15...`
- **Directory Structure:** `~/.gemini/tmp/<Project ID>/`

## Sessions
- **Scope:** Sessions are local to a project. A project may have multiple sessions.
- **Creation:** Starting `gemini` creates a new session (and a new project entry if one doesn't exist).
- **Listing:** `gemini --list-sessions` lists sessions *only* for the current project directory.
- **Resuming:** `gemini --resume <session id>` resumes a specific session within the current project.

## File Artifacts
### Session Data
- **Location:** `~/.gemini/tmp/<Project ID>/chats/`
- **Naming:** `session-<timestamp>-<hash>.json`
- **Content:** Contains the entire data of a session (messages, responses).
- **ID:** The `sessionId` field inside the JSON file (e.g., `deac22c2...`).
- **Sub-agents:** Using a sub-agent may create additional session files for the same session.

### Logs
- **Location:** `~/.gemini/tmp/<Project ID>/logs.json`
- **Type:** Append-only log.
- **Scope:** Logs all interactions within the project across *all* sessions.
- **Format:** Array of objects, e.g., `[ {"sessionId": "A", "message": "Hello"}, ... ]`.

### Checkpoints
- **Location:** `~/.gemini/tmp/<Project ID>/checkpoint-<name>.json`
- **Creation:** Created via `/chat save`.
- **Scope:** Shared across all sessions within the same project.

## Edge Cases & Anomalies
### Orphaned Projects (Renamed Directories)
- If a directory corresponding to a Gemini CLI project is renamed, the link is broken.
- The Project ID (hash of the *original* path) is not updated.
- `gemini --list-projects` in the new directory will not show the previous sessions (as they are tied to the old path's hash).
- `geminictl` must identify these as "orphans" (projects where the hashed path no longer exists) to provide a complete view of the system, including inconsistencies.

## External Resources
- **Reference Implementation:** [agent-sessions](https://github.com/jazzyalex/agent-sessions)
    - Open-source GUI for session management across various coding agents.
    - **Usage:** Consult this repository for implementation strategies, particularly for challenges like **reverse-resolving Project IDs to directory names**.
