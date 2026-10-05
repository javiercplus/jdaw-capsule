package gui

import (
	"jdaw-capsule/backend"

	"github.com/mappu/miqt/qt"
)

func NewMainWindow(app *qt.QApplication) *qt.QMainWindow {
	env, _ := backend.NewEnvironment()
	window := qt.NewQMainWindow(nil)
	window.SetWindowTitle("JDAW Capsule")
	window.SetWindowFlags(qt.FramelessWindowHint)
	window.SetMinimumSize2(860, 560)
	ApplyAppIcon(window)

	env.Config.Theme = ThemeByName(env.Config.Theme).Name
	BindThemeTarget(app, window)
	ApplyTheme(ThemeByName(env.Config.Theme), ResolveFontFamily(env.Config.FontFamily))

	centralWidget := qt.NewQWidget(window.QWidget)
	window.SetCentralWidget(centralWidget)
	mainLayout := qt.NewQVBoxLayout(centralWidget)
	mainLayout.SetContentsMargins(2, 2, 2, 2)
	mainLayout.SetSpacing(0)

	// Custom Title Bar
	titleBar := qt.NewQWidget(centralWidget)
	titleBar.SetObjectName("TitleBar")
	titleLayout := qt.NewQHBoxLayout(titleBar)
	titleLayout.SetContentsMargins(10, 5, 10, 5)

	titleIcon := qt.NewQLabel(titleBar)
	titleLabel := qt.NewQLabel(titleBar)

	titleLayout.AddWidget(titleIcon.QWidget)
	titleLayout.AddSpacing(5)
	titleLayout.AddWidget(titleLabel.QWidget)
	titleLayout.AddStretch()

	btnMin := qt.NewQPushButton(titleBar)
	btnMin.SetObjectName("TitleBarBtn")
	btnMin.SetText("—")
	btnMin.SetFixedSize2(30, 30)
	btnMin.SetCursor(qt.NewQCursor2(qt.PointingHandCursor))
	btnMin.OnClicked(func() { window.ShowMinimized() })

	btnMax := qt.NewQPushButton(titleBar)
	btnMax.SetObjectName("TitleBarBtn")
	btnMax.SetText("⬜")
	btnMax.SetFixedSize2(30, 30)
	btnMax.SetCursor(qt.NewQCursor2(qt.PointingHandCursor))
	btnMax.OnClicked(func() {
		if window.IsMaximized() {
			window.ShowNormal()
			btnMax.SetText("⬜")
		} else {
			window.ShowMaximized()
			btnMax.SetText("⬜")
		}
	})

	btnClose := qt.NewQPushButton(titleBar)
	btnClose.SetObjectName("TitleBarCloseBtn")
	btnClose.SetText("✕")
	btnClose.SetFixedSize2(30, 30)
	btnClose.SetCursor(qt.NewQCursor2(qt.PointingHandCursor))
	btnClose.OnClicked(func() { window.Close() })

	titleLayout.AddWidget(btnMin.QWidget)
	titleLayout.AddWidget(btnMax.QWidget)
	titleLayout.AddWidget(btnClose.QWidget)

	mainLayout.AddWidget(titleBar)

	// Dragging logic for custom title bar
	var dragPos *qt.QPoint
	titleBar.OnMousePressEvent(func(super func(event *qt.QMouseEvent), event *qt.QMouseEvent) {
		if event.Button() == qt.LeftButton {
			dragPos = qt.NewQPoint2(event.GlobalX(), event.GlobalY())
		}
		super(event)
	})
	titleBar.OnMouseMoveEvent(func(super func(event *qt.QMouseEvent), event *qt.QMouseEvent) {
		if dragPos != nil {
			newPos := qt.NewQPoint2(event.GlobalX(), event.GlobalY())
			dx := newPos.X() - dragPos.X()
			dy := newPos.Y() - dragPos.Y()
			window.Move(window.X()+dx, window.Y()+dy)
			dragPos = newPos
		}
		super(event)
	})
	titleBar.OnMouseReleaseEvent(func(super func(event *qt.QMouseEvent), event *qt.QMouseEvent) {
		dragPos = nil
		super(event)
	})

	contentWidget := qt.NewQWidget(centralWidget)
	layout := qt.NewQVBoxLayout(contentWidget)
	layout.SetContentsMargins(8, 8, 8, 8)
	layout.SetSpacing(0)
	mainLayout.AddWidget(contentWidget)

	statusBar := qt.NewQStatusBar(window.QWidget)
	window.SetStatusBar(statusBar)
	statusBar.ShowMessage("Ready")
	
	tabWidget := qt.NewQTabWidget(contentWidget)

	// Keep every tab label fully readable: no eliding, no stretching, and
	// scroll buttons instead of squeezing the labels when space is tight.
	tabWidget.SetElideMode(qt.ElideNone)
	tabWidget.SetUsesScrollButtons(true)
	tabWidget.TabBar().SetExpanding(false)
	tabWidget.TabBar().SetDrawBase(false)

	tabStart := CreateTabStart(env, statusBar)
	tabSettings := CreateTabSettings(env, statusBar, window)
	tabVsts := CreateTabVsts(env, statusBar)
	tabInstallWine := CreateTabInstall(env, statusBar)
	tabInstallReaper := CreateTabInstallReaper(env, statusBar)
	tabInstallYabridge := CreateTabInstallYabridge(env, statusBar)
	tabAbout := CreateTabAbout(nil, statusBar)
	// Add tabs
	tabWidget.InsertTab(0, tabStart, "Start")
	tabWidget.InsertTab(1, tabSettings, "Settings")
	tabWidget.InsertTab(2, tabVsts, "Windows VSTs")
	tabWidget.InsertTab(3, tabInstallWine, "Install Wine")
	tabWidget.InsertTab(4, tabInstallReaper, "Install REAPER")
	tabWidget.InsertTab(5, tabInstallYabridge, "Install yabridge")
	tabWidget.InsertTab(6, tabAbout, "About")
	layout.AddWidget(tabWidget.QWidget)
	window.Resize(960, 640)
	return window
}
