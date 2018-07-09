PACKAGE_URL  = github.com/mexisme/osx-env-sync/osx-env-sync
PACKAGE_PATH = $(GOPATH_SRC)/$(PACKAGE_URL)

GODEP = $(GOPATH_BIN)/dep
GODEP_URL = github.com/golang/dep/cmd/dep

GOPATH = $(shell go env GOPATH)
GOPATH_BIN = $(GOPATH)/bin
GOPATH_SRC = $(GOPATH)/src

#BIN = $(GOPATH)/bin/osx-env-sync
BIN = osx-env-sync

.PHONY: build install clean
install: $(PACKAGE_PATH) build
	cd $<; cp -v $(BIN) $(TARGET_BIN)
	chmod +x $(TARGET_BIN)
build: $(PACKAGE_PATH)
	cd $<; make $(BIN)
clean:
	rm -rvf $(PACKAGE_PATH)

.PHONY: $(PACKAGE_URL)
$(PACKAGE_PATH): $(PACKAGE_URL)
	echo $@ $<; #go get -u -d $<

$(GODEP):
	go get -u $(GODEP_URL)

$(BIN): vendor
	go build .

Gopkg.lock: Gopkg.toml | $(GODEP)
	$(GODEP) ensure -update
	@touch $@
vendor: Gopkg.lock | $(GODEP)
	$(GODEP) ensure
	@touch $@

