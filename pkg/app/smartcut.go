package app

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/mouuff/SmartCuts/pkg/orchestrator"
	"github.com/mouuff/SmartCuts/pkg/types"
)

type Item struct {
	Title   string
	Content string
}

type SmartCutApp struct {
	items         []Item
	listContainer *fyne.Container
	window        fyne.Window

	Ch     chan orchestrator.GenerationResult
	Config *types.SmartCutConfig
}

func NewSmartCutApp(w fyne.Window, config *types.SmartCutConfig, ch chan orchestrator.GenerationResult) *SmartCutApp {

	items := make([]Item, 0)
	for _, hook := range config.Hooks {
		items = append(items, Item{
			Title:   hook.Title,
			Content: "Waiting for generation...",
		})
	}

	la := &SmartCutApp{
		items:         items,
		listContainer: container.NewVBox(),
		window:        w,
		Ch:            ch,
		Config:        config,
	}

	// Render initial list
	la.RefreshList()

	return la
}

// RefreshList refreshes the UI with current items
func (la *SmartCutApp) RefreshList() {
	la.listContainer.Objects = nil
	for _, item := range la.items {
		title := widget.NewLabel(item.Title)
		content := widget.NewLabel(item.Content)
		button := widget.NewButton("Print", func(c string) func() {
			return func() {
				fmt.Println(c)
			}
		}(item.Content))

		row := container.NewBorder(nil, nil,
			container.NewVBox(title, content),
			button,
		)
		la.listContainer.Add(row)
	}
	la.listContainer.Refresh()
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
