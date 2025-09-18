package types

import (
	"context"
)

// PromptConfig represents a single prompt configuration
type PromptConfig struct {
	Title          string
	PromptTemplate string
	PropertyName   string
}

// SmartCutConfig represents the overall configuration for SmartCut
type SmartCutConfig struct {
	HostUrl        string
	Model          string
	MinRowsVisible int
	Debug          bool
	PromptConfigs  []*PromptConfig
}

// InputText represents the result of a user input
type InputText struct {
	IsExplicit bool
	Text       string
}

type Brain interface {
	GenerateString(ctx context.Context, propertyName, prompt string) (string, error)
}

type InputReader interface {
	GetChannel() chan InputText
}
