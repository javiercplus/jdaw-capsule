# jdaw-capsule
<img width="960" height="640" alt="image" src="https://github.com/user-attachments/assets/6594a53b-9911-4526-9c09-dc70ff20f287" />

Portable environment manager for Linux music production. Sets up a local Wine prefix, helps install REAPER and yabridge, and manages Windows VST plugins for use with Linux DAWs.

Built with Go + miqt (Qt) and embeds its icon and font.

## Usage

```bash
make build
./jdaw-capsule-gui
```

## Requirements

- Go 1.27+
- Qt 6 (miqt dependencies)
- Wine (for Windows VSTs/yabridge)
- `xdg-open` (to open directories from the UI)

## License

BSD-3-CLAUSE
