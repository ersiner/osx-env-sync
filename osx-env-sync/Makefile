PACKAGE_URL  = github.com/mexisme/osx-env-sync/osx-env-sync

GODEP = $(GOPATH_BIN)/dep
GODEP_URL = github.com/golang/dep/cmd/dep

GOPATH = $(shell go env GOPATH)
GOPATH_BIN = $(GOPATH)/bin
GOPATH_SRC = $(GOPATH)/src

BIN = osx-env-sync

.PHONY: build install clean
install: build
	cp -v $(BIN) $(TARGET_BIN)
	chmod +x $(TARGET_BIN)
build: $(BIN)
clean:
	go clean -i -r -cache -modcache

$(BIN):
	go build .
