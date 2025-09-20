package types

import (
	"strings"
)

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

func (c *PromptConfig) GetPrompt(input string) string {
	return strings.ReplaceAll(c.Prompt, "{{input}}", input)
}
