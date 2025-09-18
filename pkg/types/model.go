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
	configPath     string
	minRowsVisible int
	resultItems    []ResultItem

	OnChanged func()
}

func NewSmartCutModel() *SmartCutModel {
	return &SmartCutModel{
		configPath:     "",
		minRowsVisible: 6,
		resultItems:    make([]ResultItem, 0),
		OnChanged:      func() {},
	}
}

func (m *SmartCutModel) UpdateFromConfig(config *SmartCutConfig) {
	m.configPath = config.ConfigPath
	m.minRowsVisible = config.MinRowsVisible
	m.resultItems = make([]ResultItem, 0)

	for _, promptConfig := range config.PromptConfigs {
		m.resultItems = append(m.resultItems, ResultItem{
			Title:   promptConfig.Title,
			Content: "Waiting for generation...",
		})
	}
	m.OnChanged()
}

func (m *SmartCutModel) UpdateResultItem(index int, content string) {
	m.resultItems[index].Content = content
	m.OnChanged()
}

func (m *SmartCutModel) ResultItems() []ResultItem {
	return m.resultItems
}

func (m *SmartCutModel) MinRowsVisible() int {
	return m.minRowsVisible
}

func (m *SmartCutModel) ConfigPath() string {
	return m.configPath
}
