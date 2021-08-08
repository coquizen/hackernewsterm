package hackernews

type Story struct {
	By    string `json:"by"`
	ID    int    `json:"id"`
	Kids  []int  `json:"kids"`
	Score int    `json:"score"`
	Time  int    `json:"time"`
	Title string `json:"title"`
	Type  string `json:"type"`
	URL   string `json:"url"`
}

// ToStory converts the item to Story type
func (i *item) ToStory() *Story {
	return &Story{
		By:    i.By(),
		ID:    i.ID(),
		Kids:  i.Kids(),
		Score: i.Score(),
		Time:  i.Time(),
		Title: i.Title(),
		Type:  i.Type(),
		URL:   i.URL(),
	}
}
