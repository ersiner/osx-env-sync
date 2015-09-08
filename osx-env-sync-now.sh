#!/usr/bin/env bash
launchctl unload ~/Library/LaunchAgents/osx-env-sync.plist
launchctl load ~/Library/LaunchAgents/osx-env-sync.plist
echo "Environment variables reloaded. Now relaunch your GUI apps to make them aware."
echo "For command line apps, launch a new Terminal session."
