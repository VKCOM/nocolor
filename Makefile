NOW=`date +%Y%m%d%H%M%S`
OS=`uname -n -m`
AFTER_COMMIT=`git rev-parse HEAD`
GOPATH_DIR=`go env GOPATH`
BIN_NAME="nocolor"
PKG=github.com/vkcom/nocolor/cmd

build: clear
	go build -ldflags "-X '$(PKG).BuildTime=$(NOW)' -X '$(PKG).BuildOSUname=$(OS)' -X '$(PKG).BuildCommit=$(AFTER_COMMIT)'" -o build/$(BIN_NAME)

build-release:
	go run ./_scripts/release.go -build-time="$(NOW)" -build-uname="$(OS)" -build-commit="$(AFTER_COMMIT)"

check:
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH_DIR)/bin v1.39.0
	@echo "running linters..."
	@$(GOPATH_DIR)/bin/golangci-lint run ./...
	@echo "running tests..."
	@go test -tags tracing -count 3 -race -v ./tests/...
	@echo "everything is OK"

clear:
	if [ -d build ]; then rm -r build; fi

.PHONY: check
