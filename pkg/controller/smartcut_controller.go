package controller

import (
	"context"
	"sync"

	"github.com/mouuff/SmartCuts/pkg/types"
)

type SmartCutController struct {
	mu     sync.Mutex
	ctx    context.Context
	brain  types.Brain
	config *types.SmartCutConfig
	model  *types.SmartCutModel

	OnRequestFocus func()
}

func NewSmartCutController(
	ctx context.Context,
	brain types.Brain,
	model *types.SmartCutModel,
	config *types.SmartCutConfig) *SmartCutController {
	return &SmartCutController{
		ctx:            ctx,
		brain:          brain,
		config:         config,
		model:          model,
		OnRequestFocus: func() {},
	}
}

func (o *SmartCutController) UpdateItemContent(index int, content string) {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.model.UpdateResultItem(index, content)
}

func (o *SmartCutController) GenerateForInput(input types.InputText) {
	for i := range o.config.PromptConfigs {
		o.UpdateItemContent(i, "Generating...")
	}

	if input.IsExplicit {
		o.OnRequestFocus()
	}

	// For each prompt config, generate the result
	for index, promptConfig := range o.config.PromptConfigs {
		prompt := promptConfig.GetPrompt(input.Text)
		go o.generateItemContent(index, promptConfig.PropertyName, prompt)
	}
}

func (o *SmartCutController) generateItemContent(index int, propertyName, prompt string) {
	result, err := o.brain.GenerateString(o.ctx, propertyName, prompt)

	if err != nil {
		result = "Error while generating: " + err.Error()
	}

	o.UpdateItemContent(index, result)
}
