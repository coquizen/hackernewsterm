package hackernews

type Comment struct {
	By     string `json:"by"`
	ID     int    `json:"id"`
	Kids   []int  `json:"kids"`
	Parent int    `json:"parent"`
	Text   string `json:"text"`
	Time   int    `json:"time"`
	Type   string `json:"type"`
}

func (i item) ToComment() *Comment {
	return &Comment{
		By:     i.By(),
		ID:     i.ID(),
		Kids:   i.Kids(),
		Parent: i.Parent(),
		Text:   i.Text(),
		Time:   i.Time(),
		Type:   i.Type(),
	}
}
