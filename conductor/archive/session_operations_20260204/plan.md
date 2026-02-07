# Implementation Plan - Session Operations

## Phase 1: Inspection & Integration
- [x] Task: Implement Session Viewport b7e2a06
    - [x] Create a viewport component to render session JSON content.
    - [x] Implement text wrapping and scrolling for long transcripts.
- [x] Task: Implement 'Open' Action (in Gemini CLI) 3a690af
    - [x] Use `tea.ExecProcess` to launch `gemini --resume`.
    - [x] Implement confirmation dialog before opening.
    - [x] Implement directory switching before execution.

## Phase 2: File Operations
- [x] Task: Implement Delete Session 3a690af
    - [x] Add disk deletion logic for session files.
    - [x] Handle multi-file session cleanup.
- [x] Task: Implement Move Session 3a690af
    - [x] Add logic to relocate session files to a different project ID folder.

## Phase 3: TUI Integration
- [x] Task: Add Keybinds 3a690af
    - [x] Map `Space` to Inspect.
    - [x] Map `Enter` to Open (with confirmation).
    - [x] Map `d` to Delete (with confirmation).
    - [x] Map `m` to Move.

## Phase 4: UI Refinement & Polish
- [~] Task: Project Sidebar Polish
    - [ ] Display leaf directory as project name.
    - [ ] Correct pluralization ("message" vs "messages").
- [ ] Task: Inspect View Polish
    - [ ] Display message timestamps in the thread.
- [ ] Task: CLI Simplification
    - [ ] Remove `status` command and make it the default action.

## Phase 5: Verification
- [ ] Task: Manual Verification
    - [ ] Inspect a session and read its messages.
    - [ ] Resume a session and verify Gemini CLI opens correctly.
    - [ ] Move/Delete sessions and verify file system changes.
    - [ ] Task: Conductor - User Manual Verification 'Session Operations' (Protocol in workflow.md)