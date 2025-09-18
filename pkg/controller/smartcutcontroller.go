package controller

import (
	"context"
	"log"
	"strings"
	"sync"

	"github.com/mouuff/SmartCuts/pkg/types"
	"github.com/mouuff/SmartCuts/pkg/view"
)

type SmartCutController struct {
	mu      sync.Mutex
	context context.Context
	brain   types.Brain
	config  *types.SmartCutConfig
	view    *view.SmartCutView
	model   *types.SmartCutModel
}

func getModelForConfig(config *types.SmartCutConfig) *types.SmartCutModel {
	model := &types.SmartCutModel{
		MinRowsVisible: config.MinRowsVisible,
		ResultItems:    make([]types.ResultItem, 0),
	}

	for _, promptConfig := range config.PromptConfigs {
		model.ResultItems = append(model.ResultItems, types.ResultItem{
			Title:   promptConfig.Title,
			Content: "Waiting for generation...",
		})
	}

	return model
}

func NewSmartCutController(
	context context.Context,
	brain types.Brain,
	view *view.SmartCutView,
	config *types.SmartCutConfig) *SmartCutController {
	return &SmartCutController{
		context: context,
		brain:   brain,
		config:  config,
		view:    view,
		model:   getModelForConfig(config),
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

func (o *SmartCutController) Init() {
	o.view.OnAskGenerate = o.GenerateForInput
	o.view.DoRefresh(o.model)
}

func (o *SmartCutController) UpdateItemContent(index int, content string) {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.model.ResultItems[index].Content = content
	o.view.DoRefresh(o.model)
}

func (o *SmartCutController) GenerateForInput(input types.InputText) {
	if o.config.Debug {
		log.Println("GenerateForInput:", input.Text)
	}

	for i := range o.config.PromptConfigs {
		o.UpdateItemContent(i, "Generating...")
	}

	if input.IsExplicit {
		o.view.RequestFocus()
	}

	// For each prompt config, generate the result
	for _, promptConfig := range o.config.PromptConfigs {
		go func() {
			prompt := strings.ReplaceAll(promptConfig.PromptTemplate, "{{input}}", input.Text)
			result, err := o.brain.GenerateString(o.context, promptConfig.PropertyName, prompt)

			if err != nil {
				result = "Error while generating: " + err.Error()
			}

			o.UpdateItemContent(promptConfig.Index, result)
		}()
	}
}
