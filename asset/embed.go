// Package asset holds the resources compiled into the binary so the
// application runs without installing anything extra on the host system.
package asset

import _ "embed"

// DefaultFontFamily is the family name of the bundled font. It is the fallback
// used whenever the configuration does not request a different font.
const DefaultFontFamily = "CozetteVector"

//go:embed CozetteVector.ttf
var cozetteVector []byte

// CozetteVector returns the raw bytes of the bundled CozetteVector font.
func CozetteVector() []byte { return cozetteVector }
