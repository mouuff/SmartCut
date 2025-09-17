package main

import (
	"fmt"

	"github.com/mouuff/SmartCuts/pkg/inputreader"
)

func main() {
	ir := inputreader.NewShortcutInputReader()

	ir.Start()
	for data := range ir.GetChannel() {
		fmt.Println(data)
	}
}
