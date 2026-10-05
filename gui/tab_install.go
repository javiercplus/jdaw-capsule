package gui

import (
	"jdaw-capsule/backend"

	"github.com/mappu/miqt/qt"
)

func CreateTabInstall(env *backend.Environment, statusBar *qt.QStatusBar) *qt.QWidget {
	widget := qt.NewQWidget(nil)
	layout := qt.NewQVBoxLayout(widget)
	layout.SetContentsMargins(24, 24, 24, 24)
	layout.SetSpacing(12)

	// --- Info group ---
	infoGroup := qt.NewQGroupBox(widget)
	infoGroup.SetTitle("Wine Installation")
	infoLayout := qt.NewQVBoxLayout(infoGroup.QWidget)
	infoLayout.SetContentsMargins(12, 16, 12, 12)
	infoLayout.SetSpacing(8)

	descLabel := qt.NewQLabel(infoGroup.QWidget)
	descLabel.SetText("Download and install Wine 7.22 Staging (amd64) for running Windows VST plugins.\nThis will download ~70 MB and extract it to the JDAW data directory.")
	descLabel.SetWordWrap(true)
	infoLayout.AddWidget(descLabel.QWidget)

	// Target path info
	pathLabel := qt.NewQLabel(infoGroup.QWidget)
	pathLabel.SetText("Install to: " + env.JDAWDir)
	pathLabel.SetWordWrap(true)
	pathFont := qt.NewQFont()
	pathFont.SetPointSize(9)
	pathLabel.SetFont(pathFont)
	infoLayout.AddWidget(pathLabel.QWidget)

	layout.AddWidget(infoGroup.QWidget)

	// --- Progress section ---
	progressBar := qt.NewQProgressBar(widget)
	progressBar.SetMinimum(0)
	progressBar.SetMaximum(0) // Indeterminate mode
	progressBar.SetTextVisible(true)
	progressBar.SetFormat("Waiting...")
	progressBar.SetVisible(false) // Hidden until install starts

	layout.AddWidget(progressBar.QWidget)

	// --- Status label ---
	statusLabel := qt.NewQLabel(widget)
	statusLabel.SetAlignment(qt.AlignCenter)
	statusLabel.SetWordWrap(true)
	layout.AddWidget(statusLabel.QWidget)

	// --- Install button ---
	btnInstall := qt.NewQPushButton(widget)
	btnInstall.SetText("Install Wine")
	btnInstall.SetMinimumHeight(40)

	installFont := qt.NewQFont()
	installFont.SetPointSize(11)
	btnInstall.SetFont(installFont)

	btnInstall.OnClicked(func() {
		statusLabel.SetText("Downloading and extracting Wine...\nThe UI may be unresponsive during this process.")
		btnInstall.SetEnabled(false)
		progressBar.SetVisible(true)
		progressBar.SetFormat("Downloading & Extracting...")
		statusBar.ShowMessage("Installing Wine...")

		err := env.DownloadWine(nil)
		progressBar.SetVisible(false)

		if err != nil {
			statusLabel.SetText("❌ Error: " + err.Error())
			statusBar.ShowMessage("Wine installation failed")
		} else {
			statusLabel.SetText("✅ Wine installed successfully!")
			statusBar.ShowMessage2("Wine installed successfully", 5000)
		}
		btnInstall.SetEnabled(true)
	})

	layout.AddWidget3(btnInstall.QWidget, 0, qt.AlignCenter)

	// --- Stretch ---
	layout.AddStretch()

	return widget
}
