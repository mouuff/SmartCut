package view

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/mouuff/SmartCuts/pkg/types"
	"github.com/mouuff/SmartCuts/pkg/utils"
	"golang.design/x/clipboard"
)

type SmartsCutView struct {
	window            fyne.Window
	listContainer     *fyne.Container
	model             *types.SmartCutsModel
	OnRequestGenerate func(types.InputText)
}

func NewSmartsCutView(w fyne.Window, m *types.SmartCutsModel) *SmartsCutView {
	sc := &SmartsCutView{
		listContainer:     container.NewVBox(),
		window:            w,
		model:             m,
		OnRequestGenerate: func(types.InputText) {},
	}

	m.OnChanged = sc.Refresh
	sc.refreshListResults()
	return sc
}

func (sc *SmartsCutView) Refresh() {
	fyne.Do(func() {
		sc.refreshListResults()
	})
}

func (sc *SmartsCutView) RequestFocus() {
	fyne.Do(func() {
		sc.window.RequestFocus()
	})
}

func (sc *SmartsCutView) refreshListResults() {
	sc.listContainer.Objects = nil
	for _, item := range sc.model.ResultItems() {
		title := widget.NewLabelWithStyle(item.Title, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

		// Multiline selectable Entry, but locked to read-only
		content := widget.NewEntry()
		content.SetText(item.Content)
		content.MultiLine = true
		content.Wrapping = fyne.TextWrapWord
		content.OnChanged = func(_ string) {
			// Reset text if user tries to type
			content.SetText(item.Content)
		}
		content.SetMinRowsVisible(sc.model.Config().MinRowsVisible)

		// Let content expand horizontally
		contentContainer := container.NewStack(content)

		button := widget.NewButton("Copy", func(c string) func() {
			return func() {
				clipboard.Write(clipboard.FmtText, []byte(c))
			}
		}(item.Content))

		// Row = title on top, content + button side by side
		row := container.NewVBox(
			title,
			container.NewBorder(nil, nil, nil, button, contentContainer),
		)

		sc.listContainer.Add(row)
	}

	sc.listContainer.Refresh()
}

func (sc *SmartsCutView) Layout() fyne.CanvasObject {
	addBtn := widget.NewButton("Generate from clipboard", func() {
		sc.OnRequestGenerate(types.InputText{
			Text:       string(clipboard.Read(clipboard.FmtText)),
			IsExplicit: true,
		})
	})

	// Menu bar with Help
	menu := fyne.NewMainMenu(
		fyne.NewMenu("Menu",
			fyne.NewMenuItem("Help", func() {
				helpmsg := "Shortcut for processing the current clipboard: Alt+Shift+G\nConfiguration file: {{configpath}}"
				helpmsg = strings.ReplaceAll(helpmsg, "{{configpath}}", sc.model.Config().ConfigPath)
				dialog.ShowInformation("About SmartCuts", helpmsg, sc.window)
			}),
			fyne.NewMenuItem("Configure", func() {
				utils.OpenFile(sc.model.Config().ConfigPath)
				dialog.ShowInformation("Warning", "You must restart the application to apply configuration changes", sc.window)
			}),
		),
	)

	sc.window.SetMainMenu(menu)

	return container.NewBorder(nil, addBtn, nil, nil, sc.listContainer)
}
