package gui

import (
	"sync"

	"jdaw-capsule/asset"

	"github.com/mappu/miqt/qt"
)

var (
	embeddedOnce  sync.Once
	embeddedFonts []string
)

// RegisterEmbeddedFonts hands the fonts compiled into the binary to Qt so they
// become selectable the first time the application runs, with nothing installed
// on the host system. It runs once and reports the family names it made
// available.
func RegisterEmbeddedFonts() []string {
	embeddedOnce.Do(func() {
		if id := qt.QFontDatabase_AddApplicationFontFromData(asset.CozetteVector()); id >= 0 {
			embeddedFonts = append(embeddedFonts, qt.QFontDatabase_ApplicationFontFamilies(id)...)
		}
	})
	return embeddedFonts
}

// ResolveFontFamily picks the family to render with: the configuration always
// wins, and the bundled font is used only when it asks for nothing.
func ResolveFontFamily(configured string) string {
	if configured != "" {
		return configured
	}
	if families := RegisterEmbeddedFonts(); len(families) > 0 {
		return families[0]
	}
	return asset.DefaultFontFamily
}
