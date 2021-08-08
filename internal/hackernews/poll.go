package hackernews

type Poll struct {
	By    string `json:"by"`
	ID    int    `json:"id"`
	Kids  []int  `json:"kids"`
	Parts []int  `json:"parts"`
	Score int    `json:"score"`
	Text  string `json:"text"`
	Time  int    `json:"time"`
	Title string `json:"title"`
	Type  string `json:"type"`
}

// ToPoll converts item type to Poll type

func (i item) ToPoll() *Poll {
	return &Poll{
		By:    i.By(),
		ID:    i.ID(),
		Kids:  i.Kids(),
		Parts: i.Parts(),
		Score: i.Score(),
		Text:  i.Text(),
		Time:  i.Time(),
		Title: i.Title(),
		Type:  i.Type(),
	}
}
