package gui

import (
	"github.com/mappu/miqt/qt"
)

func ShowFolderPicker(parent *qt.QWidget, title string, startDir string) string {
	dialog := qt.NewQDialog(parent)
	dialog.SetWindowTitle(title)
	dialog.SetMinimumSize2(500, 400)
	
	layout := qt.NewQVBoxLayout(dialog.QWidget)

	// Current Path display
	pathLayout := qt.NewQHBoxLayout(nil)
	pathLabel := qt.NewQLabel(dialog.QWidget)
	pathLabel.SetText("Path:")
	pathEdit := qt.NewQLineEdit(dialog.QWidget)
	pathEdit.SetText(startDir)
	pathEdit.SetReadOnly(true)
	pathLayout.AddWidget(pathLabel.QWidget)
	pathLayout.AddWidget(pathEdit.QWidget)
	layout.AddLayout(pathLayout.QLayout)

	// File System Model
	model := qt.NewQFileSystemModel()
	model.SetFilter(qt.QDir__AllDirs | qt.QDir__NoDotAndDotDot | qt.QDir__Hidden)
	model.SetRootPath("")

	// Tree View
	tree := qt.NewQTreeView(dialog.QWidget)
	tree.SetModel(model.QAbstractItemModel)
	
	// Hide size, type, date columns
	tree.SetColumnHidden(1, true)
	tree.SetColumnHidden(2, true)
	tree.SetColumnHidden(3, true)

	// Set initial index
	initialIndex := model.Index2(startDir, 0)
	if initialIndex.IsValid() {
		tree.SetCurrentIndex(initialIndex)
		tree.ScrollTo(initialIndex, qt.QAbstractItemView__EnsureVisible)
		tree.Expand(initialIndex)
	}

	layout.AddWidget(tree.QWidget)

	tree.OnClicked(func(index *qt.QModelIndex) {
		pathEdit.SetText(model.FilePath(index))
	})

	// Buttons
	btnLayout := qt.NewQHBoxLayout(nil)
	btnLayout.AddStretch()

	btnCancel := qt.NewQPushButton(dialog.QWidget)
	btnCancel.SetText("Cancel")
	btnSelect := qt.NewQPushButton(dialog.QWidget)
	btnSelect.SetText("Select")
	
	btnLayout.AddWidget(btnCancel.QWidget)
	btnLayout.AddWidget(btnSelect.QWidget)
	layout.AddLayout(btnLayout.QLayout)

	selectedPath := ""

	btnCancel.OnClicked(func() {
		dialog.Reject()
	})

	btnSelect.OnClicked(func() {
		selectedPath = pathEdit.Text()
		dialog.Accept()
	})

	if dialog.Exec() == int(qt.QDialog__Accepted) {
		return selectedPath
	}

	return ""
}
