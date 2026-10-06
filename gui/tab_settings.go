package gui

import (
	"jdaw-capsule/backend"
	"strings"

	"github.com/mappu/miqt/qt"
)

func CreateTabSettings(env *backend.Environment, statusBar *qt.QStatusBar, window *qt.QMainWindow) *qt.QWidget {
	widget := qt.NewQWidget(nil)
	layout := qt.NewQVBoxLayout(widget)
	layout.SetContentsMargins(24, 24, 24, 24)
	layout.SetSpacing(12)

	// --- Appearance group ---
	groupTheme := qt.NewQGroupBox(widget)
	groupTheme.SetTitle("Appearance")
	layoutTheme := qt.NewQVBoxLayout(groupTheme.QWidget)
	layoutTheme.SetContentsMargins(12, 16, 12, 12)
	layoutTheme.SetSpacing(12)

	// Theme Selection
	themeWidget := qt.NewQWidget(groupTheme.QWidget)
	themeLayout := qt.NewQHBoxLayout(themeWidget)
	themeLayout.SetContentsMargins(0, 0, 0, 0)
	themeLayout.SetSpacing(10)

	themeNameLabel := qt.NewQLabel(themeWidget)
	themeNameLabel.SetText(ThemeByName(env.Config.Theme).Name)
	themeNameFont := qt.NewQFont()
	themeNameFont.SetPointSize(12)
	themeNameFont.SetBold(true)
	themeNameLabel.SetFont(themeNameFont)

	roller := qt.NewQSlider4(qt.Horizontal, themeWidget)
	roller.SetMinimumWidth(240)
	roller.SetRange(0, len(Themes)-1)
	roller.SetSingleStep(1)
	roller.SetPageStep(1)
	roller.SetValue(ThemeIndex(env.Config.Theme))

	roller.OnValueChanged(func(value int) {
		if value < 0 || value >= len(Themes) {
			return
		}
		theme := Themes[value]
		themeNameLabel.SetText(theme.Name)
		env.Config.Theme = theme.Name
		env.SaveConfig()
		ApplyTheme(theme, env.Config.FontFamily)
		statusBar.ShowMessage2("Theme: "+theme.Name, 2000)
	})

	themeLayout.AddWidget2(roller.QWidget, 1)
	themeLayout.AddWidget(themeNameLabel.QWidget)
	layoutTheme.AddWidget(themeWidget)

	// Window Size Selection
	sizeWidget := qt.NewQWidget(groupTheme.QWidget)
	sizeLayout := qt.NewQHBoxLayout(sizeWidget)
	sizeLayout.SetContentsMargins(0, 0, 0, 0)
	sizeLayout.SetSpacing(10)

	sizeLabel := qt.NewQLabel(sizeWidget)
	sizeLabel.SetText("Window Size (Disabled if maximized)")
	
	sizeSlider := qt.NewQSlider4(qt.Horizontal, sizeWidget)
	sizeSlider.SetRange(80, 150)
	sizeSlider.SetValue(100)
	
	sizeSlider.OnValueChanged(func(val int) {
		if !window.IsMaximized() {
			baseW, baseH := 960, 640
			newW := baseW * val / 100
			newH := baseH * val / 100
			window.Resize(newW, newH)
		}
	})

	sizeLayout.AddWidget(sizeLabel.QWidget)
	sizeLayout.AddWidget2(sizeSlider.QWidget, 1)
	layoutTheme.AddWidget(sizeWidget)

	// Font Selection
	fontWidget := qt.NewQWidget(groupTheme.QWidget)
	fontLayout := qt.NewQHBoxLayout(fontWidget)
	fontLayout.SetContentsMargins(0, 0, 0, 0)
	fontLayout.SetSpacing(10)

	fontLabel := qt.NewQLabel(fontWidget)
	fontLabel.SetText("Font Family")

	fontCombo := qt.NewQFontComboBox(fontWidget)
	fontCombo.SetCurrentFont(qt.NewQFont6(ResolveFontFamily(env.Config.FontFamily), 12))

	fontCombo.OnCurrentFontChanged(func(f *qt.QFont) {
		env.Config.FontFamily = f.Family()
		env.SaveConfig()
		ApplyTheme(ThemeByName(env.Config.Theme), ResolveFontFamily(env.Config.FontFamily))
		statusBar.ShowMessage2("Font: "+env.Config.FontFamily, 2000)
	})

	fontLayout.AddWidget(fontLabel.QWidget)
	fontLayout.AddWidget2(fontCombo.QWidget, 1)
	layoutTheme.AddWidget(fontWidget)

	themeHint := qt.NewQLabel(groupTheme.QWidget)
	themeHint.SetText("Roll to switch between: " + themeNames() + ".")
	themeHint.SetWordWrap(true)
	hintFont := qt.NewQFont()
	hintFont.SetPointSize(9)
	themeHint.SetFont(hintFont)
	layoutTheme.AddWidget(themeHint.QWidget)

	layout.AddWidget(groupTheme.QWidget)

	// --- Audio Backend group ---
	groupAudio := qt.NewQGroupBox(widget)
	groupAudio.SetTitle("Audio Backend")
	layoutAudio := qt.NewQVBoxLayout(groupAudio.QWidget)
	layoutAudio.SetContentsMargins(12, 16, 12, 12)
	layoutAudio.SetSpacing(6)

	radioAuto := qt.NewQRadioButton(groupAudio.QWidget)
	radioAuto.SetText("Automatic (detect PipeWire/PulseAudio)")

	radioPipe := qt.NewQRadioButton(groupAudio.QWidget)
	radioPipe.SetText("PipeWire (pw-jack)")

	radioPulse := qt.NewQRadioButton(groupAudio.QWidget)
	radioPulse.SetText("PulseAudio")

	if env.Config.AudioBackend == "Pipewire" {
		radioPipe.SetChecked(true)
	} else if env.Config.AudioBackend == "PulseAudio" {
		radioPulse.SetChecked(true)
	} else {
		radioAuto.SetChecked(true)
	}

	layoutAudio.AddWidget(radioAuto.QWidget)
	layoutAudio.AddWidget(radioPipe.QWidget)
	layoutAudio.AddWidget(radioPulse.QWidget)

	radioAuto.OnToggled(func(checked bool) {
		if checked {
			env.Config.AudioBackend = "Automatic"
			env.SaveConfig()
			statusBar.ShowMessage2("Audio backend: Automatic", 2000)
		}
	})
	radioPipe.OnToggled(func(checked bool) {
		if checked {
			env.Config.AudioBackend = "Pipewire"
			env.SaveConfig()
			statusBar.ShowMessage2("Audio backend: PipeWire", 2000)
		}
	})
	radioPulse.OnToggled(func(checked bool) {
		if checked {
			env.Config.AudioBackend = "PulseAudio"
			env.SaveConfig()
			statusBar.ShowMessage2("Audio backend: PulseAudio", 2000)
		}
	})

	layout.AddWidget(groupAudio.QWidget)

	// --- Plugins group ---
	groupPlugins := qt.NewQGroupBox(widget)
	groupPlugins.SetTitle("Windows VST Plugins")
	layoutPlugins := qt.NewQVBoxLayout(groupPlugins.QWidget)
	layoutPlugins.SetContentsMargins(12, 16, 12, 12)
	layoutPlugins.SetSpacing(8)

	pluginDescLabel := qt.NewQLabel(groupPlugins.QWidget)
	pluginDescLabel.SetText("Directory containing your Windows VST plugins (bridged via yabridge):")
	pluginDescLabel.SetWordWrap(true)
	layoutPlugins.AddWidget(pluginDescLabel.QWidget)

	// Path input row
	pathWidget := qt.NewQWidget(groupPlugins.QWidget)
	pathLayout := qt.NewQHBoxLayout(pathWidget)
	pathLayout.SetContentsMargins(0, 0, 0, 0)
	pathLayout.SetSpacing(6)

	inputPath := qt.NewQLineEdit(pathWidget)
	inputPath.SetText(env.Config.WindowsPluginsDir)
	inputPath.SetPlaceholderText("Select plugins directory...")

	btnBrowse := qt.NewQPushButton(pathWidget)
	btnBrowse.SetText("Browse...")

	btnBrowse.OnClicked(func() {
		dir := ShowFolderPicker(widget, "Select Windows Plugins Directory", env.Config.WindowsPluginsDir)
		if dir != "" {
			inputPath.SetText(dir)
			env.Config.WindowsPluginsDir = dir
			env.SaveConfig()
			statusBar.ShowMessage2("Plugins directory updated", 2000)
		}
	})

	inputPath.OnTextChanged(func(text string) {
		env.Config.WindowsPluginsDir = text
		env.SaveConfig()
	})

	pathLayout.AddWidget(inputPath.QWidget)
	pathLayout.AddWidget(btnBrowse.QWidget)

	layoutPlugins.AddWidget(pathWidget)

	layout.AddWidget(groupPlugins.QWidget)

	// --- Stretch at the bottom ---
	layout.AddStretch()

	return widget
}

// themeNames returns the selectable theme names as a readable list.
func themeNames() string {
	names := make([]string, 0, len(Themes))
	for _, t := range Themes {
		names = append(names, t.Name)
	}
	return strings.Join(names, ", ")
}
