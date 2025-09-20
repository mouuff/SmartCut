package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
	"github.com/mouuff/SmartCut/internal"
	"github.com/mouuff/SmartCut/pkg/brain"
	"github.com/mouuff/SmartCut/pkg/controller"
	"github.com/mouuff/SmartCut/pkg/reader"
	"github.com/mouuff/SmartCut/pkg/types"
	"github.com/mouuff/SmartCut/pkg/utils"
	"github.com/mouuff/SmartCut/pkg/view"
	"golang.design/x/clipboard"
)

type SmartCutCmd struct {
	flagSet *flag.FlagSet

	config       string
	versionFlag  bool
	noUpdateFlag bool
}

// Init initializes the command
func (cmd *SmartCutCmd) Init(args []string) error {
	cmd.flagSet = flag.NewFlagSet("smartcut", flag.ExitOnError)
	cmd.flagSet.StringVar(&cmd.config, "config", "", "configuration file override")
	cmd.flagSet.BoolVar(&cmd.versionFlag, "version", false, "prints the version and exit")
	cmd.flagSet.BoolVar(&cmd.noUpdateFlag, "no-update", false, "disable automatic updates")
	return cmd.flagSet.Parse(args)
}

// Run runs the command
func (cmd *SmartCutCmd) Run(a fyne.App, w fyne.Window) error {
	config, err := utils.GetOrCreateConfiguration(cmd.config)
	if err != nil {
		return err
	}

	b, err := brain.NewOllamaBrain(config.HostUrl)
	if err != nil {
		return err
	}

	err = clipboard.Init()
	if err != nil {
		return err
	}

	// ir := reader.NewClipboardReader(ctx.Background())
	ir := reader.NewShortcutReader()
	ir.Start()

	// MVC setup
	m := types.NewSmartCutModel(config)
	v := view.NewSmartCutView(w, m)
	c := controller.NewSmartCutController(context.Background(), b, m, config)

	// Setup View / Controller hooks
	ir.OnInput = c.GenerateForInput
	v.OnRequestGenerate = c.GenerateForInput
	c.OnRequestFocus = v.RequestFocus

	w.SetContent(v.Layout())
	return nil
}

func main() {
	cmd := &SmartCutCmd{}
	err := cmd.Init(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	if cmd.versionFlag {
		fmt.Println(internal.SmartCutVersion)
		return
	}

	a := app.NewWithID("com.mouuff.smartcut")
	w := a.NewWindow("SmartCut - " + internal.SmartCutVersion)
	w.Resize(fyne.NewSize(800, 400))
	err = cmd.Run(a, w)

	if err != nil {
		dialog.ShowError(err, w)
	}

	w.ShowAndRun()

	if !cmd.noUpdateFlag {
		u := internal.GetSmartCutUpdater()
		err = u.SelfUpdate()
		if err != nil {
			log.Println(err)
		}
	}
}
