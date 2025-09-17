package types

type SmartCutHook struct {
	Index          int
	Title          string
	PromptTemplate string
	PropertyName   string
}

type SmartCutConfig struct {
	Model string
	Hooks []*SmartCutHook
}
