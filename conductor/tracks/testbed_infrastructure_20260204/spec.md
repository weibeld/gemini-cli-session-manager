# Specification: Testbed Infrastructure for Gemini CLI Session Manager

## Overview
Establish a reproducible testbed that generates a realistic environment for `geminictl`. This involves a test data generator tool that creates ephemeral Gemini CLI data structures based on the current system path, ensuring that \"Project ID -> Directory Path\" hashing is mathematically valid for verification and isolated from real user data.

## Core Features

### 1. Storage Abstraction Layer (`internal/gemini`)
- **Goal:** Centralise all knowledge of Gemini CLI's internal data storage (directory structure, hashing, file formats) to ensure consistency and facilitate future updates.
- **Functionality**:
    - Project ID calculation (SHA-256 of absolute path).
    - Directory discovery logic (walking the storage root).
    - Session file parsing and writing (JSON marshalling).
- **Semantics**: Both the main application and the test data generator MUST use this package as their sole interface for interacting with Gemini-managed data.

### 2. Test Data Source (`testdata/src/`)
- **templates/**: Session JSON templates containing placeholders (e.g., `{{PROJECT_HASH}}`) for dynamic injection.
- **config.json**: Defines the test universe, including a list of projects and their associated sessions.

### 3. Test Data Generator (`cmd/testgen`)
- **Function**: Generates a fresh working test environment in `testdata/run/`.
- **Logic**:
    1. **Initialize**: Clear `testdata/run/`.
    2. **Create Projects**: Create dummy project directories in `testdata/run/projects/`.
    3. **Calculate Hashes**: Compute SHA-256 hashes of the absolute paths of these directories.
    4. **Generate Storage**: Create the Gemini CLI structure at `testdata/run/gemini/<HASH>/chats/` using `internal/gemini`.
    5. **Inject & Write Sessions**: Populate session files by replacing placeholders with the calculated project hashes.
- **Note**: This tool does NOT generate a `cache.json`. This forces `geminictl` to perform its full resolution and integrity logic during tests.

### 4. Isolated Mode (`--testbed`)
- **Flag**: Add a global `--testbed <path>` flag.
- **Behavior**:
    - Overrides the session discovery root to `<path>/gemini`.
    - Overrides the cache file path to `<path>/cache.json`.
- **Cleanup**: Remove the obsolete `--reset-registry` flag from the `status` command.

### 5. Build & Automation
- **Makefile Integration**:
    - `make testrun`: Builds the app, runs `testgen`, and launches `geminictl` in isolated mode.
- **Playbook**: A documented set of manual steps to verify dynamic integrity logic (renaming folders, deleting data) using the testbed.

## Technical Constraints
- **Path Standardisation**: The generator and `geminictl` must use `filepath.Abs` and consistent path normalisation (no trailing slashes) to ensure hashes match perfectly.
- **Isolation**: No modifications should be made to `~/.gemini/` or `~/.config/geminictl/` when `--testbed` is active.