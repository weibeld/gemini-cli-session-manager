# Miscellaneous TODOs

## TODO

- Check what's the 'gemini --include-directories' option
- In the projects pane, use the leaf directory of a project as the "project name" and display it separately after the hash but before the full path. Display then the full path next to the name in maybe a dimmed way (may also be truncated as it is now)
- Remove 'status' command from 'geminictl': the current functionality will be just the main functionality of the app

## Done

- In testbed data, put cache.json into geminictl subdirectory [testbed_refinement_20260205]
- In the Makefile, make the default testbed directory tmp/testbed and then just include tmp in .gitignore [testbed_refinement_20260205]
- Improve output of 'testbed' binary [testbed_refinement_20260205]
- Use same argument parsing component for 'testbed' as for 'geminictl' (same interface, usage message, etc.) [testbed_refinement_20260205]
- In the Makefile, decouple 'testbedrun' from generating the testbed data (or create an additional target, or option to the target) so that the app can be stopped and restarted on the same testbed data [testbed_refinement_20260205]
- Make larger set of testbed data (also omit the Unlocated entry, or have a separate config with an unlocated entry) [testbed_refinement_20260205]
