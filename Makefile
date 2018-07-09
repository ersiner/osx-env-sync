GO = go

TARGET_BIN = ~/.osx-env-sync
TARGET_PLIST = ~/Library/LaunchAgents/osx-env-sync.plist
TARGET_RELOAD_BIN = ~/bin/osx-env-sync-now

.PHONY: install uninstall load unload
install: $(TARGET_BIN) $(TARGET_RELOAD_BIN) load
uninstall: unload
	rm -v $(TARGET_BIN) $(TARGET_RELOAD_BIN) $(TARGET_PLIST)
clean:
	which go >/dev/null && make go-osx-env-sync-clean

load: $(TARGET_PLIST)
	launchctl load $<
unload: $(TARGET_PLIST)
	launchctl unload $<

$(TARGET_BIN):
	if which go >/dev/null; then make go-osx-env-sync; else make ruby-osx-sync-env; fi

.PHONY: go-osx-env-sync go-osx-env-sync-clean
go-osx-env-sync: osx-env-sync
	cd $<; make build install TARGET_BIN=$(TARGET_BIN)
go-osx-env-sync-clean: osx-env-sync
	cd $<; make clean

.PHONY: ruby-osx-sync-env
ruby-osx-sync-env: osx-env-sync.rb
	cp -v $< $(TARGET_BIN)

~/bin:
	mkdir -pv $@

~/bin/osx-env-sync-now: osx-env-sync-now
	cp -v $< $@
	chmod -v +x $@

~/Library/LaunchAgents/osx-env-sync.plist: osx-env-sync.plist
	cp -v $< $@

~/bin/osx-env-sync-now ~/Library/LaunchAgents/osx-env-sync.plist: ~/bin


