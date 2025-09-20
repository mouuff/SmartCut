package types

import (
	"context"
	"strings"
)

// PromptConfig represents a prompt request
type PromptRequest struct {
	Model        string
	SystemPrompt string
	Prompt       string
	PropertyName string
}

type Brain interface {
	GenerateString(ctx context.Context, r *PromptRequest) (string, error)
}

// PromptConfig represents a single prompt configuration
type PromptConfig struct {
	Model        string
	Title        string
	SystemPrompt string
	Prompt       string
	PropertyName string
}

// SmartCutConfig represents the overall configuration for SmartCut
type SmartCutConfig struct {
	ConfigPath     string `json:"-"`
	HostUrl        string
	MinRowsVisible int
	Debug          bool
	PromptConfigs  []PromptConfig
}

func (c *PromptConfig) GetPromptRequest(input string) *PromptRequest {
	return &PromptRequest{
		Model:        c.Model,
		SystemPrompt: c.SystemPrompt,
		PropertyName: c.PropertyName,
		Prompt:       strings.ReplaceAll(c.Prompt, "{{input}}", input),
	}
}
