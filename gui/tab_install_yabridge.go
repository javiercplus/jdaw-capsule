package gui

import (
	"jdaw-capsule/backend"

	"github.com/mappu/miqt/qt"
)

func CreateTabInstallYabridge(env *backend.Environment, statusBar *qt.QStatusBar) *qt.QWidget {
	widget := qt.NewQWidget(nil)
	layout := qt.NewQVBoxLayout(widget)
	layout.SetContentsMargins(24, 24, 24, 24)
	layout.SetSpacing(12)
	// --- Info group ---
	infoGroup := qt.NewQGroupBox(widget)
	infoGroup.SetTitle("yabridge Installation")
	infoLayout := qt.NewQVBoxLayout(infoGroup.QWidget)
	infoLayout.SetContentsMargins(12, 16, 12, 12)
	infoLayout.SetSpacing(8)
	descLabel := qt.NewQLabel(infoGroup.QWidget)
	descLabel.SetText("Download and install yabridge v5.1.1.\nyabridge enables using Windows VST2, VST3, and CLAP plugins in Linux DAWs.")
	descLabel.SetWordWrap(true)
	infoLayout.AddWidget(descLabel.QWidget)
	urlLabel := qt.NewQLabel(infoGroup.QWidget)
	urlLabel.SetText("Source: github.com/robbert-vdh/yabridge/releases/download/5.1.1/yabridge-5.1.1.tar.gz")
	urlLabel.SetWordWrap(true)
	urlFont := qt.NewQFont()
	urlFont.SetPointSize(9)
	urlLabel.SetFont(urlFont)
	infoLayout.AddWidget(urlLabel.QWidget)
	// Target path info
	pathLabel := qt.NewQLabel(infoGroup.QWidget)
	pathLabel.SetText("Install to: " + env.YabridgeDir)
	pathLabel.SetWordWrap(true)
	pathFont := qt.NewQFont()
	pathFont.SetPointSize(9)
	pathLabel.SetFont(pathFont)
	infoLayout.AddWidget(pathLabel.QWidget)
	layout.AddWidget(infoGroup.QWidget)
	// --- Progress section ---
	progressBar := qt.NewQProgressBar(widget)
	progressBar.SetMinimum(0)
	progressBar.SetMaximum(0)
	progressBar.SetTextVisible(true)
	progressBar.SetFormat("Waiting...")
	progressBar.SetVisible(false)
	layout.AddWidget(progressBar.QWidget)
	// --- Status label ---
	statusLabel := qt.NewQLabel(widget)
	statusLabel.SetAlignment(qt.AlignCenter)
	statusLabel.SetWordWrap(true)
	layout.AddWidget(statusLabel.QWidget)
	// --- Install button ---
	btnInstall := qt.NewQPushButton(widget)
	btnInstall.SetText("Install yabridge")
	btnInstall.SetMinimumHeight(40)
	installFont := qt.NewQFont()
	installFont.SetPointSize(11)
	btnInstall.SetFont(installFont)
	btnInstall.OnClicked(func() {
		statusLabel.SetText("Downloading and extracting yabridge...\nThe UI may be unresponsive during this process.")
		btnInstall.SetEnabled(false)
		progressBar.SetVisible(true)
		progressBar.SetFormat("Downloading & Extracting...")
		statusBar.ShowMessage("Installing yabridge...")
		err := env.DownloadYabridge(nil)
		progressBar.SetVisible(false)
		if err != nil {
			statusLabel.SetText("❌ Error: " + err.Error())
			statusBar.ShowMessage("yabridge installation failed")
		} else {
			statusLabel.SetText("✅ yabridge installed successfully!")
			statusBar.ShowMessage2("yabridge installed successfully", 5000)
		}
		btnInstall.SetEnabled(true)
	})
	layout.AddWidget3(btnInstall.QWidget, 0, qt.AlignCenter)
	// --- Stretch ---
	layout.AddStretch()
	return widget
}
