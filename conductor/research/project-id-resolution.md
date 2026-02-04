# Project ID Resolution Research

## Audit of Gemini CLI Internals
- **Storage Path:** `~/.gemini/tmp/<hash>`
- **Hashing Rule:** SHA-256 of the absolute directory path (standardized, no trailing slash).
- **Configuration Files:** Audit of `state.json` and `settings.json` confirms that Gemini CLI does **not** maintain a central registry of paths. It relies on the caller being in the correct directory to re-compute the hash.

## Competitive Analysis: `jazzyalex/agent-sessions`
- **Resolution Strategy:** Opportunistic lookup table (`hash -> path`).
- **Seeding:** Candidates are collected from other agent session stores (Claude, Codex) where the `cwd` is stored in clear text.
- **Limitation:** Cannot resolve projects if no other agent has been used in that directory.

## Validated Resolution Strategy
Based on POC benchmarks, the following multi-tier scanning strategy is recommended. It prioritizes speed for the most common cases while offering a comprehensive fallback.

### Tier 1: Desktop Scan (Instant)
- **Scope:** `~/Desktop`
- **Rationale:** High-probability location for current projects.
- **Benchmark:** < 10ms

### Tier 2: Home Directory Scan (Fast)
- **Scope:** `~` (User Home) recursively.
- **Exclusions:** `Library`, `node_modules`, `.git`, `.cache`, `.npm`, `.vscode`, `.gemini`.
- **Rationale:** Resolves >85% of projects in testing.
- **Benchmark:** ~0.8 seconds (finding 6/7 projects).

### Tier 3: External Common Paths (Fallback)
- **Scope:** `/opt`, `/var/www`, `/usr/local/src`, `/srv` (recursive).
- **Rationale:** Catches common Linux deployment/dev paths outside home.
- **Benchmark:** ~1.7 seconds.

### Tier 4: Deep Root Scan (Slow / Explicit)
- **Scope:** `/` (Root filesystem).
- **Rationale:** Last resort for identifying true orphans or non-standard paths.
- **Benchmark:** ~3 - 10 seconds.

## Recommendation for Implementation
- **Auto-Scan:** On startup (or via background goroutine), run **Tier 1 and Tier 2**. This provides near-instant resolution for the vast majority of users.
- **Manual Scan:** Provide a UI command (e.g., "Deep Scan") to trigger Tiers 3 and 4 if unresolved projects remain.
- **Registry:** Persist all resolved paths to `projects.json` to prevent re-scanning.

## External Resources
- **Reference Implementation:** [agent-sessions](https://github.com/jazzyalex/agent-sessions)
    - Open-source GUI for session management across various coding agents.
    - **Usage:** Consult this repository for implementation strategies, particularly for challenges like **reverse-resolving Project IDs to directory names**.
