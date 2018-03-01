#!/bin/sh

cp osx-env-sync.plist ~/Library/LaunchAgents/osx-env-sync.plist
cp osx-env-sync.sh ~/.osx-env-sync.sh
mkdir -p ~/bin
cp osx-env-sync-now ~/bin/osx-env-sync-now
chmod +x ~/bin/osx-env-sync-now
chmod +x ~/.osx-env-sync.sh
launchctl load ~/Library/LaunchAgents/osx-env-sync.plist
