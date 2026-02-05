# Implementation Plan - Project Operations (Deep)

## Phase 1: Core Logic (internal/gemini)
- [ ] Task: Implement Delete Logic
    - [ ] Add `gemini.DeleteProject(root, id)` to remove the directory.
- [ ] Task: Implement Move/Migrate Logic
    - [ ] Add `gemini.MoveProject(root, oldID, newPath)`:
        -   Calculate new ID.
        -   Rename directory.
        -   Update `projectHash` in all session files.

## Phase 2: TUI Integration
- [ ] Task: Update Delete Flow
    - [ ] Update `d` keybind to show stats (session count) in confirmation prompt.
    - [ ] Call `gemini.DeleteProject` on confirmation.
- [ ] Task: Update Change Directory Flow
    - [ ] Update `c` keybind to capture new path.
    - [ ] Call `gemini.MoveProject`.
    - [ ] Update Cache: Delete old, Set new.

## Phase 3: Verification
- [ ] Task: Manual Verification (Testbed)
    - [ ] Test Delete: Verify folder is gone from `testdata/run/gemini`.
    - [ ] Test Move: Verify folder is renamed and session JSONs contain new hash.
    - [ ] Task: Conductor - User Manual Verification 'Deep Project Operations' (Protocol in workflow.md)