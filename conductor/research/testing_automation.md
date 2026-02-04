# Testing Automation Research & Findings

## Overview
This document summarizes our investigation into automating integration and UI tests for `geminictl`, a TUI-based application.

## 1. Challenges with TUI Automation
- **Interactivity:** TUI apps are inherently designed for human keyboard interaction. Capturing and asserting on the "View" (terminal output) is complex and brittle.
- **Asynchronicity:** Background operations (like project ID resolution) happen in separate Goroutines. A "one-shot" command-line execution often exits before these operations finish, making them hard to verify without explicit "wait" mechanisms or polling.
- **State Complexity:** Verifying complex interactions (e.g., "press 'c', type path, press Enter") requires a robust E2E framework.

## 2. Investigated Approaches

### Approach A: Headless JSON Mode
- **Mechanism:** Add a `--format json` flag to output the application state as a machine-readable object.
- **Pros:** Easy to parse and assert on the core data model.
- **Cons:** It captures a "snapshot" of the state. Verifying background updates or user interaction requires adding even more complexity, such as a `--wait` flag or streaming line-delimited JSON.

### Approach B: E2E UI Testing
- **Tools:** `expect`, `go-expect`, or `teatest` (for Bubbletea).
- **Pros:** True black-box testing of the actual user experience.
- **Cons:** Highly coupled to the exact layout and rendering logic. Small UI tweaks (like changing a bracket style) can break all tests.

## 3. Recommended Strategy: Deterministic Testbed & Manual Playbook
Given the speed of development and the complexity of robust TUI automation, we have adopted a **Deterministic Testbed** strategy:

1. **Test Data Generator (`cmd/testgen`):** A tool that generates a perfectly reproducible set of Gemini CLI data (`~/.gemini/tmp`) and working directories. It uses the same `internal/gemini` logic as the app to ensure hash validity.
2. **Isolated Execution (`--testbed`):** The application supports a dedicated flag to run against the mock data instead of real user data.
3. **Manual Verification:** Developers use the testbed to manually walk through a defined set of scenarios (renaming folders, deleting projects) to verify the TUI's responsiveness and integrity logic.

## 4. Conclusion on Automation
While a headless JSON mode is useful for verifying "startup logic" (like garbage collection), it doesn't adequately capture the "natural operation" of the TUI. For future sophisticated features (like session movement or inspection), manual testing against the generated testbed remains the most effective and efficient approach.
