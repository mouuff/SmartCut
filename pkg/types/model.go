package types

type InputText struct {
	IsExplicit bool
	Text       string
}

type ResultItem struct {
	Title   string
	Content string
}

type SmartCutModel struct {
	config      *SmartCutConfig
	resultItems []ResultItem

	OnChanged func()
}

func NewSmartCutModel(config *SmartCutConfig) *SmartCutModel {
	resultItems := make([]ResultItem, 0)
	for _, promptConfig := range config.PromptConfigs {
		resultItems = append(resultItems, ResultItem{
			Title:   promptConfig.Title,
			Content: "Waiting for generation...",
		})
	}

	return &SmartCutModel{
		config:      config,
		resultItems: resultItems,
		OnChanged:   func() {},
	}
}

func (m *SmartCutModel) UpdateResultItem(index int, content string) {
	m.resultItems[index].Content = content
	m.OnChanged()
}

func (m *SmartCutModel) ResultItems() []ResultItem {
	return m.resultItems
}

func (m *SmartCutModel) Config() SmartCutConfig {
	return *m.config
}
