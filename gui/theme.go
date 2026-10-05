package gui

import (
	"strings"

	"github.com/mappu/miqt/qt"
)

// Theme holds every color used by the application stylesheet.
type Theme struct {
	Name        string
	Bg          string
	Surface     string
	Border      string
	Hover       string
	Text        string
	Dim         string
	Placeholder string
	Accent      string
	AccentHover string
}

// Themes lists every selectable color theme, in roller order.
var Themes = []Theme{
	{
		Name:        "Green",
		Bg:          "#050605",
		Surface:     "#0f100e",
		Border:      "#262824",
		Hover:       "#40423e",
		Text:        "#ffffff",
		Dim:         "#e0e0e0",
		Placeholder: "#9a9c94",
		Accent:      "#8da079",
		AccentHover: "#a6bb90",
	},
	{
		Name:        "Blue",
		Bg:          "#05070c",
		Surface:     "#0e1116",
		Border:      "#222a36",
		Hover:       "#38455a",
		Text:        "#ffffff",
		Dim:         "#d5dde8",
		Placeholder: "#98a4b5",
		Accent:      "#6f9fd8",
		AccentHover: "#8bb6e6",
	},
	{
		Name:        "Red",
		Bg:          "#0c0505",
		Surface:     "#160e0e",
		Border:      "#362222",
		Hover:       "#543030",
		Text:        "#ffffff",
		Dim:         "#e8d5d5",
		Placeholder: "#b89a9a",
		Accent:      "#d86f6f",
		AccentHover: "#e68a8a",
	},
	{
		Name:        "Yellow",
		Bg:          "#0b0a04",
		Surface:     "#16140d",
		Border:      "#36311f",
		Hover:       "#544b2f",
		Text:        "#ffffff",
		Dim:         "#ece5cf",
		Placeholder: "#b8b094",
		Accent:      "#d6bf6f",
		AccentHover: "#e6d18a",
	},
	{
		Name:        "Dark",
		Bg:          "#0a0a0c",
		Surface:     "#141416",
		Border:      "#2a2a2e",
		Hover:       "#3d3d43",
		Text:        "#f2f2f2",
		Dim:         "#b8b8bc",
		Placeholder: "#8e8e94",
		Accent:      "#8a8a90",
		AccentHover: "#a5a5ab",
	},
	{
		Name:        "White",
		Bg:          "#f4f4f2",
		Surface:     "#ffffff",
		Border:      "#d0d0cc",
		Hover:       "#e8e8e4",
		Text:        "#1a1a1a",
		Dim:         "#4a4a4a",
		Placeholder: "#8b8b87",
		Accent:      "#6d8f5a",
		AccentHover: "#55753f",
	},
	{
		Name:        "Orange",
		Bg:          "#0c0803",
		Surface:     "#16100a",
		Border:      "#362a1c",
		Hover:       "#54402a",
		Text:        "#ffffff",
		Dim:         "#ecdfd0",
		Placeholder: "#b8a795",
		Accent:      "#d68f4f",
		AccentHover: "#e8a468",
	},
	{
		Name:        "Purple",
		Bg:          "#08050c",
		Surface:     "#120e17",
		Border:      "#2a2236",
		Hover:       "#423355",
		Text:        "#ffffff",
		Dim:         "#ddd2ec",
		Placeholder: "#a99bbb",
		Accent:      "#9d7fd6",
		AccentHover: "#b295e6",
	},
}

// DefaultThemeName is used when no theme has been chosen yet.
const DefaultThemeName = "Dark"

// ThemeIndex returns the roller position of a theme name
func ThemeIndex(name string) int {
	for i, t := range Themes {
		if t.Name == name {
			return i
		}
	}
	for i, t := range Themes {
		if t.Name == DefaultThemeName {
			return i
		}
	}
	return 0
}

// ThemeByName returns the theme with the given name, falling back to the default.
func ThemeByName(name string) Theme {
	return Themes[ThemeIndex(name)]
}

// themeTarget holds the widgets that receive the stylesheet on every theme change.
var themeTarget struct {
	app    *qt.QApplication
	window *qt.QMainWindow
}

// BindThemeTarget registers the application and main window that ApplyTheme updates.
func BindThemeTarget(app *qt.QApplication, window *qt.QMainWindow) {
	themeTarget.app = app
	themeTarget.window = window
}

// ApplyTheme installs the theme stylesheet application-wide and repolishes the window.
func ApplyTheme(theme Theme, fontFamily string) {
	sheet := theme.StyleSheet(fontFamily)
	if themeTarget.app != nil {
		themeTarget.app.SetStyleSheet(sheet)
	}
	if themeTarget.window != nil {
		themeTarget.window.SetStyleSheet(sheet)
	}
}

// StyleSheet renders the application stylesheet for this theme.
func (t Theme) StyleSheet(fontFamily string) string {
	if fontFamily == "" {
		fontFamily = "Roboto"
	}
	return strings.NewReplacer(
		"{{Bg}}", t.Bg,
		"{{Surface}}", t.Surface,
		"{{Border}}", t.Border,
		"{{Hover}}", t.Hover,
		"{{Text}}", t.Text,
		"{{Dim}}", t.Dim,
		"{{Placeholder}}", t.Placeholder,
		"{{Accent}}", t.Accent,
		"{{AccentHover}}", t.AccentHover,
		"{{FontFamily}}", fontFamily,
	).Replace(styleSheetTemplate)
}

const styleSheetTemplate = `
	QWidget {
		background-color: {{Bg}};
		color: {{Text}};
		font-family: '{{FontFamily}}', sans-serif;
	}
	QMainWindow {
		background-color: {{Bg}};
	}
	QLabel {
		color: {{Text}};
		background: transparent;
	}
	QPushButton {
		background-color: {{Surface}};
		border: 1px solid {{Border}};
		border-radius: 0px;
		padding: 8px 16px;
		color: {{Text}};
		font-weight: bold;
		font-size: 14px;
	}
	QPushButton:hover {
		background-color: {{Hover}};
		border: 1px solid {{Accent}};
	}
	QPushButton:pressed {
		background-color: {{Bg}};
	}
	QPushButton:disabled {
		background-color: {{Bg}};
		border: 1px solid {{Border}};
		color: {{Dim}};
	}
	QTabWidget::pane {
		border: none;
		background: {{Bg}};
		color: {{Text}};
	}
	QTabBar {
		font-size: 13px;
		font-weight: bold;
	}
	QTabBar::tab {
		background: {{Surface}};
		color: {{Dim}};
		padding: 10px 16px;
		border-bottom: 2px solid transparent;
	}
	QTabBar::tab:selected {
		color: {{Accent}};
		background: {{Bg}};
		border-bottom: 2px solid {{Accent}};
	}
	QTabBar::tab:hover:!selected {
		background: {{Hover}};
	}
	QGroupBox {
		font-weight: bold;
		font-size: 14px;
		color: {{Accent}};
		border: 1px solid {{Border}};
		margin-top: 18px;
		border-radius: 0px;
	}
	QGroupBox::title {
		subcontrol-origin: margin;
		subcontrol-position: top left;
		padding: 0 5px;
		left: 10px;
		color: {{Accent}};
	}
	QLineEdit {
		background-color: {{Surface}};
		border: 1px solid {{Border}};
		padding: 6px 10px;
		color: {{Text}};
		border-radius: 0px;
	}
	QLineEdit:focus {
		border: 1px solid {{Accent}};
	}
	QLineEdit::placeholder {
		color: {{Placeholder}};
	}
	QComboBox {
		background-color: {{Surface}};
		border: 1px solid {{Border}};
		padding: 6px 10px;
		color: {{Text}};
		border-radius: 0px;
		min-height: 20px;
	}
	QComboBox::drop-down {
		border-left: 1px solid {{Border}};
		width: 30px;
	}
	QComboBox:focus {
		border: 1px solid {{Accent}};
	}
	QComboBox QAbstractItemView {
		background-color: {{Surface}};
		color: {{Text}};
		selection-background-color: {{Hover}};
		border: 1px solid {{Border}};
		border-radius: 0px;
	}
	QListWidget {
		background-color: {{Surface}};
		border: 1px solid {{Border}};
		padding: 4px;
		color: {{Text}};
		outline: none;
	}
	QListWidget::item {
		padding: 4px 6px;
	}
	QListWidget::item:selected {
		background-color: {{Hover}};
		color: {{Text}};
	}
	QListWidget::item:hover {
		background-color: {{Hover}};
	}
	QLabel#VstDetail {
		background-color: {{Surface}};
		border: 1px solid {{Border}};
		padding: 8px 10px;
		color: {{Dim}};
		font-size: 11px;
	}
	QProgressBar {
		border: 1px solid {{Border}};
		background: {{Surface}};
		text-align: center;
		color: {{Text}};
		border-radius: 0px;
	}
	QProgressBar::chunk {
		background-color: {{Accent}};
		border-radius: 0px;
	}
	QRadioButton {
		color: {{Text}};
	}
	QRadioButton::indicator {
		width: 14px;
		height: 14px;
		border-radius: 0px;
		border: 2px solid {{Border}};
		background: {{Surface}};
	}
	QRadioButton::indicator:checked {
		background: {{Accent}};
		border: 2px solid {{Accent}};
	}
	QSlider::groove:horizontal {
		height: 6px;
		background: {{Border}};
		border-radius: 0px;
	}
	QSlider::sub-page:horizontal {
		background: {{Accent}};
		border-radius: 0px;
	}
	QSlider::handle:horizontal {
		width: 18px;
		height: 18px;
		margin: -6px 0;
		border-radius: 0px;
		background: {{Accent}};
		border: 2px solid {{Bg}};
	}
	QSlider::handle:horizontal:hover {
		background: {{AccentHover}};
	}
	QSlider::add-page {
		background: {{Border}};
	}
	QStatusBar {
		background: {{Surface}};
		color: {{Dim}};
	}
	QStatusBar::item {
		border: none;
	}
	QToolTip {
		background-color: {{Surface}};
		color: {{Text}};
		border: 1px solid {{Accent}};
		padding: 4px;
	}
	QFrame {
		color: {{Text}};
		background: transparent;
		border: none;
	}
	QFrame#separator {
		background: {{Border}};
		border: none;
		min-height: 1px;
		max-height: 1px;
	}
	QScrollBar:vertical {
		background: {{Bg}};
		width: 10px;
		margin: 0;
	}
	QScrollBar::handle:vertical {
		background: {{Border}};
		min-height: 24px;
		border-radius: 0px;
	}
	QScrollBar::handle:vertical:hover {
		background: {{AccentHover}};
	}
	QScrollBar:horizontal {
		background: {{Bg}};
		height: 10px;
		margin: 0;
	}
	QScrollBar::handle:horizontal {
		background: {{Border}};
		min-width: 24px;
		border-radius: 0px;
	}
	QScrollBar::handle:horizontal:hover {
		background: {{AccentHover}};
	}
	QScrollBar::add-line, QScrollBar::sub-line {
		height: 0;
		width: 0;
	}
	QWidget#TitleBar {
		background-color: {{Surface}};
		border-bottom: 1px solid {{Border}};
	}
	QPushButton#TitleBarBtn {
		background-color: transparent;
		border: none;
		border-radius: 0px;
		padding: 0px;
		font-size: 16px;
	}
	QPushButton#TitleBarBtn:hover {
		background-color: {{Hover}};
	}
	QPushButton#TitleBarCloseBtn {
		background-color: transparent;
		border: none;
		border-radius: 0px;
		padding: 0px;
		font-size: 16px;
	}
	QPushButton#TitleBarCloseBtn:hover {
		background-color: #d86f6f;
		color: #ffffff;
	}
`
