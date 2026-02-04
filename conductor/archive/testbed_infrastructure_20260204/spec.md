# Specification: Testbed Infrastructure for Gemini CLI Session Manager

## Overview
Establish a reproducible testbed that generates a realistic environment for `geminictl`. This involves a standalone testbed tool that creates ephemeral Gemini CLI data structures based on the current system path, ensuring that \"Project ID -> Directory Path\" hashing is mathematically valid for verification and isolated from real user data.

## Core Features

### 1. Storage Abstraction Layer (`internal/gemini`)
- **Goal:** Centralise all knowledge of Gemini CLI's internal data storage (directory structure, hashing, file formats) to ensure consistency and facilitate future updates.
- **Functionality**:
    - Project ID calculation (SHA-256 of absolute path).
    - Directory discovery logic (walking the storage root).
    - Session file parsing and writing (JSON marshalling).
- **Semantics**: Both the main application and the testbed tool MUST use this package as their sole interface for interacting with Gemini-managed data.

### 2. Testbed Configuration (`cmd/testbed/config/`)
- **templates/**: Session JSON templates containing placeholders (e.g., `{{PROJECT_HASH}}`) for dynamic injection.
- **config/default.json**: Defines the test universe, using path types to simulate different resolution states:
    - **Relative/Absolute Path**: Directory is created (Valid).
    - **Empty Path**: Virtual path used for hash and directory is NOT created (Unlocated).

### 3. Testbed Tool (`bin/testbed`)
- **Function**: Generates a fresh working test environment at a specified location.
- **Arguments**: Requires mandatory `--config` (path to config JSON) and `--dir` (target directory) flags.
- **Logic**:
    1. **Initialize**: Clears and (re)creates the target directory.
    2. **Create Projects**: Creates dummy working directories in `workdirs/` for relative paths.
    3. **Calculate Hashes**: Computes SHA-256 hashes of the absolute paths.
    4. **Generate Storage**: Creates the Gemini CLI structure at `gemini/<HASH>/chats/` using `internal/gemini`.
    5. **Inject & Write Sessions**: Populates session files by replacing placeholders with the calculated project hashes.

### 4. Isolated Mode (`--testbed`)
- **Flag**: Add a global `--testbed <path>` flag.
- **Behavior**:
    - Overrides the session discovery root to `<path>/gemini`.
    - Overrides the cache file path to `<path>/cache.json`.
- **Cleanup**: Remove the obsolete `--reset-registry` flag from the `status` command.

### 5. Build & Automation
- **Makefile Integration**:
    - `make testbed`: Builds the app and builds the testbed data in the configured directory.
    - `make testbedrun`: Builds everything and launches `geminictl` in isolated mode against the testbed.

## Technical Constraints
- **Path Standardisation**: The testbed tool and `geminictl` must use `filepath.Abs` and consistent path normalisation (no trailing slashes) to ensure hashes match perfectly.
- **Isolation**: No modifications should be made to `~/.gemini/` or `~/.config/geminictl/` when `--testbed` is active.