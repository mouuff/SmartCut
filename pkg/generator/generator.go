package generator

import (
	"context"
	"log"
	"strings"

	"github.com/mouuff/SmartCuts/pkg/types"
	"golang.design/x/clipboard"
)

type Brain interface {
	GenerateString(ctx context.Context, propertyName, prompt string) (string, error)
}

type GenerationResult struct {
	PromptConfig  *types.PromptConfig
	ClipboardText string
	Text          string
}

type ClipboardGenerator struct {
	Context context.Context
	Brain   Brain
	Config  *types.SmartCutConfig
	Out     chan GenerationResult
}

func NewClipboardGenerator(context context.Context, brain Brain, config *types.SmartCutConfig) *ClipboardGenerator {
	return &ClipboardGenerator{
		Context: context,
		Config:  config,
		Brain:   brain,
		Out:     make(chan GenerationResult),
	}
}

func (o *ClipboardGenerator) Start() {
	// Init returns an error if the package is not ready for use.
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}

	ch := clipboard.Watch(context.TODO(), clipboard.FmtText)
	for data := range ch {
		// print out clipboard data whenever it is changed
		println(string(data))
		o.GenerateForString(string(data))
	}

	log.Println("Stopped feeding from clipboard.")
}

func (o *ClipboardGenerator) GenerateForString(clipboardText string) {
	for _, promptConfig := range o.Config.PromptConfigs {
		prompt := strings.ReplaceAll(promptConfig.PromptTemplate, "{{input}}", clipboardText)
		result, err := o.Brain.GenerateString(o.Context, promptConfig.PropertyName, prompt)

		if err != nil {
			log.Println("Error generating:", err)
			continue
		} else {
			log.Println("Generated:", result)
			o.Out <- GenerationResult{
				Text:         result,
				PromptConfig: promptConfig,
			}
		}
	}
}
