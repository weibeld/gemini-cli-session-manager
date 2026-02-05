# Implementation Plan - Project Operations (Deep)

## Phase 1: Core Logic (internal/gemini)
- [x] Task: Implement Delete Logic b7e2a06
    - [x] Add `gemini.DeleteProject(root, id)` to remove the directory.
- [x] Task: Implement Move/Migrate Logic b7e2a06
    - [x] Add `gemini.MoveProject(root, oldID, newPath)`:
        -   Calculate new ID.
        -   Rename directory.
        -   Update `projectHash` in all session files.

## Phase 2: TUI Integration
- [x] Task: Update Delete Flow b7e2a06
    - [x] Update `d` keybind to show stats (session count) in confirmation prompt.
    - [x] Call `gemini.DeleteProject` on confirmation.
- [x] Task: Update Change Directory Flow b7e2a06
    - [x] Update `c` keybind to capture new path.
    - [x] Call `gemini.MoveProject`.
    - [x] Update Cache: Delete old, Set new.

## Phase 3: Verification
- [x] Task: Manual Verification (Testbed)
    - [x] Test Delete: Verify folder is gone from `testdata/run/gemini`.
    - [x] Test Move: Verify folder is renamed and session JSONs contain new hash.
    - [x] Task: Conductor - User Manual Verification 'Deep Project Operations' (Protocol in workflow.md)