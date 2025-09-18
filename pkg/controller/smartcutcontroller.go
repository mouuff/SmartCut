package controller

import (
	"context"
	"log"
	"strings"
	"sync"

	"github.com/mouuff/SmartCuts/pkg/types"
)

type SmartCutController struct {
	mu      sync.Mutex
	context context.Context
	brain   types.Brain
	config  *types.SmartCutConfig
	model   *types.SmartCutModel

	OnRequestFocus func()
}

func NewSmartCutController(
	context context.Context,
	brain types.Brain,
	model *types.SmartCutModel,
	config *types.SmartCutConfig) *SmartCutController {
	return &SmartCutController{
		context:        context,
		brain:          brain,
		config:         config,
		model:          model,
		OnRequestFocus: func() {},
	}
}

func (o *SmartCutController) ListenTo(inputReader types.InputReader) {
	go func() {
		// Listen to clipboard changes
		for data := range inputReader.GetChannel() {
			o.GenerateForInput(data)
		}

		panic("unreachable")
	}()
}

func (o *SmartCutController) UpdateItemContent(index int, content string) {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.model.UpdateResultItem(index, content)
}

func (o *SmartCutController) GenerateForInput(input types.InputText) {
	if o.config.Debug {
		log.Println("GenerateForInput:", input.Text)
	}

	for i := range o.config.PromptConfigs {
		o.UpdateItemContent(i, "Generating...")
	}

	if input.IsExplicit {
		o.OnRequestFocus()
	}

	// For each prompt config, generate the result
	for index, promptConfig := range o.config.PromptConfigs {
		go func() {
			prompt := strings.ReplaceAll(promptConfig.PromptTemplate, "{{input}}", input.Text)
			result, err := o.brain.GenerateString(o.context, promptConfig.PropertyName, prompt)

			if err != nil {
				result = "Error while generating: " + err.Error()
			}

			o.UpdateItemContent(index, result)
		}()
	}
}
