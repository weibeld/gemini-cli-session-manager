# Implementation Plan - Testbed Infrastructure

## Phase 1: Test Data Source & Seeder Logic
- [ ] Task: Create Test Data Source
    - [ ] Create `testdata/src/templates/` with session JSONs containing placeholders.
    - [ ] Create `testdata/src/config.json` defining the test scenario (projects, sessions per project, status).
- [ ] Task: Implement Seeder Tool
    - [ ] Create `cmd/seeder/main.go`.
    - [ ] Logic: Read config -> Create project dirs -> Calc Hashes -> Gen Gemini Structure -> Inject Hashes -> Gen Cache.

## Phase 2: Application Support
- [ ] Task: Add Testbed Flag
    - [ ] Update `cmd/geminictl/root.go` to add global `--testbed <path>` flag.
    - [ ] Update `internal/scanner` and `internal/cache` constructors to accept custom root paths.
- [ ] Task: Remove Obsolete Flag
    - [ ] Remove `--reset-registry` logic and flag from `cmd/geminictl/status.go`.

## Phase 3: Integration & Verification
- [ ] Task: Update Makefile
    - [ ] Add `testrun` target: builds app, runs seeder, and launches `geminictl --testbed`.
- [ ] Task: Final Verification
    - [ ] Execute `make testrun`.
    - [ ] Verify TUI correctly identifies and resolves projects within the ephemeral `testdata/run/` environment.
    - [ ] Task: Conductor - User Manual Verification 'Testbed Integrity' (Protocol in workflow.md)
