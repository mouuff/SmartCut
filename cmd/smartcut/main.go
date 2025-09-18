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
	"github.com/mouuff/SmartCuts/pkg/brain"
	"github.com/mouuff/SmartCuts/pkg/controller"
	"github.com/mouuff/SmartCuts/pkg/inputreader"
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
func (cmd *SmartCutCmd) Run() error {
	if cmd.printBin {
		fmt.Printf("smartcut_%s_%s\n", runtime.GOOS, runtime.GOARCH)
		return nil
	}

	config, err := utils.GetOrCreateConfiguration(cmd.config)
	if err != nil {
		return err
	}

	b, err := brain.NewOllamaBrain(config.HostUrl, config.Model)
	if err != nil {
		panic(err)
	}

	err = clipboard.Init()
	if err != nil {
		panic(err)
	}

	// ir := inputreader.NewClipboardInputReader(context.Background())
	ir := inputreader.NewShortcutInputReader()
	ir.Start()

	a := app.New()
	w := a.NewWindow("SmartCuts - " + SmartCutVersion)

	smartcutview := view.NewSmartCutView(w)

	rg := controller.NewSmartCutController(context.Background(), b, smartcutview, config)
	smartcutview.OnAskGenerate = rg.GenerateForInput
	rg.ListenTo(ir)
	rg.RefreshView()

	w.SetContent(smartcutview.Layout())
	w.Resize(fyne.NewSize(800, 400))
	w.ShowAndRun()

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

	return nil
}

func main() {
	cmd := &SmartCutCmd{}
	err := cmd.Init(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
