# Implementation Plan - Testbed Infrastructure

## Phase 1: Storage Abstraction & Refactoring
- [x] Task: Implement `internal/gemini` package b7e2a06
    - [x] Move hashing logic to `internal/gemini`.
    - [x] Implement directory discovery and session I/O methods.
- [x] Task: Refactor Application to use `internal/gemini` b7e2a06
    - [x] Update `internal/scanner` to use abstracted discovery logic.
    - [x] Update `internal/cache` to use abstracted hashing.

## Phase 2: Testbed Tool & Data Logic
- [x] Task: Consolidate Testbed Source b7e2a06
    - [x] Move config and templates into `cmd/testbed/`.
    - [x] Reorganize configuration into `cmd/testbed/config/default.json`.
- [x] Task: Implement `cmd/testbed` Tool b7e2a06
    - [x] Create `cmd/testbed/main.go` using `internal/gemini`.
    - [x] Implement mandatory flags: `--config` and `--dir`.
    - [x] Logic: Read config -> Create project dirs -> Calc Hashes -> Gen Gemini Structure -> Inject Hashes.
    - [x] Handle empty paths to simulate unlocated projects.

## Phase 3: Application Support
- [x] Task: Add Testbed Flag b7e2a06
    - [x] Update `cmd/geminictl/root.go` to add global `--testbed <path>` flag.
    - [x] Update `internal/scanner` and `internal/cache` constructors to accept custom root paths.
- [x] Task: Remove Obsolete Flag b7e2a06
    - [x] Remove `--reset-registry` logic and flag from `cmd/geminictl/status.go`.

## Phase 4: Integration & Verification
- [x] Task: Update Makefile b7e2a06
    - [x] Add `testbed` and `testbedrun` targets with dynamic path support.
- [x] Task: Final Verification b7e2a06
    - [x] Execute `make testbedrun`.
    - [x] Verify TUI correctly identifies and resolves projects within the isolated testbed.
    - [x] Task: Conductor - User Manual Verification 'Testbed Integrity' (Protocol in workflow.md)

## Phase 5: Refinement & Research
- [x] Task: Housekeeping & Terminology b7e2a06
    - [x] Update `.gitignore` to exclude `testbed/`.
    - [x] Standardise all naming to 'testbed'.
- [x] Task: Enhance Gemini Data Model b7e2a06
    - [x] Update `internal/gemini` with realistic session filenames and expanded JSON struct.
- [x] Task: Research Testing Automation b7e2a06
    - [x] Investigate TUI automation and document findings in `conductor/research/testing_automation.md`.
