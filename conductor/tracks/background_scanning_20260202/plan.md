# Implementation Plan - Background Project Scanning

## Phase 1: Core Logic & Concurrency
- [~] Task: Refactor Scanner for Streaming
    - [ ] Modify `internal/scanner` to support a "Resolve" mode that emits matches over a channel.
    - [ ] Implement the 4-tier directory walking logic from the POC.
    - [ ] Add logic to skip already resolved projects.
- [ ] Task: Integrate Background Scanning in Bubbletea
    - [ ] Define a `ResolutionMsg` to update the Model when a project is found.
    - [ ] Implement a Bubbletea command to start the background Goroutine on `Init`.
    - [ ] Implement a `ScanFinishedMsg` to stop all indicators.

## Phase 2: UI Implementation
- [ ] Task: Add Scanning Indicators
    - [ ] Update the `projectView` struct to track `isScanning` state.
    - [ ] Add a simple spinner or "..." next to unresolved hashes in the Sidebar.
- [ ] Task: Real-time Sidebar Updates
    - [ ] Update the `View` method to swap hashes for paths dynamically as `ResolutionMsg` arrives.
    - [ ] Ensure the "Selected" project remains stable even as the list content changes.
- [ ] Task: Styling & Orphan Indicators
    - [ ] Implement dimming/styling for projects that finish Tier 4 without resolution.

## Phase 3: Persistence & Verification
- [ ] Task: Implement Auto-Save
    - [ ] Hook into the `ResolutionMsg` handler to update and save the `Registry`.
- [ ] Task: Final Verification
    - [ ] Verify background scanning doesn't block UI input.
    - [ ] Verify `projects.json` is correctly updated.
    - [ ] Task: Conductor - User Manual Verification 'UI Responsiveness & Resolution' (Protocol in workflow.md)
