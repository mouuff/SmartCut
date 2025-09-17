package app

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Item struct {
	Title   string
	Content string
}

type SmartCutApp struct {
	items         []Item
	listContainer *fyne.Container
	window        fyne.Window
}

func NewSmartCutApp(w fyne.Window) *SmartCutApp {
	la := &SmartCutApp{
		items: []Item{
			{"Item 1", "This is the content of item 1"},
			{"Item 2", "This is the content of item 2"},
			{"Item 3", "This is the content of item 3"},
		},
		listContainer: container.NewVBox(),
		window:        w,
	}

	// Render initial list
	la.UpdateList()

	return la
}

// UpdateList refreshes the UI with current items
func (la *SmartCutApp) UpdateList() {
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
	la.UpdateList()
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
