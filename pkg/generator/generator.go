package generator

import (
	"context"
	"log"
	"strings"

	"github.com/mouuff/SmartCuts/pkg/types"
)

type ResultsGeneratorImpl struct {
	lastTextRead string
	inputReader  types.InputReader
	ch           chan types.GenerationResult
	context      context.Context
	brain        types.Brain
	Config       *types.SmartCutConfig
}

func NewResultGenerator(
	context context.Context,
	brain types.Brain,
	inputReader types.InputReader,
	config *types.SmartCutConfig) *ResultsGeneratorImpl {
	return &ResultsGeneratorImpl{
		inputReader: inputReader,
		ch:          make(chan types.GenerationResult),
		context:     context,
		Config:      config,
		brain:       brain,
	}
}

func (o *ResultsGeneratorImpl) Start() {
	go func() {
		// Listen to clipboard changes
		for data := range o.inputReader.GetChannel() {
			o.generateForString(data.Text)
		}

		panic("unreachable")
	}()
}

func (o *ResultsGeneratorImpl) GetChannel() chan types.GenerationResult {
	return o.ch
}

func (o *ResultsGeneratorImpl) ReGenerate() {
	o.generateForString(o.lastTextRead)
}

func (o *ResultsGeneratorImpl) generateForString(data string) {
	o.lastTextRead = data

	if o.Config.Debug {
		log.Println("Clipboard changed:", data)
	}

	// For each prompt config, generate the result
	for _, promptConfig := range o.Config.PromptConfigs {
		go o.generateForPromptConfig(data, promptConfig)
	}
}

func (o *ResultsGeneratorImpl) generateForPromptConfig(clipboardText string, promptConfig *types.PromptConfig) {
	o.ch <- types.GenerationResult{
		Text:         "Generating...",
		PromptConfig: promptConfig,
	}

	prompt := strings.ReplaceAll(promptConfig.PromptTemplate, "{{input}}", clipboardText)
	result, err := o.brain.GenerateString(o.context, promptConfig.PropertyName, prompt)

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
