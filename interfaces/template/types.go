package template

type ContntElement struct {
	IsTag   bool
	Content string
}

type TPLMemoViewData struct {
	CreatedAt       string
	ContentElements []ContntElement
	UserName        string
}
