# Implementation Plan - Project ID Resolution Research

## Phase 1: Research & Discovery
- [ ] Task: Audit Gemini CLI internals
    - [ ] Search for undocumented config files or databases in `~/.gemini`.
    - [ ] Inspect common global locations (e.g., `~/Library/Application Support/` or `~/.config/`).
- [ ] Task: Review `agent-sessions` repository
    - [ ] Locate the logic for Gemini CLI session discovery.
    - [ ] Identify if they use a registry, a global database, or a scanning mechanism.
- [ ] Task: Document Research Findings
    - [ ] Create `conductor/research/resolution_strategies.md`.
    - [ ] Summarize findings and propose a scanning strategy for the POC.

## Phase 2: Evaluation & POC Development
- [ ] Task: Implement Performance Benchmarks
    - [ ] Write a Go script to measure traversal speed of common directories (`~/Desktop`, `~/src`, etc.).
    - [ ] Test various walk strategies (e.g., `filepath.Walk` vs `os.ReadDir`).
- [ ] Task: Develop Resolution POC
    - [ ] Implement a script that takes a Project ID from `~/.gemini/tmp` and attempts to find its matching directory.
    - [ ] Implement optimization: Stop scanning once the target Project ID is matched.
    - [ ] Implement optimization: Allow a user-provided base directory.
- [ ] Task: Final Evaluation & Report
    - [ ] Run the POC against the current environment.
    - [ ] Update `conductor/research/resolution_strategies.md` with benchmark data and the final recommended implementation path.
    - [ ] Task: Conductor - User Manual Verification 'Evaluation & POC Development' (Protocol in workflow.md)
