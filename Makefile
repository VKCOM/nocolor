NOW=`date '+%Y.%m.%d %H:%M:%S'`
OS=`uname -n -m`
AFTER_COMMIT=`git rev-parse HEAD`
GOPATH_DIR=`go env GOPATH`
BIN_NAME=nocolor
PKG=github.com/vkcom/nocolor/cmd
VERSION=1.0.3

build: clear
	go build -ldflags "-X '$(PKG).Version=$(VERSION)' -X '$(PKG).BuildTime=$(NOW)' -X '$(PKG).BuildOSUname=$(OS)' -X '$(PKG).BuildCommit=$(AFTER_COMMIT)'" -o build/$(BIN_NAME)

release:
	go run ./_scripts/release.go -version="$(VERSION)" -build-time="$(NOW)" -build-uname="$(OS)" -build-commit="$(AFTER_COMMIT)"

check: lint test

test:
	@echo "running tests..."
	@go test -tags tracing -count 3 -race -v ./tests/...
	@echo "tests passed"

lint:
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH_DIR)/bin v1.39.0
	@echo "running linters..."
	@$(GOPATH_DIR)/bin/golangci-lint run ./...
	@echo "no linter errors found"

clear:
	if [ -d build ]; then rm -r build; fi

.PHONY: check
