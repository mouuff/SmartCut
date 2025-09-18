package inputreader

import (
	"context"

	"github.com/mouuff/SmartCuts/pkg/types"
	"golang.design/x/clipboard"
)

type ClipboardInputReader struct {
	context context.Context
	OnInput func(types.InputText)
}

func NewClipboardInputReader(context context.Context) *ClipboardInputReader {
	return &ClipboardInputReader{
		context: context,
		OnInput: func(types.InputText) {},
	}
}

func (c *ClipboardInputReader) Start() {

	clipch := clipboard.Watch(c.context, clipboard.FmtText)

	go func() {
		for data := range clipch {
			c.OnInput(types.InputText{
				Text:       string(data),
				IsExplicit: false,
			})
		}
	}()
}
