package generator

import (
	"context"
	"log"
	"strings"

	"github.com/mouuff/SmartCuts/pkg/types"
	"golang.design/x/clipboard"
)

type ClipboardGenerator struct {
	currentClipboard string
	ch               chan types.GenerationResult
	Context          context.Context
	Brain            types.Brain
	Config           *types.SmartCutConfig
}

func NewClipboardGenerator(context context.Context, brain types.Brain, config *types.SmartCutConfig) *ClipboardGenerator {
	return &ClipboardGenerator{
		ch:      make(chan types.GenerationResult),
		Context: context,
		Config:  config,
		Brain:   brain,
	}
}

func (o *ClipboardGenerator) Start() {
	// Init returns an error if the package is not ready for use.
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}

	ch := clipboard.Watch(o.Context, clipboard.FmtText)

	// Listen to clipboard changes
	for data := range ch {
		o.generateForString(string(data))
	}

	panic("unreachable")
}

func (o *ClipboardGenerator) GetChannel() chan types.GenerationResult {
	return o.ch
}

func (o *ClipboardGenerator) ReGenerate() {
	o.generateForString(o.currentClipboard)
}

func (o *ClipboardGenerator) generateForString(data string) {
	o.currentClipboard = data

	if o.Config.Debug {
		log.Println("Clipboard changed:", data)
	}

	// For each prompt config, generate the result
	for _, promptConfig := range o.Config.PromptConfigs {
		go o.generateForPromptConfig(data, promptConfig)
	}
}

func (o *ClipboardGenerator) generateForPromptConfig(clipboardText string, promptConfig *types.PromptConfig) {
	o.ch <- types.GenerationResult{
		Text:         "Generating...",
		PromptConfig: promptConfig,
	}

	prompt := strings.ReplaceAll(promptConfig.PromptTemplate, "{{input}}", clipboardText)
	result, err := o.Brain.GenerateString(o.Context, promptConfig.PropertyName, prompt)

	if o.Config.Debug {
		if err != nil {
			log.Println("Error generating:", err)
		} else {
			log.Println("Generated: ", result)
		}
	}

	if err != nil {
		result = "Error while generating: " + err.Error()
	}

	o.ch <- types.GenerationResult{
		Text:         result,
		PromptConfig: promptConfig,
	}
}
