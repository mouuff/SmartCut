package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
	"github.com/mouuff/SmartCuts/pkg/brain"
	"github.com/mouuff/SmartCuts/pkg/controller"
	"github.com/mouuff/SmartCuts/pkg/inputreader"
	"github.com/mouuff/SmartCuts/pkg/types"
	"github.com/mouuff/SmartCuts/pkg/utils"
	"github.com/mouuff/SmartCuts/pkg/view"
	"github.com/mouuff/go-rocket-update/pkg/provider"
	"github.com/mouuff/go-rocket-update/pkg/updater"
	"golang.design/x/clipboard"
)

const SmartCutVersion string = "v0.0.3"

type SmartCutCmd struct {
	flagSet *flag.FlagSet

	config   string
	printBin bool
}

// Init initializes the command
func (cmd *SmartCutCmd) Init(args []string) error {
	cmd.flagSet = flag.NewFlagSet("smartcut", flag.ExitOnError)
	cmd.flagSet.StringVar(&cmd.config, "config", "", "configuration file override")
	cmd.flagSet.BoolVar(&cmd.printBin, "printBin", false, "prints the expected binary name")
	return cmd.flagSet.Parse(args)
}

// Run runs the command
func (cmd *SmartCutCmd) Run(a fyne.App, w fyne.Window) {
	if cmd.printBin {
		fmt.Printf("smartcut_%s_%s\n", runtime.GOOS, runtime.GOARCH)
		return
	}

	m := types.NewSmartCutModel()
	v := view.NewSmartCutView(w, m)
	w.SetContent(v.Layout())

	config, err := utils.GetOrCreateConfiguration(cmd.config)
	if err != nil {
		dialog.ShowError(err, w)
		return
	}

	m.UpdateFromConfig(config)

	b, err := brain.NewOllamaBrain(config.HostUrl, config.Model)
	if err != nil {
		dialog.ShowError(err, w)
		return
	}

	err = clipboard.Init()
	if err != nil {
		dialog.ShowError(err, w)
		return
	}

	// ir := inputreader.NewClipboardInputReader(ctx.Background())
	ir := inputreader.NewShortcutInputReader()
	ir.Start()

	c := controller.NewSmartCutController(context.Background(), b, m, config)

	// Setup View / Controller hooks
	ir.OnInput = c.GenerateForInput
	v.OnRequestGenerate = c.GenerateForInput
	c.OnRequestFocus = v.RequestFocus

	u := &updater.Updater{
		Provider: &provider.Github{
			RepositoryURL: "github.com/mouuff/SmartCuts",
			ArchiveName:   fmt.Sprintf("binaries_%s.zip", runtime.GOOS),
		},
		ExecutableName: fmt.Sprintf("smartcut_%s_%s", runtime.GOOS, runtime.GOARCH),
		Version:        SmartCutVersion,
	}

	if _, err := u.Update(); err != nil {
		log.Println(err)
	}
}

func main() {
	cmd := &SmartCutCmd{}
	err := cmd.Init(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	a := app.New()
	w := a.NewWindow("SmartCuts - " + SmartCutVersion)
	w.Resize(fyne.NewSize(800, 400))
	cmd.Run(a, w)
	w.ShowAndRun()
}
