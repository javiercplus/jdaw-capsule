package gui

import (
	"github.com/mappu/miqt/qt"
)

func runAsyncInstall(parent *qt.QWidget, title string, statusLabel *qt.QLabel, statusBar *qt.QStatusBar, btnInstall *qt.QPushButton, installFunc func() error, successMsg string) {
	statusLabel.SetText("Downloading and extracting...\nThe UI may be unresponsive during this process.")
	btnInstall.SetEnabled(false)
	statusBar.ShowMessage("Installing " + title + "...")

	dialog := qt.NewQProgressDialog(parent)
	dialog.SetLabelText("Downloading and extracting...")
	dialog.SetMinimum(0)
	dialog.SetMaximum(0)
	dialog.SetWindowTitle(title)
	dialog.SetCancelButton(nil)
	dialog.Show()

	doneCh := make(chan error, 1)

	go func() {
		err := installFunc()
		doneCh <- err
	}()

	var checkTimer *qt.QTimer
	checkTimer = qt.NewQTimer()
	checkTimer.OnTimeout(func() {
		select {
		case err := <-doneCh:
			checkTimer.Stop()
			dialog.Close()
			if err != nil {
				statusLabel.SetText("❌ Error: " + err.Error())
				statusBar.ShowMessage(title + " failed")
			} else {
				statusLabel.SetText("✅ " + successMsg)
				statusBar.ShowMessage2(successMsg, 5000)
			}
			btnInstall.SetEnabled(true)
		default:
			// Still running
		}
	})
	checkTimer.Start(100) // check every 100ms
}
