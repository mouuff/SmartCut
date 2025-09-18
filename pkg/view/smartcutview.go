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

type SmartCutView struct {
	listContainer *fyne.Container
	window        fyne.Window
	OnAskGenerate func(types.InputText)

	model *types.SmartCutModel
}

func NewSmartCutView(w fyne.Window, m *types.SmartCutModel) *SmartCutView {
	sc := &SmartCutView{
		listContainer: container.NewVBox(),
		window:        w,
		model:         m,
		OnAskGenerate: func(types.InputText) {},
	}

	m.OnChanged = sc.Refresh
	return sc
}

func (sc *SmartCutView) Refresh() {
	fyne.Do(func() {
		sc.refreshListContainer()
	})
}

func (sc *SmartCutView) RequestFocus() {
	fyne.Do(func() {
		sc.window.RequestFocus()
	})
}

func (sc *SmartCutView) refreshListContainer() {
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
		content.SetMinRowsVisible(sc.model.MinRowsVisible())

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

func (sc *SmartCutView) Layout() fyne.CanvasObject {
	addBtn := widget.NewButton("Generate from clipboard", func() {
		sc.OnAskGenerate(types.InputText{
			Text:       string(clipboard.Read(clipboard.FmtText)),
			IsExplicit: true,
		})
	})

	helpmsg := `Shortcut for processing the current clipboard: Alt+Shift+G
	Configuration file: {{configpath}}`

	helpmsg = strings.ReplaceAll(helpmsg, "{{configpath}}", utils.GetConfigurationFilePath())

	// Menu bar with Help
	menu := fyne.NewMainMenu(
		fyne.NewMenu("Menu",
			fyne.NewMenuItem("Help", func() {
				dialog.ShowInformation("About SmartCut", helpmsg, sc.window)
			}),
			fyne.NewMenuItem("Configure", func() {
				utils.OpenFile(utils.GetConfigurationFilePath())
				dialog.ShowInformation("Warning", "You must restart the application to apply configuration changes", sc.window)
			}),
		),
	)

	sc.window.SetMainMenu(menu)

	return container.NewBorder(nil, addBtn, nil, nil, sc.listContainer)
}
