package main

import (
	"fmt"
	"strings"

	"github.com/go-vgo/robotgo"
)

func main() {
	fmt.Println("Text replacer app started. Press Ctrl+Shift+X to grab selected text, print it, and replace it with uppercase version.")
	fmt.Println("Press Ctrl+Shift+Q to quit.")

	// Register hotkeys
	quitHotkey := []string{"ctrl", "shift", "q"}
	replaceHotkey := []string{"ctrl", "shift", "x"}

	// Run hotkey listener in goroutine
	go func() {
		for {
			if robotgo.AddEvent(quitHotkey...) {
				fmt.Println("Quitting...")
				return
			}

			if robotgo.AddEvent(replaceHotkey...) {
				// Simulate Ctrl+C to copy selected text
				robotgo.KeyTap("c", "ctrl")

				// Read from clipboard
				text, err := robotgo.ReadAll()
				if err != nil {
					fmt.Println("Error reading clipboard:", err)
					continue
				}

				if text == "" {
					fmt.Println("No text selected.")
					continue
				}

				fmt.Println("Selected text:", text)

				// Replace with uppercase
				newText := strings.ToUpper(text)
				fmt.Println("Replacing with:", newText)

				// Type the new text
				robotgo.TypeStr(newText)
			}
		}
	}()

	// Keep program alive
	select {}
}
