#!/bin/bash -e
## NOTE: '-e' means the script will fail at the first error, hopefully avoiding irreparable damage.

PROG_DIR="$(dirname $0)"
mkdir -pv ~/bin

if which go >/dev/null; then
  go get -v github.com/mexisme/osx-env-sync/osx-env-sync
  cp -v $(go env GOPATH)/bin/osx-env-sync ~/.osx-env-sync
else
  (
    cd "${PROG_DIR}"
    cp -v osx-env-sync.rb ~/.osx-env-sync
  )
fi

(
  cd "${PROG_DIR}"
  cp -v osx-env-sync.plist ~/Library/LaunchAgents/osx-env-sync.plist
  cp -v osx-env-sync-now ~/bin/osx-env-sync-now
)

chmod -v +x ~/bin/osx-env-sync-now
chmod -v +x ~/.osx-env-sync

launchctl load ~/Library/LaunchAgents/osx-env-sync.plist
