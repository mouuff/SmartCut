package types

import "strings"

// PromptConfig represents a single prompt configuration
type PromptConfig struct {
	Title          string
	PromptTemplate string
	PropertyName   string
}

// SmartCutsConfig represents the overall configuration for SmartCut
type SmartCutsConfig struct {
	ConfigPath     string `json:"-"`
	HostUrl        string
	Model          string
	MinRowsVisible int
	Debug          bool
	PromptConfigs  []PromptConfig
}

func (c *PromptConfig) GetPrompt(input string) string {
	return strings.ReplaceAll(c.PromptTemplate, "{{input}}", input)
}
