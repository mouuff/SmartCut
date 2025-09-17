package main

import (
	"context"
	"flag"
	"log"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	smartcutapp "github.com/mouuff/SmartCuts/pkg/app"
	"github.com/mouuff/SmartCuts/pkg/brain"
	"github.com/mouuff/SmartCuts/pkg/generator"
	"github.com/mouuff/SmartCuts/pkg/inputreader"
	"github.com/mouuff/SmartCuts/pkg/utils"
)

type SmartCutCmd struct {
	flagSet *flag.FlagSet

	config string
}

// Init initializes the command
func (cmd *SmartCutCmd) Init(args []string) error {
	cmd.flagSet = flag.NewFlagSet("smartcut", flag.ExitOnError)
	cmd.flagSet.StringVar(&cmd.config, "config", "", "configuration file (required)")
	return cmd.flagSet.Parse(args)
}

// Run runs the command
func (cmd *SmartCutCmd) Run() error {
	config, err := utils.GetOrCreateConfiguration(cmd.config)
	if err != nil {
		return err
	}

	b, err := brain.NewOllamaBrain(config.Model)

	if err != nil {
		panic(err)
	}

	ir := inputreader.NewClipboardInputReader(context.Background())
	rg := generator.NewResultGenerator(context.Background(), b, ir, config)

	a := app.New()
	w := a.NewWindow("SmartCuts")

	smartcut := smartcutapp.NewSmartCutApp(w, config, rg)

	// Start the input reader
	go ir.Start()

	// Start listening to clipboard results
	go smartcut.Start()

	// Start the generator
	go rg.Start()

	w.SetContent(smartcut.Layout())
	w.Resize(fyne.NewSize(800, 400))
	w.ShowAndRun()
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
