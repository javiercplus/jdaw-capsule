package gui

import (
	"os"
	"path/filepath"
	"jdaw-capsule/backend"

	"github.com/mappu/miqt/qt"
)

func CreateTabInstall(env *backend.Environment, statusBar *qt.QStatusBar) *qt.QWidget {
	widget := qt.NewQWidget(nil)
	layout := qt.NewQVBoxLayout(widget)
	layout.SetContentsMargins(24, 24, 24, 24)
	layout.SetSpacing(12)

	infoGroup := qt.NewQGroupBox(widget)
	infoGroup.SetTitle("Wine Installation")
	infoLayout := qt.NewQVBoxLayout(infoGroup.QWidget)
	infoLayout.SetContentsMargins(12, 16, 12, 12)
	infoLayout.SetSpacing(8)

	descLabel := qt.NewQLabel(infoGroup.QWidget)
	descLabel.SetText("Download and install Wine 7.22 Staging (amd64) for running Windows VST plugins.\nThis will download ~70 MB and extract it to the JDAW data directory.")
	descLabel.SetWordWrap(true)
	infoLayout.AddWidget(descLabel.QWidget)

	pathLabel := qt.NewQLabel(infoGroup.QWidget)
	pathLabel.SetText("Install to: " + env.JDAWDir)
	pathLabel.SetWordWrap(true)
	pathFont := qt.NewQFont()
	pathFont.SetPointSize(9)
	pathLabel.SetFont(pathFont)
	infoLayout.AddWidget(pathLabel.QWidget)

	layout.AddWidget(infoGroup.QWidget)

	statusLabel := qt.NewQLabel(widget)
	statusLabel.SetAlignment(qt.AlignCenter)
	statusLabel.SetWordWrap(true)
	layout.AddWidget(statusLabel.QWidget)

	// Pre-detect installation
	if _, err := os.Stat(filepath.Join(env.WineBin, "wine")); err == nil {
		statusLabel.SetText("Wine is already installed in this environment.")
	}

	btnInstall := qt.NewQPushButton(widget)
	btnInstall.SetText("Install Wine")
	btnInstall.SetMinimumHeight(40)

	installFont := qt.NewQFont()
	installFont.SetPointSize(11)
	btnInstall.SetFont(installFont)

	btnInstall.OnClicked(func() {
		runAsyncInstall(widget, "Wine", statusLabel, statusBar, btnInstall, func() error {
			return env.DownloadWine(nil)
		}, "Wine installed successfully!")
	})

	layout.AddWidget3(btnInstall.QWidget, 0, qt.AlignCenter)

	layout.AddStretch()

	return widget
}
