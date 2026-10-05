package gui

import (
	"jdaw-capsule/backend"

	"github.com/mappu/miqt/qt"
)

func CreateTabInstallReaper(env *backend.Environment, statusBar *qt.QStatusBar) *qt.QWidget {
	widget := qt.NewQWidget(nil)
	layout := qt.NewQVBoxLayout(widget)
	layout.SetContentsMargins(24, 24, 24, 24)
	layout.SetSpacing(12)
	// --- Info group ---
	infoGroup := qt.NewQGroupBox(widget)
	infoGroup.SetTitle("REAPER Installation")
	infoLayout := qt.NewQVBoxLayout(infoGroup.QWidget)
	infoLayout.SetContentsMargins(12, 16, 12, 12)
	infoLayout.SetSpacing(8)
	descLabel := qt.NewQLabel(infoGroup.QWidget)
	descLabel.SetText("Download and install REAPER v7.82 (Linux x86_64).\nREAPER is a professional digital audio workstation.")
	descLabel.SetWordWrap(true)
	infoLayout.AddWidget(descLabel.QWidget)
	urlLabel := qt.NewQLabel(infoGroup.QWidget)
	urlLabel.SetText("Source: reaper.fm/files/7.x/reaper782_linux_x86_64.tar.xz")
	urlLabel.SetWordWrap(true)
	urlFont := qt.NewQFont()
	urlFont.SetPointSize(9)
	urlLabel.SetFont(urlFont)
	infoLayout.AddWidget(urlLabel.QWidget)
	// Target path info
	pathLabel := qt.NewQLabel(infoGroup.QWidget)
	pathLabel.SetText("Install to: " + env.ReaperDir)
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
	btnInstall.SetText("Install REAPER")
	btnInstall.SetMinimumHeight(40)
	installFont := qt.NewQFont()
	installFont.SetPointSize(11)
	btnInstall.SetFont(installFont)
	btnInstall.OnClicked(func() {
		statusLabel.SetText("Downloading and extracting REAPER...\nThe UI may be unresponsive during this process.")
		btnInstall.SetEnabled(false)
		progressBar.SetVisible(true)
		progressBar.SetFormat("Downloading & Extracting...")
		statusBar.ShowMessage("Installing REAPER...")
		err := env.DownloadReaper(nil)
		progressBar.SetVisible(false)
		if err != nil {
			statusLabel.SetText("❌ Error: " + err.Error())
			statusBar.ShowMessage("REAPER installation failed")
		} else {
			statusLabel.SetText("✅ REAPER installed successfully!")
			statusBar.ShowMessage2("REAPER installed successfully", 5000)
		}
		btnInstall.SetEnabled(true)
	})
	layout.AddWidget3(btnInstall.QWidget, 0, qt.AlignCenter)
	// --- Stretch ---
	layout.AddStretch()
	return widget
}
