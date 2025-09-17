package inputreader

import (
	"bytes"
	"fmt"
	"syscall"
	"unsafe"

	"github.com/mouuff/SmartCuts/pkg/types"
	"golang.design/x/clipboard"
)

const (
	ModAlt = 1 << iota
	ModCtrl
	ModShift
	ModWin
)

type MSG struct {
	HWND   uintptr
	UINT   uintptr
	WPARAM int16
	LPARAM int64
	DWORD  int32
	POINT  struct{ X, Y int64 }
}

type Hotkey struct {
	Id        int // Unique id
	Modifiers int // Mask of modifiers
	KeyCode   int // Key code, e.g. 'A'
}

// String returns a human-friendly display name of the hotkey
// such as "Hotkey[Id: 1, Alt+Ctrl+O]"
func (h *Hotkey) String() string {
	mod := &bytes.Buffer{}
	if h.Modifiers&ModAlt != 0 {
		mod.WriteString("Alt+")
	}
	if h.Modifiers&ModCtrl != 0 {
		mod.WriteString("Ctrl+")
	}
	if h.Modifiers&ModShift != 0 {
		mod.WriteString("Shift+")
	}
	if h.Modifiers&ModWin != 0 {
		mod.WriteString("Win+")
	}
	return fmt.Sprintf("Hotkey[Id: %d, %s%c]", h.Id, mod, h.KeyCode)
}

type ShortcutInputReader struct {
	ch chan types.InputResult
}

func NewShortcutInputReader() *ShortcutInputReader {
	return &ShortcutInputReader{
		ch: make(chan types.InputResult),
	}
}

func (s *ShortcutInputReader) GetChannel() chan types.InputResult {
	return s.ch
}

func (s *ShortcutInputReader) Start() {
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}

	go func() {
		user32 := syscall.MustLoadDLL("user32")

		reghotkey := user32.MustFindProc("RegisterHotKey")
		// Hotkeys to listen to:
		keys := map[int16]*Hotkey{
			1: {1, ModAlt + ModShift, 'G'}, // ALT+SHIFT+G
		}

		// Register hotkeys:
		for _, v := range keys {
			r1, _, err := reghotkey.Call(
				0, uintptr(v.Id), uintptr(v.Modifiers), uintptr(v.KeyCode))
			if r1 == 1 {
				fmt.Println("Registered", v)
			} else {
				fmt.Println("Failed to register", v, ", error:", err)
			}
		}

		getmsg := user32.MustFindProc("GetMessageW")

		for {
			var msg = &MSG{}
			getmsg.Call(uintptr(unsafe.Pointer(msg)), 0, 0, 0)

			// Registered id is in the WPARAM field:
			if id := msg.WPARAM; id != 0 {
				fmt.Println("Hotkey pressed:", keys[id])

				if id == 1 {
					rawclip := clipboard.Read(clipboard.FmtText)

					if rawclip != nil {
						s.ch <- types.InputResult{
							Text:       string(rawclip),
							IsExplicit: true,
						}
					}
				}
			}
		}
	}()
}
