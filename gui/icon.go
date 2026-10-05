package gui

import (
	"sync"

	"jdaw-capsule/asset"

	"github.com/mappu/miqt/qt"
)

// DesktopFileName identifies this binary to the desktop shell, so the window
// icon and the launcher entry resolve to the same app.
const DesktopFileName = "jdaw-capsule"

// AppLogoSize is the logical edge, in pixels, of the logo shown above the title.
const AppLogoSize = 196

var (
	appIconOnce sync.Once
	appIcon     *qt.QIcon
)

// AppIcon builds the application icon from the raster sizes compiled into the
// binary, so the window and the taskbar show it without touching the host
// icon theme. It builds once and returns the same icon on later calls.
func AppIcon() *qt.QIcon {
	appIconOnce.Do(func() {
		appIcon = qt.NewQIcon()
		for _, size := range asset.IconSizes {
			data, ok := asset.IconPNG(size)
			if !ok {
				continue
			}
			pixmap := qt.NewQPixmap()
			if !pixmap.LoadFromDataWithData(data) {
				continue
			}
			appIcon.AddPixmap2(pixmap, qt.QIcon__Normal)
		}
	})
	return appIcon
}

// ApplyAppIcon sets the app icon on the window and names the desktop entry so
// the window manager can match the taskbar icon to the launcher.
func ApplyAppIcon(window *qt.QMainWindow) {
	window.SetWindowIcon(AppIcon())
	qt.QGuiApplication_SetDesktopFileName(DesktopFileName)
}

// logoSourceSize is the embedded size the logo is rasterised from, big enough to
// cover any display scale.
const logoSourceSize = 256

// AppLogoPixmap draws the app icon at the given logical size, rendered at the
// device pixel ratio of the widget it will be shown on so it stays sharp on
// HiDPI screens. It returns nil when the icon is unavailable.
func AppLogoPixmap(widget *qt.QWidget, logicalSize int) *qt.QPixmap {
	data, ok := asset.IconPNG(logoSourceSize)
	if !ok {
		return nil
	}
	source := qt.NewQPixmap()
	if !source.LoadFromDataWithData(data) {
		return nil
	}
	ratio := widget.DevicePixelRatioF()
	if ratio <= 0 {
		ratio = 1
	}
	physical := int(float64(logicalSize)*ratio + 0.5)
	logo := source.Scaled3(physical, physical, qt.KeepAspectRatio, qt.SmoothTransformation)
	logo.SetDevicePixelRatio(ratio)
	return logo
}
