package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"golang.design/x/clipboard"
)

type SmartCutConfig struct {
	Model string
	Debug bool
}

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

	// Init returns an error if the package is not ready for use.
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}

	ch := clipboard.Watch(context.TODO(), clipboard.FmtText)
	for data := range ch {
		// print out clipboard data whenever it is changed
		println(string(data))
	}

	return nil
}

func printConfigurationTemplate() {
	configTemplate := &SmartCutConfig{
		Model: "mistral",
		Debug: false,
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
