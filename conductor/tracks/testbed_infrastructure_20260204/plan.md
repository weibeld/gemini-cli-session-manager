# Implementation Plan - Testbed Infrastructure

## Phase 1: Storage Abstraction & Refactoring
- [ ] Task: Implement `internal/gemini` package
    - [ ] Move hashing logic to `internal/gemini`.
    - [ ] Implement directory discovery and session I/O methods.
- [ ] Task: Refactor Application to use `internal/gemini`
    - [ ] Update `internal/scanner` to use abstracted discovery logic.
    - [ ] Update `internal/cache` to use abstracted hashing.

## Phase 2: Test Data Source & Generator Logic
- [x] Task: Create Test Data Source
    - [x] Create `testdata/src/templates/` with session JSONs containing placeholders.
    - [x] Create `testdata/src/config.json` defining the test scenario (projects, sessions per project).
- [ ] Task: Implement `cmd/testgen` Tool
    - [ ] Create `cmd/testgen/main.go` using `internal/gemini`.
    - [ ] Logic: Read config -> Create project dirs -> Calc Hashes -> Gen Gemini Structure -> Inject Hashes.

## Phase 3: Application Support
- [ ] Task: Add Testbed Flag
    - [ ] Update `cmd/geminictl/root.go` to add global `--testbed <path>` flag.
    - [ ] Update `internal/scanner` and `internal/cache` constructors to accept custom root paths.
- [ ] Task: Remove Obsolete Flag
    - [ ] Remove `--reset-registry` logic and flag from `cmd/geminictl/status.go`.

## Phase 4: Integration & Verification
- [ ] Task: Update Makefile
    - [ ] Add `testrun` target: builds app, runs testgen, and launches `geminictl --testbed`.
- [ ] Task: Final Verification
    - [ ] Execute `make testrun`.
    - [ ] Verify TUI correctly identifies and resolves projects within the ephemeral `testdata/run/` environment.
    - [ ] Task: Conductor - User Manual Verification 'Testbed Integrity' (Protocol in workflow.md)
