# Miscellaneous TODOs

## TODO

- Check what's the 'gemini --include-directories' option
- In testbed data, put cache.json into geminictl subdirectory
- In the Makefile, make the default testbed directory tmp/testbed and then just include tmp in .gitignore
- Improve output of 'testbed' binary
- Use same argument parsing component for 'testbed' as for 'geminictl' (same interface, usage message, etc.)
- In the Makefile, decouple 'testbedrun' from generating the testbed data (or create an additional target, or option to the target) so that the app can be stopped and restarted on the same testbed data
- Make larger set of testbed data (also omit the Unlocated entry, or have a separate config with an unlocated entry)
- In the projects pane, use the leaf directory of a project as the "project name" and display it separately after the hash but before the full path. Display then the full path next to the name in maybe a dimmed way (may also be truncated as it is now)

## Done
