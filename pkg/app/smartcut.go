package app

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/mouuff/SmartCuts/pkg/types"
	"github.com/mouuff/SmartCuts/pkg/utils"
	"golang.design/x/clipboard"
)

type Item struct {
	Title   string
	Content string
}

type SmartCutApp struct {
	rg            types.ResultsGenerator
	items         []Item
	listContainer *fyne.Container
	window        fyne.Window

	Config *types.SmartCutConfig
}

func NewSmartCutApp(
	w fyne.Window,
	config *types.SmartCutConfig,
	rg types.ResultsGenerator) *SmartCutApp {

	items := make([]Item, 0)
	for _, hook := range config.PromptConfigs {
		items = append(items, Item{
			Title:   hook.Title,
			Content: "Waiting for generation...",
		})
	}

	sc := &SmartCutApp{
		items:         items,
		listContainer: container.NewVBox(),
		window:        w,
		Config:        config,
		rg:            rg,
	}

	// Render initial list
	sc.RefreshList()

	return sc
}

func (sc *SmartCutApp) Start() {
	go func() {
		for res := range sc.rg.GetChannel() {
			fyne.Do(func() {
				sc.UpdateItem(res)
			})
		}
	}()
}

func (sc *SmartCutApp) RefreshList() {
	sc.listContainer.Objects = nil
	for _, item := range sc.items {
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
		content.SetMinRowsVisible(sc.Config.MinRowsVisible)

		// Let content expand horizontally
		contentContainer := container.NewStack(content)

		button := widget.NewButton("Copy", func(c string) func() {
			return func() {
				clipboard.Write(clipboard.FmtText, []byte(c))
				fmt.Println(c)
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

// AddItem appends a new item and refreshes the view
func (sc *SmartCutApp) UpdateItem(result types.GenerationResult) {
	sc.items[result.PromptConfig.Index].Content = result.Text
	sc.RefreshList()

	if result.IsExplicit {
		sc.window.RequestFocus()
	}
}

func (sc *SmartCutApp) Layout() fyne.CanvasObject {
	addBtn := widget.NewButton("Generate from clipboard", func() {
		rawclip := clipboard.Read(clipboard.FmtText)

		if rawclip != nil {
			sc.rg.GenerateForInput(types.InputResult{
				Text:       string(rawclip),
				IsExplicit: true,
			})
		}
	})

	helpmsg := `Shortcut for processing the current clipboard: Alt+Shift+G
	Configuration file: {{configpath}}`

	configpath, err := utils.GetConfigurationFilePath()

	if err != nil {
		panic(err)
	}

	helpmsg = strings.ReplaceAll(helpmsg, "{{configpath}}", configpath)

	// Menu bar with Help
	menu := fyne.NewMainMenu(
		fyne.NewMenu("Help",
			fyne.NewMenuItem("About", func() {
				dialog.ShowInformation("About SmartCut", helpmsg, sc.window)
			}),
		),
	)

	sc.window.SetMainMenu(menu)

	return container.NewBorder(nil, addBtn, nil, nil, sc.listContainer)
}
