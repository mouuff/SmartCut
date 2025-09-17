package inputreader

import (
	"context"

	"github.com/mouuff/SmartCuts/pkg/types"
	"golang.design/x/clipboard"
)

type ClipboardInputReader struct {
	context context.Context
	ch      chan types.InputResult
}

func NewClipboardInputReader(context context.Context) *ClipboardInputReader {
	return &ClipboardInputReader{
		context: context,
		ch:      make(chan types.InputResult),
	}
}

func (c *ClipboardInputReader) Start() {

	clipch := clipboard.Watch(c.context, clipboard.FmtText)

	go func() {
		for data := range clipch {
			c.ch <- types.InputResult{
				Text:       string(data),
				IsExplicit: false,
			}
		}
	}()
}

func (c *ClipboardInputReader) GetChannel() chan types.InputResult {
	return c.ch
}
