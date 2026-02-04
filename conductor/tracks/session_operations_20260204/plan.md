# Implementation Plan - Session Operations

## Phase 1: Inspection & Integration
- [ ] Task: Implement Session Viewport
    - [ ] Create a viewport component to render session JSON content.
    - [ ] Implement text wrapping and scrolling for long transcripts.
- [ ] Task: Implement 'Resume' Action
    - [ ] Use `tea.ExecProcess` to launch `gemini --resume`.
    - [ ] Implement directory switching before execution.

## Phase 2: File Operations
- [ ] Task: Implement Delete Session
    - [ ] Add disk deletion logic for session files.
    - [ ] Handle multi-file session cleanup.
- [ ] Task: Implement Move Session
    - [ ] Add logic to relocate session files to a different project ID folder.

## Phase 3: TUI Integration
- [ ] Task: Add Keybinds
    - [ ] Map `v` or `Enter` to Inspect.
    - [ ] Map `r` to Resume.
    - [ ] Map `d` to Delete.
    - [ ] Map `m` to Move.

## Phase 4: Verification
- [ ] Task: Manual Verification
    - [ ] Inspect a session and read its messages.
    - [ ] Resume a session and verify Gemini CLI opens correctly.
    - [ ] Move/Delete sessions and verify file system changes.
    - [ ] Task: Conductor - User Manual Verification 'Session Operations' (Protocol in workflow.md)
