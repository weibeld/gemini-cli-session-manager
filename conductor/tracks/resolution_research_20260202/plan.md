# Implementation Plan - Project ID Resolution Research

## Phase 1: Research & Discovery
- [x] Task: Audit Gemini CLI internals
    - [x] Search for undocumented config files or databases in `~/.gemini`.
    - [x] Inspect common global locations (e.g., `~/Library/Application Support/` or `~/.config/`).
- [x] Task: Review `agent-sessions` repository
    - [x] Locate the logic for Gemini CLI session discovery.
    - [x] Identify if they use a registry, a global database, or a scanning mechanism.
- [x] Task: Document Research Findings
    - [x] Create `conductor/research/resolution_strategies.md`.
    - [x] Summarize findings and propose a scanning strategy for the POC.

## Phase 2: Evaluation & POC Development [checkpoint: e6f8746]
- [x] Task: Implement Performance Benchmarks e6f8746
    - [x] Write a Go script to measure traversal speed of common directories (`~/Desktop`, `~/src`, etc.).
    - [x] Test various walk strategies (e.g., `filepath.Walk` vs `os.ReadDir`).
- [x] Task: Develop Resolution POC e6f8746
    - [x] Implement a script that takes a Project ID from `~/.gemini/tmp` and attempts to find its matching directory.
    - [x] Implement optimization: Stop scanning once the target Project ID is matched.
    - [x] Implement optimization: Allow a user-provided base directory.
- [x] Task: Final Evaluation & Report e6f8746
    - [x] Run the POC against the current environment.
    - [x] Update `conductor/research/resolution_strategies.md` with benchmark data and the final recommended implementation path.
    - [x] Task: Conductor - User Manual Verification 'Evaluation & POC Development' (Protocol in workflow.md) e6f8746
