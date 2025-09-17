package types

import (
	"context"
)

// PromptConfig represents a single prompt configuration
type PromptConfig struct {
	Index          int
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

// GenerationResult represents the result of a generation
type GenerationResult struct {
	PromptConfig *PromptConfig
	OriginalText string
	Text         string
	IsExplicit   bool
}

// InputResult represents the result of a user input
type InputResult struct {
	IsExplicit bool
	Text       string
}

type Brain interface {
	GenerateString(ctx context.Context, propertyName, prompt string) (string, error)
}

type ResultsGenerator interface {
	GetChannel() chan GenerationResult
	GenerateForInput(input InputResult)
}

type InputReader interface {
	GetChannel() chan InputResult
}
