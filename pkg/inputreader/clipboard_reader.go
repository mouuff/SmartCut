package inputreader

import (
	"context"

	"github.com/mouuff/SmartCut/pkg/types"
	"golang.design/x/clipboard"
)

type ClipboardReader struct {
	ctx     context.Context
	OnInput func(types.InputText)
}

func NewClipboardReader(ctx context.Context) *ClipboardReader {
	return &ClipboardReader{
		ctx:     ctx,
		OnInput: func(types.InputText) {},
	}
}

func (c *ClipboardReader) Start() {

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
