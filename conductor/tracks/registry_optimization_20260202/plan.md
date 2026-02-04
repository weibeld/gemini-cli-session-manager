# Implementation Plan - Cache Optimization & UI

## Phase 1: Core Refactoring
- [x] Task: Rename and Refactor internal/registry b7e2a06
- [x] Task: Update Dependencies b7e2a06
- [x] Task: Remove Auto-Registration Logic b7e2a06

## Phase 2: Logic & Integrity Checks
- [x] Task: Implement Runtime State Derivation
    - [x] In `cmd/geminictl/status.go`, implement the 4-step logic:
        1. Load Cache.
        2. Scan `tmp` for active IDs.
        3. **GC:** Remove stale cache entries.
        4. **Derive:** Map Active IDs to States (Valid, Orphaned, Unlocated, New).
- [x] Task: Trigger Scanning for 'New' Projects
    - [x] Update `status.go` to only trigger scan for IDs not in cache.
    - [x] Update `tui` to handle the initial derived states.

## Phase 3: UI & Verification
- [x] Task: UI Updates
    - [x] Display truncated hash next to path.
    - [x] Style 'Orphaned' (strikethrough/red) and 'Unlocated' (red).
- [x] Task: Implement Auto-Selection
    - [x] Automatically update project selection when cursor moves.
- [x] Task: Manual Verification
    - [x] Verify `cache.json` structure is now simple key-value.
    - [x] Verify 'Unlocated' persists (empty string in JSON).
    - [x] Verify 'Orphaned' is detected dynamically (path in JSON, missing on disk).

## Phase 4: Final Polish
- [x] Task: Refine TUI Aesthetics
    - [x] Update brackets `()` -> `[]`.
    - [x] Reset colors for orphaned/unlocated tags to normal.
    - [x] Format session rows (remove separators and prefixes).
- [x] Task: Dual Pane Navigation
    - [x] Add state for pane focus.
    - [x] Implement `H` and `L` for focus switching.
    - [x] Implement session list cursor.
- [x] Task: Layout Stabilization
    - [x] Ensure consistent pane heights.

