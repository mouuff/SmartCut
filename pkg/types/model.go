package types

type ResultItem struct {
	Title   string
	Content string
}

type SmartCutModel struct {
	ConfigFilePath string
	MinRowsVisible int
	ResultItems    []ResultItem
}
