package types

import (
	"context"
)

type PromptConfig struct {
	Index          int
	Title          string
	PromptTemplate string
	PropertyName   string
}

type SmartCutConfig struct {
	Model          string
	MinRowsVisible int
	Debug          bool
	PromptConfigs  []*PromptConfig
}

type Brain interface {
	GenerateString(ctx context.Context, propertyName, prompt string) (string, error)
}

type GenerationResult struct {
	PromptConfig  *PromptConfig
	ClipboardText string
	Text          string
}
