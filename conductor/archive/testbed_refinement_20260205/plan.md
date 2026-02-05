# Implementation Plan - Testbed Refinement & Housekeeping

## Phase 1: Testbed Generator Refactoring
- [x] Task: Port `testbed` to Cobra
    - [x] Initialize Cobra in `cmd/testbedgen/`.
    - [x] Migrate flags (`--config`, `--dir`) to Cobra command structure.
    - [x] Implement summarized output logic (less verbose).
- [x] Task: Update Output Layout
    - [x] Update `cmd/testbedgen/main.go` to write `cache.json` to `testbed/geminictl/cache.json`.

## Phase 2: Configuration & Data Diversity
- [x] Task: Expand Test Scenarios
    - [x] Update `default.json` to include 5 valid projects (`workdir-a` to `workdir-e`).
    - [x] Create/Rename `unlocated.json` for edge-case testing.
- [x] Task: Update Application Storage Logic
    - [x] Update `internal/cache/cache.go` to look for `geminictl/cache.json` when in testbed mode.

## Phase 3: Build & Automation
- [x] Task: Refactor Makefile
    - [x] Implement decoupled `testbed`, `run`, and `run-testbed` targets.
    - [x] Set `tmp/testbed` as the default `TESTBED_DIR`.
    - [x] Renamed build targets to `build-testbedgen` for consistency.
- [x] Task: Update .gitignore
    - [x] Ignore `tmp/` and remove old `testbed/` entries.

## Phase 4: Verification
- [x] Task: Manual Verification
    - [x] Run `make testbed` and verify `tmp/testbed` layout and 5-project content.
    - [x] Run `make run-testbed` and verify auto-selection/navigation works across more projects.
    - [x] Task: Conductor - User Manual Verification 'Housekeeping' (Protocol in workflow.md)