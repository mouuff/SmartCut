package types

type InputText struct {
	IsExplicit bool
	Text       string
}

type ResultItem struct {
	Title   string
	Content string
}

type SmartCutsModel struct {
	config      *SmartCutsConfig
	resultItems []ResultItem

	OnChanged func()
}

func NewSmartCutsModel(config *SmartCutsConfig) *SmartCutsModel {
	resultItems := make([]ResultItem, 0)
	for _, promptConfig := range config.PromptConfigs {
		resultItems = append(resultItems, ResultItem{
			Title:   promptConfig.Title,
			Content: "Waiting for generation...",
		})
	}

	return &SmartCutsModel{
		config:      config,
		resultItems: resultItems,
		OnChanged:   func() {},
	}
}

func (m *SmartCutsModel) UpdateResultItem(index int, content string) {
	m.resultItems[index].Content = content
	m.OnChanged()
}

func (m *SmartCutsModel) ResultItems() []ResultItem {
	return m.resultItems
}

func (m *SmartCutsModel) Config() SmartCutsConfig {
	return *m.config
}
