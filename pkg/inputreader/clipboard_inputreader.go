package inputreader

import (
	"context"

	"github.com/mouuff/SmartCuts/pkg/types"
	"golang.design/x/clipboard"
)

type ClipboardInputReader struct {
	ctx     context.Context
	OnInput func(types.InputText)
}

func NewClipboardInputReader(ctx context.Context) *ClipboardInputReader {
	return &ClipboardInputReader{
		ctx:     ctx,
		OnInput: func(types.InputText) {},
	}
}

func (c *ClipboardInputReader) Start() {

	clipch := clipboard.Watch(c.ctx, clipboard.FmtText)

	go func() {
		for data := range clipch {
			c.OnInput(types.InputText{
				Text:       string(data),
				IsExplicit: false,
			})
		}
	}()
}
