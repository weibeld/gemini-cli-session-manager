# Implementation Plan - Cache Optimization & UI

## Phase 1: Core Refactoring
- [ ] Task: Rename and Refactor internal/registry
    - [ ] Rename `internal/registry` to `internal/cache`.
    - [ ] Simplify data model to `map[string]string` (Hash -> Path).
    - [ ] Update persistence to load/save this simple map.
- [ ] Task: Update Dependencies
    - [ ] Update `cmd/geminictl` and `internal/tui` to use `internal/cache`.

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
