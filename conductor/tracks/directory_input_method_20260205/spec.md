# Specification: Advanced Directory Input Method

## Overview
Develop a superior method for directory selection within `geminictl`, prioritizing speed, intuition, and error prevention. This track evaluates two primary paradigms: optimized text-based input and a visual tree explorer.

**Reference:** [Research: TUI File Tree Explorer](../../research/file_tree_explorer.md)

## Core Features (to be evaluated)

### 1. Optimized Text-Based Input
- **Auto-completion:** Tab-completion for path segments.
- **Chunk Navigation:** Ability to delete or skip entire path segments (directories) at once.
- **Validation:** Real-time feedback on path existence.

### 2. File Tree Explorer
- **Visual Tree:** Render a hierarchical directory tree (similar to `tree -d`).
- **Interaction:** Expand/collapse nodes with `l`/`h`, navigate lines with `j`/`k`.
- **Intelligent Rooting:** Automatically root the tree at a sensible location (e.g., User Home or Project Root).

## Goals
- Determine the most "self-explaining" and intuitive method for directory selection.
- Implement the selected method as a reusable TUI component.

## Acceptance Criteria
- User can select a directory without manual typing errors.
- Navigation is fast and follows consistent TUI patterns (Vim keys).
- The chosen method is integrated into the "Move Project" operation.
