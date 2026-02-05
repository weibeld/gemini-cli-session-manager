# Specification: Testbed Refinement & Housekeeping

## Overview
Perform a series of technical debt resolution and housekeeping tasks to solidify the testbed infrastructure and improve codebase consistency before tackling advanced session operations.

## Core Refinements

### 1. Unified Argument Parsing
- **Goal:** Switch the `testbed` generator tool to use the **Cobra** library for CLI argument parsing.
- **Benefit:** Ensures consistent flag handling, help message formatting, and user experience across both `geminictl` and `testbed` binaries.
- **Output:** Reduce verbosity. Provide a high-level summary of the generation (e.g., "Generated 5 projects in 1.2s") instead of logging every individual file operation.

### 2. Testbed Data Layout & Content
- **Goal:** Isolate `geminictl` artifacts and expand test coverage.
- **Data Layout:** Move `cache.json` to a subdirectory: `<testbed_dir>/geminictl/cache.json`.
- **Default Scenario (`default.json`):**
    - 5 projects, all in the **Valid** state (directories exist).
    - Working directory names: `workdir-a`, `workdir-b`, `workdir-c`, `workdir-d`, `workdir-e`.
- **Special Scenarios:**
    - `unlocated.json`: Retain the scenario with unresolved project IDs for testing resolution logic.

### 3. Makefile Decoupling & Target Renaming
- **Goal:** Separate build, seed, and execution logic.
- **Targets:**
    - `testbed`: Build `bin/testbed` and generate data in `tmp/testbed/`.
    - `run`: Build `bin/geminictl` and run against real data.
    - `run-testbed`: Build `bin/geminictl` and run against `tmp/testbed/` (using `--testbed` flag).

### 4. Temporary Directory Management
- **Goal:** Move all ephemeral test artifacts out of the project root to reduce noise.
- **Change:** Change the default testbed output location to `tmp/testbed/`.
- **Git:** Update `.gitignore` to ignore the entire `tmp/` directory.

## Success Criteria
- `bin/testbed --help` shows standard Cobra help output.
- `make testbed` generates 5 valid projects in `tmp/testbed/`.
- `make run-testbed` successfully launches the TUI against the generated data without needing a manual seed step if data already exists.
- The `testbed` tool is quiet by default, printing only a summary.
