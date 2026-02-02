# Product Guidelines

## Visual Identity & UX
- **Self-Explaining Interface:** A primary focus is on the tool being entirely self-explaining. A user should be able to grasp the full scope and usage of the tool naturally through a single entry point, without needing to consult external documentation or complex help flags.
- **Intuitive Re-entry:** The design should ensure that a user returning after several months can immediately "drop back in" and use the tool effectively, with its functionality remaining evident and logical.
- **Minimalist Aesthetic:** CLI output should be clean and uncluttered. Use plain text by default, with sparing use of colour and emojis to highlight status or hierarchy.
- **Interactivity for Productivity:** While output is minimalist, the tool should offer rich interactive features when managing data. This includes:
    - Fuzzy search for projects and sessions.
    - Tab completion for IDs and paths.
    - Interactive selection menus for browsing and managing sessions.
- **Parsability:** Ensure core commands support non-interactive modes for use in scripts (e.g., standard UNIX flags).

## Tone & Documentation
- **Minimalist & Direct:** Documentation and help messages should be concise. Use direct language and avoid conversational filler.
- **Terminology Consistency:** Strictly adhere to Gemini CLI terminology (Project, Session, Checkpoint) to avoid user confusion.

## Error Handling
- **Brief & Informative:** Errors should be concise and clearly state what went wrong.
- **Silent Failures:** Prefer brief error messages followed by non-zero exit codes. Avoid stack traces unless a `--debug` flag is provided.

## Operational Principles
- **Safety First:** Prompt for confirmation before destructive actions (like deleting sessions).
- **Efficiency:** Maximise speed of information retrieval from `~/.gemini/tmp`.
