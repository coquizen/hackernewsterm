package models

// Message struct for requesting data from server
type Request struct {
	RequestType string `json:"requesttype"`
	Payload     string `json:"payload"`
}

// Story is the posting data type
type Story struct {
	ID    int32  `json:"id"`
	Score int32  `json:"score"`
	Title string `json:"title"`
	By    string `json:"by"`
	Time  int32  `json:"time"`
	URL   string `json:"url"`
	Type  string `json:"type"`
}

// Item is a structure represents all the possible fields
type Item struct {
	ID          int      `json:"id"`
	Deleted     bool     `json:"deleted,omitempty"`
	Type        string   `json:"type,omitempty"`
	By          string   `json:"by,omitempty"`
	Time        int      `json:"time,omitempty"`
	Text        string   `json:"text,omitempty"`
	Title       string   `json:"title,omitempty"`
	Dead        bool     `json:"dead,omitempty"`
	Parent      int      `json:"parent,omitempty"`
	Poll        int      `json:"poll,omitempty"`
	Kids        []int    `json:"kids,omitempty"`
	URL         string   `json:"url,omitempty"`
	Score       int      `json:"score,omitempty"`
	Parts       []string `json:"parts,omitempty"`
	Descendants int      `json:"descendants,omitempty"`
}


