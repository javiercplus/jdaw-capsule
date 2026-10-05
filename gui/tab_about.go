package gui

import (

	"github.com/mappu/miqt/qt"
)

const aboutText = `JDAW Capsule

Portable DAW environment helper for Linux.

This tool sets up a contained Wine prefix, helps install REAPER,
manages yabridge, and points the UI to a Windows VST plugins directory.

Created for convenience and portability.`

func CreateTabAbout(_ *qt.QWidget, _ *qt.QStatusBar) *qt.QWidget {
	widget := qt.NewQWidget(nil)
	layout := qt.NewQVBoxLayout(widget)
	layout.SetContentsMargins(24, 24, 24, 24)
	layout.SetSpacing(12)

	// Logo
	logo := qt.NewQLabel(widget)
	logo.SetAlignment(qt.AlignCenter)
	logo.SetPixmap(AppLogoPixmap(logo.QWidget, AppLogoSize))
	layout.AddWidget(logo.QWidget)

	title := qt.NewQLabel(widget)
	title.SetText("JDAW Capsule")
	title.SetAlignment(qt.AlignCenter)
	titleFont := qt.NewQFont()
	titleFont.SetPointSize(18)
	titleFont.SetBold(true)
	title.SetFont(titleFont)
	layout.AddWidget(title.QWidget)

	text := qt.NewQLabel(widget)
	text.SetText(aboutText)
	text.SetAlignment(qt.AlignTop | qt.AlignLeft)
	text.SetWordWrap(true)
	text.SetOpenExternalLinks(true)
	layout.AddWidget(text.QWidget)

	repoLabel := qt.NewQLabel(widget)
	repoLabel.SetText("<a href=\"https://github.com/javiercplus/jdaw-capsule\">https://github.com/javiercplus/jdaw-capsule</a>")
	repoLabel.SetAlignment(qt.AlignCenter)
	repoLabel.SetOpenExternalLinks(true)
	layout.AddWidget(repoLabel.QWidget)

	layout.AddStretch()
	return widget
}