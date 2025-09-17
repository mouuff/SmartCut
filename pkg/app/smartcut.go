package app

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/mouuff/SmartCuts/pkg/generator"
	"github.com/mouuff/SmartCuts/pkg/types"
	"golang.design/x/clipboard"
)

type Item struct {
	Title   string
	Content string
}

type SmartCutApp struct {
	items         []Item
	listContainer *fyne.Container
	window        fyne.Window

	Config *types.SmartCutConfig
}

func NewSmartCutApp(w fyne.Window, config *types.SmartCutConfig) *SmartCutApp {

	items := make([]Item, 0)
	for _, hook := range config.PromptConfigs {
		items = append(items, Item{
			Title:   hook.Title,
			Content: "Waiting for generation...",
		})
	}

	la := &SmartCutApp{
		items:         items,
		listContainer: container.NewVBox(),
		window:        w,
		Config:        config,
	}

	// Render initial list
	la.RefreshList()

	return la
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
func (la *SmartCutApp) UpdateItem(result generator.GenerationResult) {
	la.items[result.PromptConfig.Index].Content = result.Text
	la.RefreshList()
}

// AddItem appends a new item and refreshes the view
func (la *SmartCutApp) AddItem(title, content string) {
	la.items = append(la.items, Item{Title: title, Content: content})
	la.RefreshList()
}

// Layout builds the full UI
func (la *SmartCutApp) Layout() fyne.CanvasObject {
	addBtn := widget.NewButton("Add Item", func() {
		newIndex := len(la.items) + 1
		la.AddItem(fmt.Sprintf("Item %d", newIndex),
			fmt.Sprintf("This is the content of item %d", newIndex))
	})

	return container.NewBorder(nil, addBtn, nil, nil, la.listContainer)
}
