# Implementation Plan - Cache Optimization & UI

## Phase 1: Core Refactoring
- [x] Task: Rename and Refactor internal/registry b7e2a06
- [x] Task: Update Dependencies b7e2a06
- [x] Task: Remove Auto-Registration Logic b7e2a06

## Phase 2: Logic & Integrity Checks
- [ ] Task: Implement Runtime State Derivation
    - [ ] In `cmd/geminictl/status.go`, implement the 4-step logic:
        1. Load Cache.
        2. Scan `tmp` for active IDs.
        3. **GC:** Remove stale cache entries.
        4. **Derive:** Map Active IDs to States (Valid, Orphaned, Unlocated, New).
- [ ] Task: Trigger Scanning for 'New' Projects
    - [ ] Update `status.go` to only trigger scan for IDs not in cache.
    - [ ] Update `tui` to handle the initial derived states.

## Phase 3: UI & Verification
- [ ] Task: UI Updates
    - [ ] Display truncated hash next to path.
    - [ ] Style 'Orphaned' (strikethrough/red) and 'Unlocated' (red).
- [ ] Task: Manual Verification
    - [ ] Verify `cache.json` structure is now simple key-value.
    - [ ] Verify 'Unlocated' persists (empty string in JSON).
    - [ ] Verify 'Orphaned' is detected dynamically (path in JSON, missing on disk).
