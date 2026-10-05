BIN      := jdaw-capsule-gui
PREFIX   ?= $(HOME)/.local
GO       ?= go
SOURCES  := main.go $(wildcard gui/*.go backend/*.go)

.PHONY: all build run fmt vet tidy clean install help

all: build

## build: compile the GUI binary
build: $(BIN)

$(BIN): $(SOURCES) go.mod go.sum
	$(GO) build -o $(BIN) .

## run: build and launch the GUI
run: build
	./$(BIN)

## fmt: format all Go sources
fmt:
	gofmt -w $(SOURCES)

## vet: run go vet
vet:
	$(GO) vet ./...

## tidy: sync go.mod / go.sum
tidy:
	$(GO) mod tidy

## install: copy the binary to $(PREFIX)/bin
install: build
	install -Dm755 $(BIN) $(PREFIX)/bin/$(BIN)

## clean: remove build artifacts
clean:
	rm -f $(BIN)

## help: list available targets
help:
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/^## /  /'
