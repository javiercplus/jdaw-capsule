// Package asset holds the resources compiled into the binary so the
// application runs without installing anything extra on the host system.
package asset

import (
	"embed"
	"fmt"
	"io/fs"
)

// DefaultFontFamily is the family name of the bundled font. It is the fallback
// used whenever the configuration does not request a different font.
const DefaultFontFamily = "CozetteVector"

//go:embed CozetteVector.ttf
var cozetteVector []byte

// CozetteVector returns the raw bytes of the bundled CozetteVector font.
func CozetteVector() []byte { return cozetteVector }

//go:embed icons
var icons embed.FS

// IconSizes lists every raster size of the app icon, smallest first.
var IconSizes = []int{16, 24, 32, 48, 64, 128, 256, 512}

// IconPNG returns the app icon rendered at the requested size.
func IconPNG(size int) ([]byte, bool) {
	data, err := icons.ReadFile(fmt.Sprintf("icons/jdaw-%d.png", size))
	if err != nil {
		return nil, false
	}
	return data, true
}

// IconICO returns the multi-resolution app icon, ready for a .desktop entry.
func IconICO() ([]byte, bool) {
	data, err := icons.ReadFile("icons/jdaw.ico")
	if err != nil {
		return nil, false
	}
	return data, true
}

// IconFS exposes the icon directory for packaging tools that walk it.
func IconFS() fs.FS {
	sub, err := fs.Sub(icons, "icons")
	if err != nil {
		return icons
	}
	return sub
}
