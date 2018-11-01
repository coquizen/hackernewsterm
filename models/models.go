package models


import "time"
// Message struct for requesting data from server
type Message struct {
	RequestType string
	Payload     string
}

// Story is the posting data type
type Story struct {
	ID    int32  `json:"id"`
	Rank  int32  `json:"score"`
	Title string `json:"title"`
	By    string `json:"by"`
	Time  int32  `json:"time"`
	URL   string `json:"url"`
	Type  string `json:"type"`
}

// This structure represents all the possible fields
type Item struct {
	ID           int32
	Deleted      bool
	Type         string
	By           string
	PublishedAt time.Time
	Text         string
	Dead         bool
	Parent       int32
	Poll         int32
	Kids         int32
	URL          string
	Score        int32
	Parts        []string
	Descendants  int32
}
