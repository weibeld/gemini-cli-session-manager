# Implementation Plan - Project Operations

## Phase 1: Input Infrastructure
- [ ] Task: Implement TUI Text Input
    - [ ] Integrate `bubbles/textinput` for capturing paths.
    - [ ] Create a "Prompt" component for confirmations.

## Phase 2: Core Logic
- [ ] Task: Implement Path Re-assignment
    - [ ] Create logic to update `internal/cache` with a new path for an existing ID.
    - [ ] Implement hash validation check.
- [ ] Task: Implement Delete Logic
    - [ ] Create logic to remove an ID from `internal/cache`.

## Phase 3: TUI Integration
- [ ] Task: Add Keybinds and UI Flow
    - [ ] Map `c` to "Change Directory" flow.
    - [ ] Map `d` to "Delete Project" confirmation flow.
    - [ ] Update TUI states to reflect changes immediately.

## Phase 4: Verification
- [ ] Task: Manual Verification
    - [ ] Successfully change directory of an Orphaned project and see it become Valid.
    - [ ] Delete a project and verify it disappears from `cache.json` (until next scan discovery).
    - [ ] Task: Conductor - User Manual Verification 'Project Operations' (Protocol in workflow.md)
