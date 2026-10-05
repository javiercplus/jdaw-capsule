package main

import (
	"jdaw-capsule/gui"
	"os"
	"github.com/mappu/miqt/qt"
)

func main() {
	app := qt.NewQApplication(os.Args)
	gui.RegisterEmbeddedFonts()
	mainWindow := gui.NewMainWindow(app)
	mainWindow.Show()

	os.Exit(qt.QApplication_Exec())
}
