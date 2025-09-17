package inputreader

import (
	"context"

	"golang.design/x/clipboard"
)

type ClipboardInputReader struct {
	context context.Context
	ch      chan string
}

func NewClipboardInputReader(context context.Context) *ClipboardInputReader {
	return &ClipboardInputReader{
		context: context,
		ch:      make(chan string),
	}
}

func (c *ClipboardInputReader) Start() {
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}

	clipch := clipboard.Watch(c.context, clipboard.FmtText)

	go func() {
		for data := range clipch {
			c.ch <- string(data)
		}
	}()
}

func (c *ClipboardInputReader) GetChannel() chan string {
	return c.ch
}
