package reader

import (
	"github.com/mouuff/SmartCut/pkg/types"
	"golang.design/x/clipboard"
	"golang.design/x/hotkey"
)

type ShortcutReaderV2 struct {
	OnInput func(types.InputText)
}

func NewShortcutReaderV2() *ShortcutReaderV2 {
	return &ShortcutReaderV2{
		OnInput: func(types.InputText) {},
	}
}

func (s *ShortcutReaderV2) Start() {

	go func() {
		hk := hotkey.New([]hotkey.Modifier{hotkey.ModAlt, hotkey.ModShift}, hotkey.KeyO)
		err := hk.Register()
		if err != nil {
			panic(err)
		}

		for range hk.Keydown() {
			rawclip := clipboard.Read(clipboard.FmtText)

			if rawclip != nil {
				s.OnInput(types.InputText{
					Text:       string(rawclip),
					IsExplicit: true,
				})
			}
		}
	}()

}
