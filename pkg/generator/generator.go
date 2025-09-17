package generator

import (
	"context"
	"log"
	"strings"

	"github.com/mouuff/SmartCuts/pkg/types"
)

type ResultsGeneratorImpl struct {
	lastInput   types.InputResult
	inputReader types.InputReader
	ch          chan types.GenerationResult
	context     context.Context
	brain       types.Brain
	Config      *types.SmartCutConfig
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
			o.generateForInput(data)
		}

		panic("unreachable")
	}()
}

func (o *ResultsGeneratorImpl) GetChannel() chan types.GenerationResult {
	return o.ch
}

func (o *ResultsGeneratorImpl) ReGenerate() {
	o.generateForInput(o.lastInput)
}

func (o *ResultsGeneratorImpl) generateForInput(input types.InputResult) {
	o.lastInput = input

	if o.Config.Debug {
		log.Println("Clipboard changed:", input.Text)
	}

	// For each prompt config, generate the result
	for _, promptConfig := range o.Config.PromptConfigs {
		go o.generateForPromptConfig(input, promptConfig)
	}
}

func (o *ResultsGeneratorImpl) generateForPromptConfig(input types.InputResult, promptConfig *types.PromptConfig) {
	o.ch <- types.GenerationResult{
		Text:         "Generating...",
		PromptConfig: promptConfig,
		IsExplicit:   input.IsExplicit,
	}

	prompt := strings.ReplaceAll(promptConfig.PromptTemplate, "{{input}}", input.Text)
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
		IsExplicit:   false,
	}
}
