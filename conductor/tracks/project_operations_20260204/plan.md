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
- [x] Task: Manual Verification (Testbed) b7e2a06
    - [x] Test Delete: Verify folder is gone from `testdata/run/gemini`.
    - [x] Test Move: Verify folder is renamed and session JSONs contain new hash.
    - [x] Task: Conductor - User Manual Verification 'Deep Project Operations' (Protocol in workflow.md)

## Phase 4: Unified Modals & Selectors
- [x] Task: Design Unified Modal Component b7e2a06
    - [x] Create a reusable modal frame with consistent styling.
    - [x] Implement state management for modal overlays.
- [x] Task: Implement TextInputModal
    - [x] Re-implement project migration via a modal text input for stability.
- [x] Task: Implement Project List Selector b7e2a06
    - [x] Create a modal for selecting from known projects (for future session moves).
- [x] Task: Migrating Existing Prompts b7e2a06
    - [x] Refactor the Delete Project confirmation to use the unified modal.