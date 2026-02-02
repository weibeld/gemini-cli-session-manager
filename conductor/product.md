# Initial Concept
A session manager for Gemini CLI. The user will provide details on how Gemini CLI manages sessions to iteratively define requirements.

# Product Definition

## Vision
`geminictl` is a CLI utility designed to provide observability and management capabilities for Gemini CLI sessions and projects. It bridges the gap between the opaque storage mechanism of Gemini CLI (`~/.gemini/tmp`) and the user's need to understand their active contexts, creating a transparent layer for session inspection.

## Core Concepts
- **Project:** A local directory where Gemini CLI is utilized. The Project ID is the SHA-256 hash of the directory's absolute path.
- **Session:** A distinct interaction context within a project. Sessions are strictly local to their parent project.
- **Storage:** All state is maintained in `~/.gemini/tmp/<Project ID>/`, containing session files (`chats/session-*.json`), logs (`logs.json`), and checkpoints (`checkpoint-*.json`).

## Initial Features (MVP)
### `geminictl status`
A command to provide a comprehensive system-wide view:
- **Project Listing:** Displays all known projects as their readable directory paths (resolving hashes).
- **Session Breakdown:** For each project, lists active sessions.
  - **Display:** Session ID (short hash).
  - **Metrics:** Message count per session.
  - **Recency:** Relative time of the last message (e.g., "2 hours ago").

## Target Audience
- **Developers:** Who need to switch contexts or resume specific past sessions without guessing IDs.
- **Power Users:** Who want to audit their usage or clean up old session data.