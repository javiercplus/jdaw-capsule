# DEV.md

Guide for developers working on the `jdaw-capsule` project.

## Architecture

`jdaw-capsule` is a Go application with a GUI built using the `miqt` bindings for Qt (PySide-like syntax in Go). 
The goal of this application is to act as a portable environment launcher for DAWs on Linux, currently specifically targeting REAPER, and setting up `yabridge` to manage Windows VSTs.

### Directory Structure

* `backend/` - Contains the logic that interfaces with the system (environment setup, downloading packages, launching REAPER).
* `gui/` - Contains all GUI components constructed using `miqt`.
* `asset/` - Contains embedded assets like icons.
* `main.go` - Entry point that initializes the environment and GUI.

## Core Components

### `backend.Environment`
This structure resolves paths and acts as the state manager.
* **JDAWDir**: Resolves to `~/.local/share/jdaw/`. This acts as the container where Wine, REAPER, and yabridge are installed.
* **WinePrefix**: Specifically `~/.local/share/jdaw/data`.

### GUI Tabs
* `tab_start.go`: The main landing tab. Verifies if REAPER is installed and manages asynchronous launching.
* `tab_install.go`, `tab_install_reaper.go`, `tab_install_yabridge.go`: Provide the installation flows. Pre-detects existing installations on load and leverages `async_util.go` to provide popup progress dialogs without freezing the main UI thread.
* `async_util.go`: Contains `runAsyncInstall`, a helper utility that executes download/extract tasks via a background goroutine and monitors them with a `QTimer` to safely interact with the Qt event loop.

## Working with the GUI

When adding new buttons or asynchronous actions, ensure they do not block the main thread.
If you need to execute blocking operations (like downloading a file or running `tar`), use a goroutine and communicate its completion back to the main thread via a channel polled by a `qt.NewQTimer()`.

## Build Instructions

To build the application, run:
```bash
make
```

Ensure you have the required Qt dependencies installed as defined by `miqt`.
