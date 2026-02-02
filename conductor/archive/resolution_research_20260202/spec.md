# Specification: Project ID Resolution Research & Evaluation

## Overview
The core challenge for `geminictl` is resolving one-way SHA-256 hashes (Project IDs) back to absolute directory paths. This track focuses on researching optimal resolution strategies and evaluating their feasibility through performance testing and proof-of-concept (POC) development.

## Goals
1. **Magic Bullet Discovery:** Investigate if Gemini CLI stores the hash-to-path mapping in an undocumented location (e.g., SQLite DB, global config, or specific metadata files).
2. **Competitive Analysis:** Analyze the `agent-sessions` repository to determine how they handle Project ID resolution for Gemini CLI.
3. **Scanning Performance Evaluation:** Measure the performance of brute-force and optimized scanning strategies to determine if they are viable within a ~30-second window.

## Functional Requirements
- **Research Report:** A comprehensive markdown document summarizing:
    - Findings from Gemini CLI internals audit.
    - Findings from `agent-sessions` code review.
    - Recommended resolution strategy (e.g., deep scan, user-guided base path, etc.).
- **POC Implementation:** A standalone Go script or internal package to:
    - Benchmark filesystem traversal speed.
    - Test the recommended strategy against a real `~/.gemini/tmp` environment.
    - Demonstrate Project ID matching for discovered directories.

## Acceptance Criteria
- A resolution strategy is identified and documented.
- The POC demonstrates the ability to resolve at least one known Project ID hash to its correct path.
- Benchmark data confirms the recommended scanning strategy completes in < 30 seconds for a typical user directory structure.

## Out of Scope
- Full TUI integration of the scanning feature (this will be a follow-up track).
- Automatic modification of the production `projects.json` registry (the POC should use a temporary or separate data store).
