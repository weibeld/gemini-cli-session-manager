# Implementation Plan - Directory Input Method

**Reference:** [Research: TUI File Tree Explorer](../../research/file_tree_explorer.md)

## Phase 1: Research & Prototyping
- [ ] Task: Prototype Optimized Text Input
    - [ ] Experiment with segment-based navigation and auto-completion.
- [ ] Task: Prototype File Tree Explorer
    - [ ] Build a basic hierarchical rendering logic using `os.ReadDir`.
    - [ ] Implement expansion/collapsing state.

## Phase 2: Evaluation & Selection
- [ ] Task: Compare Prototypes
    - [ ] Evaluate based on "Self-explaining" and "Intuitive Re-entry" guidelines.
- [ ] Task: Finalise Specification
    - [ ] Commit to one paradigm and update `spec.md` with final requirements.

## Phase 3: Component Implementation
- [ ] Task: Build Reusable Selector Component
    - [ ] Implement the selected method as a standalone modal component.

## Phase 4: Integration & Verification
- [ ] Task: Update Move Project Flow
    - [ ] Replace simple text input with the new selector.
- [ ] Task: Verification
    - [ ] Verify successful project migration using the new input method.
    - [ ] Task: Conductor - User Manual Verification 'Advanced Directory Input' (Protocol in workflow.md)
