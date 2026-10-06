package gui

import (
	"fmt"
	"jdaw-capsule/backend"
	"os"
	"path/filepath"

	"github.com/mappu/miqt/qt"
)

func CreateTabStart(env *backend.Environment, statusBar *qt.QStatusBar) *qt.QWidget {
	widget := qt.NewQWidget(nil)
	layout := qt.NewQVBoxLayout(widget)
	layout.SetContentsMargins(24, 24, 24, 24)
	layout.SetSpacing(8)

	logoLabel := qt.NewQLabel(widget)
	logoLabel.SetObjectName("AppLogo")
	logoLabel.SetAlignment(qt.AlignCenter)
	logoLabel.SetPixmap(AppLogoPixmap(logoLabel.QWidget, AppLogoSize))

	titleLabel := qt.NewQLabel(widget)
	titleLabel.SetText("JDAW Capsule")
	titleLabel.SetAlignment(qt.AlignCenter)
	titleFont := qt.NewQFont()
	titleFont.SetPointSize(22)
	titleFont.SetBold(true)
	titleLabel.SetFont(titleFont)

	subtitleLabel := qt.NewQLabel(widget)
	subtitleLabel.SetText("Portable DAW Environment for Linux")
	subtitleLabel.SetAlignment(qt.AlignCenter)
	subtitleFont := qt.NewQFont()
	subtitleFont.SetPointSize(11)
	subtitleLabel.SetFont(subtitleFont)

	layout.AddWidget(logoLabel.QWidget)
	layout.AddWidget(titleLabel.QWidget)
	layout.AddWidget(subtitleLabel.QWidget)

	layout.AddSpacing(8)
	separator := qt.NewQFrame(widget)
	separator.SetObjectName("separator")
	separator.SetFrameShape(qt.QFrame__HLine)
	separator.SetFrameShadow(qt.QFrame__Sunken)
	layout.AddWidget(separator.QWidget)
	layout.AddSpacing(8)

	infoGroup := qt.NewQGroupBox(widget)
	infoGroup.SetTitle("Environment")
	infoLayout := qt.NewQFormLayout(infoGroup.QWidget)
	infoLayout.SetContentsMargins(12, 16, 12, 12)
	infoLayout.SetSpacing(6)

	addInfoRow := func(label, value string) {
		valLabel := qt.NewQLabel3(value)
		valLabel.SetWordWrap(true)
		infoLayout.AddRow3(label, valLabel.QWidget)
	}

	addInfoRow("Wine:", env.WineBin)
	addInfoRow("REAPER:", env.ReaperDir)
	addInfoRow("Plugins:", env.Config.WindowsPluginsDir)
	addInfoRow("Audio:", env.Config.AudioBackend)
	addInfoRow("WINEPREFIX:", env.WinePrefix)

	layout.AddWidget(infoGroup.QWidget)

	layout.AddStretch()

	btnStart := qt.NewQPushButton(widget)
	btnStart.SetText("Start DAW")
	btnStart.SetMinimumHeight(48)

	startBtnFont := qt.NewQFont()
	startBtnFont.SetPointSize(13)
	startBtnFont.SetBold(true)
	btnStart.SetFont(startBtnFont)

	btnStart.OnClicked(func() {
		reaperExec := filepath.Join(env.ReaperDir, "reaper")
		if _, err := os.Stat(reaperExec); os.IsNotExist(err) {
			qt.QMessageBox_Warning(widget, "REAPER Not Installed", "REAPER is not installed.\nPlease navigate to the Install tabs and install Wine, REAPER, and Yabridge one by one.")
			return
		}

		statusBar.ShowMessage("Launching REAPER...")
		btnStart.SetEnabled(false)

		doneCh := make(chan error, 1)
		go func() {
			err := env.LaunchReaper()
			doneCh <- err
		}()

		var checkTimer *qt.QTimer
		checkTimer = qt.NewQTimer()
		checkTimer.OnTimeout(func() {
			select {
			case err := <-doneCh:
				checkTimer.Stop()
				btnStart.SetEnabled(true)
				if err != nil {
					statusBar.ShowMessage(fmt.Sprintf("Error launching REAPER: %v", err))
				} else {
					statusBar.ShowMessage2("REAPER launched successfully", 5000)
				}
			default:
				// Still launching
			}
		})
		checkTimer.Start(100)
	})

	layout.AddWidget3(btnStart.QWidget, 0, qt.AlignCenter)

	layout.AddSpacing(8)
	versionLabel := qt.NewQLabel(widget)
	hostname, _ := os.Hostname()
	versionLabel.SetText(fmt.Sprintf("Host: %s  •  Session: %s", hostname, os.Getenv("XDG_SESSION_TYPE")))
	versionLabel.SetAlignment(qt.AlignCenter)
	versionFont := qt.NewQFont()
	versionFont.SetPointSize(9)
	versionLabel.SetFont(versionFont)
	layout.AddWidget(versionLabel.QWidget)

	return widget
}
