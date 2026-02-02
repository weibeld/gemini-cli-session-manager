# Implementation Plan - Background Project Scanning

## Phase 1: Core Logic & Concurrency
- [x] Task: Refactor Scanner for Streaming 084c20c
    - [ ] Modify `internal/scanner` to support a "Resolve" mode that emits matches over a channel.
    - [ ] Implement the 4-tier directory walking logic from the POC.
    - [ ] Add logic to skip already resolved projects.
- [x] Task: Integrate Background Scanning in Bubbletea b7e2a06
    - [x] Define a `ResolutionMsg` to update the Model when a project is found.
    - [x] Implement a Bubbletea command to start the background Goroutine on `Init`.
    - [x] Implement a `ScanFinishedMsg` to stop all indicators.

## Phase 2: UI Implementation [checkpoint: b7e2a06]
- [x] Task: Add Scanning Indicators b7e2a06
    - [x] Update the `projectView` struct to track `isScanning` state.
    - [x] Add a simple spinner or "..." next to unresolved hashes in the Sidebar.
- [x] Task: Real-time Sidebar Updates b7e2a06
    - [x] Update the `View` method to swap hashes for paths dynamically as `ResolutionMsg` arrives.
    - [x] Ensure the "Selected" project remains stable even as the list content changes.
- [x] Task: Styling & Orphan Indicators b7e2a06
    - [x] Implement dimming/styling for projects that finish Tier 4 without resolution.

## Phase 3: Persistence & Verification [checkpoint: b7e2a06]
- [x] Task: Implement Auto-Save b7e2a06
    - [x] Hook into the `ResolutionMsg` handler to update and save the `Registry`.
- [x] Task: Final Verification b7e2a06
    - [x] Verify background scanning doesn't block UI input.
    - [x] Verify `projects.json` is correctly updated.
    - [x] Task: Conductor - User Manual Verification 'UI Responsiveness & Resolution' (Protocol in workflow.md) b7e2a06
