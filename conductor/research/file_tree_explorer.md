# Research: TUI File Tree Explorer for Bubbletea

## Overview
Investigation into providing an intuitive, visual directory selection experience for `geminictl`.

## 1. Options Explored

### Option A: `bubbles/filepicker` (Official)
- **Status:** Evaluated and Rejected.
- **Findings:** Primarily a file selector for a single directory. Poor support for hierarchical navigation and lacks a "Tree" view. Does not support Vim-like (`jkhl`) navigation out of the box.

### Option B: External Tools (e.g., `xplr`)
- **Status:** Considered and Rejected.
- **Findings:** Extremely powerful but adds a heavy external dependency (Rust-based). Breaks the goal of a standalone, zero-dependency Go binary.

### Option C: Custom Directory Navigator (POC)
- **Status:** Implemented as POC.
- **Findings:** Handled `jkhl` navigation and directory traversal effectively. However, it was a "Spatial" navigator (only showed current folder) rather than a true "Tree" (visualizing parent/sibling context).

## 2. Requirements for a "True" Tree Explorer
A future implementation should support:
- **Hierarchical Rendering:** Indented tree lines (├─, └─) similar to the `tree` command.
- **Collapsible Nodes:** Expand/Collapse subdirectories with `l`/`h`.
- **Intelligent Rooting:** Automatically determine a sensible root (e.g., Home or 2 levels above the current directory) and expand the path leading to the current project.
- **Performance:** Non-blocking directory reading to prevent UI stutter during expansion of large folders.

## 3. Recommended Path Forward
Building a robust, reusable Tree Component is a non-trivial task. It should be developed as a standalone track: **Track [ui_tree_explorer]: Implement Visual File Tree Component**. This component could later be reused for JSON inspection or other hierarchical data views.
