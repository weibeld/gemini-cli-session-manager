# Implementation Plan - Global Search

## Phase 1: Search Infrastructure
- [ ] Task: Implement Search State
    - [ ] Add `Query` string to the TUI Model.
    - [ ] Implement filtering logic for the Project and Session slices.

## Phase 2: UI Implementation
- [ ] Task: Add Search Bar
    - [ ] Integrate `bubbles/textinput` for the search query.
    - [ ] Map `/` to trigger search mode.
- [ ] Task: Real-time Filtering
    - [ ] Update the `View` to only display items matching the current `Query`.

## Phase 3: Deep Search (Optional/Advanced)
- [ ] Task: Implement Session Content Search
    - [ ] Add logic to parse and grep session JSON contents in the background.

## Phase 4: Verification
- [ ] Task: Manual Verification
    - [ ] Verify projects are filtered as you type.
    - [ ] Verify sessions are filtered within a project.
    - [ ] Task: Conductor - User Manual Verification 'Global Search' (Protocol in workflow.md)
