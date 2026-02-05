# Technology Stack

## Core Language & Frameworks
- **Language:** Go (Golang) - Chosen for its performance, ease of creating single binaries, and superior TUI libraries.
- **TUI Framework:** [Bubbletea](https://github.com/charmbracelet/bubbletea) - A powerful, Elm-inspired framework for building self-explaining, interactive terminal user interfaces.
- **CLI Arguments:** [Cobra](https://github.com/spf13/cobra) - For standard command-line interface structure and flag handling.

## Data & Persistence
- **Cache:** A simple JSON file located at `~/.config/geminictl/cache.json`.
    - **Purpose:** Maps SHA-256 hashes (Project IDs) back to absolute directory paths to facilitate instant loading and prevent redundant scanning.
- **Data Source:** Read-only access to `~/.gemini/tmp/` for session data, logs, and checkpoints.

## Development & Tooling
- **Build System:** Standard Go toolchain.
- **Testing:** Go's built-in `testing` package for unit and integration tests.

## Distribution
- **Target OS:** macOS and Linux (UNIX-like systems) only. Windows is not supported.
- **Primary:** Static binaries.
- **macOS:** Homebrew formula for easy installation and updates.
