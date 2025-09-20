package types

import "strings"

// PromptConfig represents a single prompt configuration
type PromptConfig struct {
	Title          string
	SystemPrompt   string
	TemplatePrompt string
	PropertyName   string
}

// SmartCutConfig represents the overall configuration for SmartCut
type SmartCutConfig struct {
	ConfigPath     string `json:"-"`
	HostUrl        string
	Model          string
	MinRowsVisible int
	Debug          bool
	PromptConfigs  []PromptConfig
}

func (c *PromptConfig) GetPrompt(input string) string {
	return strings.ReplaceAll(c.TemplatePrompt, "{{input}}", input)
}
