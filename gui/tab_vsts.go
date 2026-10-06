package gui

import (
	"fmt"
	"jdaw-capsule/backend"

	"github.com/mappu/miqt/qt"
)

func CreateTabVsts(env *backend.Environment, statusBar *qt.QStatusBar) *qt.QWidget {
	widget := qt.NewQWidget(nil)
	layout := qt.NewQVBoxLayout(widget)
	layout.SetContentsMargins(24, 24, 24, 24)
	layout.SetSpacing(12)

	// --- Windows VST group ---
	group := qt.NewQGroupBox(widget)
	group.SetTitle("Windows VST Plugins")
	groupLayout := qt.NewQVBoxLayout(group.QWidget)
	groupLayout.SetContentsMargins(12, 16, 12, 12)
	groupLayout.SetSpacing(8)

	pathLabel := qt.NewQLabel(group.QWidget)
	pathLabel.SetText(env.Config.WindowsPluginsDir)
	pathLabel.SetWordWrap(true)
	pathLabel.SetTextInteractionFlags(qt.TextSelectableByMouse)
	pathFont := qt.NewQFont()
	pathFont.SetPointSize(9)
	pathLabel.SetFont(pathFont)
	groupLayout.AddWidget(pathLabel.QWidget)

	// --- Plugin list ---
	list := qt.NewQListWidget(group.QWidget)
	list.SetSelectionMode(qt.QAbstractItemView__SingleSelection)
	list.SetUniformItemSizes(true)
	list.SetAlternatingRowColors(false)
	list.SetMinimumHeight(240)
	groupLayout.AddWidget(list.QWidget)

	// --- Buttons ---
	buttonRow := qt.NewQWidget(group.QWidget)
	buttonLayout := qt.NewQHBoxLayout(buttonRow)
	buttonLayout.SetContentsMargins(0, 0, 0, 0)
	buttonLayout.SetSpacing(8)

	btnOpenFolder := qt.NewQPushButton(buttonRow)
	btnOpenFolder.SetText("Open VST Folder")
	btnOpenFolder.SetCursor(qt.NewQCursor2(qt.PointingHandCursor))

	btnRefresh := qt.NewQPushButton(buttonRow)
	btnRefresh.SetText("Refresh List")
	btnRefresh.SetCursor(qt.NewQCursor2(qt.PointingHandCursor))

	buttonLayout.AddWidget(btnOpenFolder.QWidget)
	buttonLayout.AddWidget(btnRefresh.QWidget)
	buttonLayout.AddStretch()
	groupLayout.AddWidget(buttonRow)

	layout.AddWidget(group.QWidget)

	// --- Selected plugin details ---
	detailLabel := qt.NewQLabel(widget)
	detailLabel.SetObjectName("VstDetail")
	detailLabel.SetWordWrap(true)
	detailLabel.SetMinimumHeight(40)
	detailLabel.SetAlignment(qt.AlignTop | qt.AlignLeft)
	detailLabel.SetText("No plugin selected.")
	layout.AddWidget(detailLabel.QWidget)

	// --- Summary ---
	summaryLabel := qt.NewQLabel(widget)
	summaryLabel.SetAlignment(qt.AlignCenter)
	summaryFont := qt.NewQFont()
	summaryFont.SetPointSize(9)
	summaryLabel.SetFont(summaryFont)
	layout.AddWidget(summaryLabel.QWidget)

	layout.AddStretch()

	// The scan result is kept next to the list so a selection maps back to the
	// plugin it came from without storing anything inside the item.
	var plugins []backend.Plugin

	refresh := func() {
		list.Clear()
		plugins = nil
		detailLabel.SetText("No plugin selected.")

		found, err := env.ScanPlugins()
		if err != nil {
			list.AddItem("Could not read the plugins directory")
			summaryLabel.SetText(err.Error())
			statusBar.ShowMessage2("VST scan failed", 3000)
			return
		}

		plugins = found
		pathLabel.SetText(env.Config.WindowsPluginsDir)
		if len(plugins) == 0 {
			list.AddItem("No plugins found in this directory")
			summaryLabel.SetText("0 plugins")
			statusBar.ShowMessage2("No Windows plugins found", 3000)
			return
		}

		for _, plugin := range plugins {
			list.AddItem(fmt.Sprintf("%-40s  %-5s  %s", plugin.Name, plugin.Format, humanSize(plugin.Size)))
		}
		summaryLabel.SetText(fmt.Sprintf("%d plugins  •  %s total", len(plugins), humanSize(totalSize(plugins))))
		statusBar.ShowMessage2(fmt.Sprintf("%d Windows plugins found", len(plugins)), 3000)
	}

	list.OnCurrentRowChanged(func(row int) {
		if row < 0 || row >= len(plugins) {
			detailLabel.SetText("No plugin selected.")
			return
		}
		plugin := plugins[row]
		detailLabel.SetText(fmt.Sprintf("%s\n%s\n%s  •  modified %s",
			plugin.Path, plugin.Format, humanSize(plugin.Size), plugin.Modified.Format("2006-01-02 15:04")))
	})

	btnRefresh.OnClicked(func() {
		statusBar.ShowMessage("Scanning plugins...")
		refresh()
	})

	btnOpenFolder.OnClicked(func() {
		dir := env.Config.WindowsPluginsDir
		if dir == "" {
			statusBar.ShowMessage2("No plugins directory configured", 3000)
			return
		}
		if err := backend.OpenInFileManager(dir); err != nil {
			statusBar.ShowMessage2("Could not open: "+err.Error(), 4000)
			return
		}
		statusBar.ShowMessage2("Opened "+dir, 3000)
	})

	refresh()
	return widget
}

// totalSize adds up the size of every plugin in the list.
func totalSize(plugins []backend.Plugin) int64 {
	var total int64
	for _, plugin := range plugins {
		total += plugin.Size
	}
	return total
}

// humanSize renders a byte count with the largest unit that keeps it readable.
func humanSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for size := bytes / unit; size >= unit && exp < 3; size /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(bytes)/float64(div), "KMGT"[exp])
}
