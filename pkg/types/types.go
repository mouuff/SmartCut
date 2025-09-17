package types

type PromptConfig struct {
	Index          int
	Title          string
	PromptTemplate string
	PropertyName   string
}

type SmartCutConfig struct {
	Model          string
	MinRowsVisible int
	PromptConfigs  []*PromptConfig
}
