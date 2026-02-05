# Implementation Plan - Testbed Refinement & Housekeeping

## Phase 1: Testbed Generator Refactoring
- [~] Task: Port `testbed` to Cobra
    - [ ] Initialize Cobra in `cmd/testbed/`.
    - [ ] Migrate flags (`--config`, `--dir`) to Cobra command structure.
    - [ ] Implement summarized output logic (less verbose).
- [ ] Task: Update Output Layout
    - [ ] Update `cmd/testbed/main.go` to write `cache.json` to `testbed/geminictl/cache.json`.

## Phase 2: Configuration & Data Diversity
- [ ] Task: Expand Test Scenarios
    - [ ] Update `default.json` to include 5 valid projects (`workdir-a` to `workdir-e`).
    - [ ] Create/Rename `unlocated.json` for edge-case testing.
- [ ] Task: Update Application Storage Logic
    - [ ] Update `internal/cache/cache.go` to look for `geminictl/cache.json` when in testbed mode.

## Phase 3: Build & Automation
- [ ] Task: Refactor Makefile
    - [ ] Implement decoupled `testbed`, `run`, and `run-testbed` targets.
    - [ ] Set `tmp/testbed` as the default `TESTBED_DIR`.
- [ ] Task: Update .gitignore
    - [ ] Ignore `tmp/` and remove old `testbed/` entries.

## Phase 4: Verification
- [ ] Task: Manual Verification
    - [ ] Run `make testbed` and verify `tmp/testbed` layout and 5-project content.
    - [ ] Run `make run-testbed` and verify auto-selection/navigation works across more projects.
    - [ ] Task: Conductor - User Manual Verification 'Housekeeping' (Protocol in workflow.md)
