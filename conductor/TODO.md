# Miscellaneous TODOs

## TODO

- Check what's the 'gemini --include-directories' option
- In the projects pane, use the leaf directory of a project as the "project name" and display it separately after the hash but before the full path. Display then the full path next to the name in maybe a dimmed way (may also be truncated as it is now)
- Remove 'status' command from 'geminictl': the current functionality will be just the main functionality of the app
- In the directory input method or in the Change Directory operation, normalise the passed directory paths (e.g. either add or remove the trailing slash)
- Additional use case for this app: find out the time and date of past messages in any sessions. The Gemini app itself doesn't show timestamps and asking Gemini about the dta of past messages rarely gives the answer one expects. So, the timestamps could be integrated simply into the Inspect screen (which already shows all the message), or we might think about a dedicated UI for showing like a timeline of the messages within a sessions (or even globally across sessions). 
- Search within a session: in addition to the global search, we definitely also need a session-specific search. The most natrual location for this would be the Inspect screen where we display the message thread of a session. There could be a search bar to search the in the messages within this session, i.e. a similar experience like in the Gemini or ChatGPT apps.
- In the Inspect view of a session, allow displaying the session file JSON
  - TBD: what to do if there are multiple files for a session?
  - Also allow opening the session file directory in Finder (or just the current OS's file explorer)
- Add operations for opening the project directory or the session directory of the currently selected project or session in the OS GUI file explorer (e.g. Finder)
  - TBD: think about an operation for copying the corresponding path (project directory or session file) to the clipboard for command line usage
- UI tweak: highlight background of currently selected line in either projects pane or sessions pane (if the cursor is actively in the projects pane, then the currently selected line should be highlighted, but if the cursor switches to the sessions pane, then only the current line in the sessions pane should be highlighted and the corresponding line in the projects pane is just marked with the green text and the cursor, exactly as it is now, but it is not highlighted)
- Add a status line at the bottom (think of Vim status line, Gemini CLI status line, etc.) that can display various information. Information for display to think about:
  - Command help: commands that are applicable to the currently selected element, e.g. "d (delete), m (move), ..." etc. 
  - Path of the currently selected element (not sure how useful this is, paths are quite long)
  - Operational messages, notification (e.g. that a session has been moved, that a session or project has been deleted, etc.)
- Correct pluralisation in UI (e.g. "message" vs "messages")
- Sessions of unlocated projects cannot be opened in Gemini CLI since we can't change to the unknown directory. Maybe we should gracefully handle this, e.g. by instantly showing a message when the user attempts to open such a session (the Inspect view, on the other hand, still works)
- UI tweak: make it clearer in which pane the cursor is, e.g. change the frame colour
- Check whether the Inspect view labels messages only as User and Gemini. It seems that other types of messages (e.g. "info") are not labelled as such

## Done

- In testbed data, put cache.json into geminictl subdirectory [testbed_refinement_20260205]
- In the Makefile, make the default testbed directory tmp/testbed and then just include tmp in .gitignore [testbed_refinement_20260205]
- Improve output of 'testbed' binary [testbed_refinement_20260205]
- Use same argument parsing component for 'testbed' as for 'geminictl' (same interface, usage message, etc.) [testbed_refinement_20260205]
- In the Makefile, decouple 'testbedrun' from generating the testbed data (or create an additional target, or option to the target) so that the app can be stopped and restarted on the same testbed data [testbed_refinement_20260205]
- Make larger set of testbed data (also omit the Unlocated entry, or have a separate config with an unlocated entry) [testbed_refinement_20260205]
