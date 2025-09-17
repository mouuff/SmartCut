package generator

import (
	"context"
	"log"
	"strings"

	"github.com/mouuff/SmartCuts/pkg/types"
	"golang.design/x/clipboard"
)

type ClipboardGenerator struct {
	Context context.Context
	Brain   types.Brain
	Config  *types.SmartCutConfig
	Out     chan types.GenerationResult
}

func NewClipboardGenerator(context context.Context, brain types.Brain, config *types.SmartCutConfig) *ClipboardGenerator {
	return &ClipboardGenerator{
		Context: context,
		Config:  config,
		Brain:   brain,
		Out:     make(chan types.GenerationResult),
	}
}

func (o *ClipboardGenerator) Start() {
	// Init returns an error if the package is not ready for use.
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}

	ch := clipboard.Watch(context.TODO(), clipboard.FmtText)

	// Listen to clipboard changes
	for data := range ch {
		o.GenerateForString(string(data))
	}

	panic("unreachable")
}

func (o *ClipboardGenerator) GenerateForString(data string) {
	// For each prompt config, generate the result
	for _, promptConfig := range o.Config.PromptConfigs {
		go o.GenerateForPromptConfig(data, promptConfig)
	}
}

func (o *ClipboardGenerator) GenerateForPromptConfig(clipboardText string, promptConfig *types.PromptConfig) {
	o.Out <- types.GenerationResult{
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

	o.Out <- types.GenerationResult{
		Text:         result,
		PromptConfig: promptConfig,
	}
}
