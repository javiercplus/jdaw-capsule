BIN      := jdaw-capsule-gui
PREFIX   ?= $(HOME)/.local
GO       ?= go
PYTHON   ?= python3
GOSOURCES := main.go $(wildcard gui/*.go backend/*.go asset/*.go)
ASSETS    := $(wildcard asset/*.ttf asset/icons/*.png asset/icons/*.ico)
SOURCES   := $(GOSOURCES) $(ASSETS)

.PHONY: all build run fmt vet tidy clean install icons help

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
	gofmt -w $(GOSOURCES)

## vet: run go vet
vet:
	$(GO) vet ./...

## tidy: sync go.mod / go.sum
tidy:
	$(GO) mod tidy

## install: copy the binary to $(PREFIX)/bin
install: build
	install -Dm755 $(BIN) $(PREFIX)/bin/$(BIN)

## icons: regenerate the app icon set from asset/jdaw.svg
icons:
	@mkdir -p asset/icons
	@for s in 16 24 32 48 64 128 256 512; do \
		inkscape asset/jdaw.svg -w $$s -h $$s -o asset/icons/jdaw-$$s.png || exit 1; \
	done
	@$(PYTHON) -c 'import itertools, struct; \
sizes = (16, 24, 32, 48, 64, 128, 256); \
data = [open("asset/icons/jdaw-%d.png" % s, "rb").read() for s in sizes]; \
offs = list(itertools.accumulate([6 + 16 * len(sizes)] + [len(d) for d in data[:-1]])); \
entries = b"".join(struct.pack("<BBBBHHII", s % 256, s % 256, 0, 0, 1, 32, len(d), o) \
	for s, d, o in zip(sizes, data, offs)); \
head = struct.pack("<HHH", 0, 1, len(sizes)); \
open("asset/icons/jdaw.ico", "wb").write(head + entries + b"".join(data))'

## clean: remove build artifacts
clean:
	rm -f $(BIN)

## help: list available targets
help:
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/^## /  /'
