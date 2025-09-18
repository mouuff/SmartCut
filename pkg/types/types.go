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

// InputResult represents the result of a user input
type InputResult struct {
	IsExplicit bool
	Text       string
}

type Brain interface {
	GenerateString(ctx context.Context, propertyName, prompt string) (string, error)
}

type InputReader interface {
	GetChannel() chan InputResult
}
