package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	smartcutapp "github.com/mouuff/SmartCuts/pkg/app"
	"github.com/mouuff/SmartCuts/pkg/brain"
	"github.com/mouuff/SmartCuts/pkg/orchestrator"
	"github.com/mouuff/SmartCuts/pkg/types"
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
	if cmd.config == "" {
		log.Println("Please pass the configuration file using -config")
		log.Println("Here is an example configuration:")
		printConfigurationTemplate()
		return errors.New("-config parameter required")
	}

	var config types.SmartCutConfig
	err := utils.ReadFromJson(cmd.config, &config)
	if err != nil {
		return err
	}

	b, err := brain.NewOllamaBrain("llama3.2")

	if err != nil {
		panic(err)
	}

	o := orchestrator.NewOrchestrator(context.Background(), b, &config)

	go o.StartFeedFromClipboard()

	a := app.New()
	w := a.NewWindow("SmartCuts")

	smartcut := smartcutapp.NewSmartCutApp(w, &config)

	go func() {
		for res := range o.Out {
			fyne.Do(func() {
				smartcut.UpdateItem(res)
			})
		}
	}()

	w.SetContent(smartcut.Layout())
	w.Resize(fyne.NewSize(800, 400))
	w.ShowAndRun()

	return nil
}

func printConfigurationTemplate() {
	configTemplate := &types.SmartCutConfig{
		Model:          "mistral",
		MinRowsVisible: 7,
		PromptConfigs: []*types.PromptConfig{
			&types.PromptConfig{
				Title:          "Translate to French",
				PropertyName:   "translated_text",
				PromptTemplate: "Translate this to french: '{{input}}'",
			},
		},
	}

	jsonData, err := json.MarshalIndent(configTemplate, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(jsonData))
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
