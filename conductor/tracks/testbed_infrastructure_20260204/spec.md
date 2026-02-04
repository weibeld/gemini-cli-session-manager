# Specification: Testbed Infrastructure for Gemini CLI Session Manager

## Overview
Establish a reproducible testbed that generates a realistic environment for `geminictl`. This involves a seeder tool that creates ephemeral test data based on the current system path, ensuring that "Project ID -> Directory Path" hashing is mathematically valid for verification and isolated from real user data.

## Core Features

### 1. Test Data Source (`testdata/src/`)
- **templates/**: Session JSON templates containing placeholders (e.g., `{{PROJECT_HASH}}`) for dynamic injection.
- **config.json**: Defines the test scenario, including a list of projects, their expected states (Valid, Unlocated, Orphaned), and associated sessions.
- **seeder.go**: (Optional location) Logic to generate the test environment.

### 2. Seeder Tool (`cmd/seeder`)
- **Function**: Generates a fresh working test environment in `testdata/run/`.
- **Logic**:
    1. **Initialize**: Clear `testdata/run/`.
    2. **Create Projects**: Create dummy project directories in `testdata/run/projects/`.
    3. **Calculate Hashes**: Compute SHA-256 hashes of the absolute paths of these directories.
    4. **Generate Storage**: Create the Gemini CLI structure at `testdata/run/gemini/<HASH>/chats/`.
    5. **Inject & Write Sessions**: Populate session files by replacing placeholders with the calculated project hashes.
    6. **Generate Cache**: Create `testdata/run/cache.json` with the calculated hashes and paths, including deliberate inconsistencies (Orphaned/Unlocated).

### 3. Isolated Mode (`--testbed`)
- **Flag**: Add a global `--testbed <path>` flag.
- **Behavior**:
    - Overrides the session discovery root to `<path>/gemini`.
    - Overrides the cache file path to `<path>/cache.json`.
- **Cleanup**: Remove the obsolete `--reset-registry` flag from the `status` command.

### 4. Build & Automation
- **Makefile Integration**:
    - `make testrun`: Builds the app, runs the seeder, and launches `geminictl` in isolated mode.
    - `make test-seed`: Manually triggers the seeder to refresh `testdata/run/`.

## Technical Constraints
- **Path Standardisation**: The seeder and `geminictl` must use `filepath.Abs` and consistent path normalisation (no trailing slashes) to ensure hashes match perfectly.
- **Isolation**: No modifications should be made to `~/.gemini/` or `~/.config/geminictl/` when `--testbed` is active.
