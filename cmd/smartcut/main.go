package main

import (
	"fmt"
	"strings"

	"github.com/go-vgo/robotgo"
	"github.com/robotn/gohook"
)

func main() {
	fmt.Println("Text replacer app started. Press Ctrl+Shift+X to grab selected text, print it, and replace it with uppercase version.")
	fmt.Println("Press Ctrl+Shift+Q to quit.")

	// Register quit hotkey
	gohook.Register(gohook.KeyDown, []string{"q", "ctrl", "shift"}, func(e gohook.Event) {
		fmt.Println("Quitting...")
		gohook.End()
	})

	// Register the main hotkey
	gohook.Register(gohook.KeyDown, []string{"x", "ctrl", "shift"}, func(e gohook.Event) {
		// Simulate Ctrl+C to copy selected text
		robotgo.KeyTap("c", "ctrl")

		// Read from clipboard
		text, err := robotgo.ReadAll()
		if err != nil {
			fmt.Println("Error reading clipboard:", err)
			return
		}

		if text == "" {
			fmt.Println("No text selected.")
			return
		}

		fmt.Println("Selected text:", text)

		// Replace with something else (e.g., uppercase)
		newText := strings.ToUpper(text)
		fmt.Println("Replacing with:", newText)

		// Type the new text to replace the selection
		robotgo.TypeStr(newText)
	})

	s := gohook.Start()
	<-gohook.Process(s)
}
